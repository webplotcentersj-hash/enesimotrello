#!/bin/bash

# Setup Nginx Reverse Proxy with SSL for TaskBoard
# Run this after deploying with docker-compose

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

echo "🔒 TaskBoard Nginx + SSL Setup"
echo "=============================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    print_error "Please run this script as root (use sudo)"
    exit 1
fi

# Get domain information
read -p "Enter your frontend domain (e.g., taskboard.com): " FRONTEND_DOMAIN
read -p "Enter your backend/API domain (e.g., api.taskboard.com): " BACKEND_DOMAIN
read -p "Enter your email for SSL certificates: " EMAIL

if [ -z "$FRONTEND_DOMAIN" ] || [ -z "$BACKEND_DOMAIN" ] || [ -z "$EMAIL" ]; then
    print_error "All fields are required"
    exit 1
fi

# Install Nginx
echo ""
echo "📦 Installing Nginx..."
apt-get update
apt-get install -y nginx

print_success "Nginx installed"

# Install Certbot
echo ""
echo "📦 Installing Certbot..."
apt-get install -y certbot python3-certbot-nginx

print_success "Certbot installed"

# Create Nginx configuration
echo ""
echo "📝 Creating Nginx configuration..."

cat > /etc/nginx/sites-available/taskboard << EOF
# TaskBoard Backend API
server {
    listen 80;
    server_name $BACKEND_DOMAIN;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
        
        # WebSocket support
        proxy_read_timeout 86400;
    }
}

# TaskBoard Frontend
server {
    listen 80;
    server_name $FRONTEND_DOMAIN;

    location / {
        proxy_pass http://localhost:80;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
        
        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;
    }
}
EOF

# Enable site
ln -sf /etc/nginx/sites-available/taskboard /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Test Nginx configuration
echo ""
echo "🔍 Testing Nginx configuration..."
if nginx -t; then
    print_success "Nginx configuration is valid"
else
    print_error "Nginx configuration has errors"
    exit 1
fi

# Reload Nginx
systemctl reload nginx
print_success "Nginx reloaded"

# Setup SSL certificates
echo ""
echo "🔒 Setting up SSL certificates..."
echo "This may take a minute..."

certbot --nginx -d $FRONTEND_DOMAIN -d $BACKEND_DOMAIN \
    --non-interactive \
    --agree-tos \
    --email $EMAIL \
    --redirect

if [ $? -eq 0 ]; then
    print_success "SSL certificates installed successfully"
else
    print_error "Failed to install SSL certificates"
    print_warning "Make sure your domains point to this server's IP"
    exit 1
fi

# Setup auto-renewal
echo ""
echo "⚙️  Setting up automatic certificate renewal..."
systemctl enable certbot.timer
systemctl start certbot.timer

print_success "Auto-renewal configured"

# Configure firewall
echo ""
echo "🔥 Configuring firewall..."
if command -v ufw &> /dev/null; then
    ufw allow 'Nginx Full'
    ufw delete allow 'Nginx HTTP'
    print_success "Firewall configured"
else
    print_warning "ufw not found, please configure firewall manually"
fi

# Display summary
echo ""
echo "=============================="
print_success "Setup Complete!"
echo "=============================="
echo ""
echo "🌐 Your TaskBoard is now accessible at:"
echo "   Frontend: https://$FRONTEND_DOMAIN"
echo "   Backend:  https://$BACKEND_DOMAIN"
echo ""
echo "🔒 SSL certificates will auto-renew"
echo ""
echo "📝 Update your .env file with the HTTPS URLs:"
echo "   REACT_APP_API_URL=https://$BACKEND_DOMAIN/api/v1"
echo "   CORS_ORIGIN=https://$FRONTEND_DOMAIN"
echo ""
echo "Then restart the services:"
echo "   docker compose -f docker-compose.prod.yml restart"
echo ""

