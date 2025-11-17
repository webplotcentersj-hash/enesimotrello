#!/bin/bash

# Setup Nginx Reverse Proxy for n8n + TaskBoard Backend (Frontend on Vercel)
# This script configures Nginx to serve n8n and TaskBoard backend API
# Frontend is deployed separately on Vercel

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }

echo "🔒 Nginx Setup (n8n + TaskBoard Backend)"
echo "========================================"
echo "Frontend will be deployed on Vercel"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    print_error "Please run this script as root (use sudo)"
    exit 1
fi

# Get domain information
echo "📝 Domain Configuration"
echo "----------------------"
read -p "Enter n8n domain (e.g., n8n.tudominio.com): " N8N_DOMAIN
read -p "Enter TaskBoard API domain (e.g., api.taskboard.tudominio.com): " TASKBOARD_API_DOMAIN
read -p "Enter Vercel frontend domain (e.g., taskboard.vercel.app or taskboard.tudominio.com): " VERCEL_DOMAIN
read -p "Enter your email for SSL certificates: " EMAIL

if [ -z "$N8N_DOMAIN" ] || [ -z "$TASKBOARD_API_DOMAIN" ] || [ -z "$VERCEL_DOMAIN" ] || [ -z "$EMAIL" ]; then
    print_error "All fields are required"
    exit 1
fi

# Detect n8n port (default is 5678)
read -p "Enter n8n port (default: 5678): " N8N_PORT
N8N_PORT=${N8N_PORT:-5678}

# Install Nginx if not installed
if ! command -v nginx &> /dev/null; then
    echo ""
    echo "📦 Installing Nginx..."
    apt-get update
    apt-get install -y nginx
    print_success "Nginx installed"
else
    print_success "Nginx is already installed"
fi

# Install Certbot if not installed
if ! command -v certbot &> /dev/null; then
    echo ""
    echo "📦 Installing Certbot..."
    apt-get install -y certbot python3-certbot-nginx
    print_success "Certbot installed"
else
    print_success "Certbot is already installed"
fi

# Create Nginx configuration for n8n
echo ""
echo "📝 Creating Nginx configuration for n8n..."
cat > /etc/nginx/sites-available/n8n << EOF
# n8n Application
server {
    listen 80;
    server_name $N8N_DOMAIN;

    # Increase body size for file uploads
    client_max_body_size 50M;

    location / {
        proxy_pass http://127.0.0.1:$N8N_PORT;
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
EOF

# Create Nginx configuration for TaskBoard Backend
echo "📝 Creating Nginx configuration for TaskBoard Backend..."
cat > /etc/nginx/sites-available/taskboard-api << EOF
# TaskBoard Backend API
server {
    listen 80;
    server_name $TASKBOARD_API_DOMAIN;

    # CORS headers (will be handled by backend, but good to have here too)
    add_header 'Access-Control-Allow-Origin' '$VERCEL_DOMAIN' always;
    add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
    add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, X-Anonymous-User-Id' always;
    add_header 'Access-Control-Allow-Credentials' 'true' always;

    # Handle preflight requests
    if (\$request_method = 'OPTIONS') {
        add_header 'Access-Control-Allow-Origin' '$VERCEL_DOMAIN' always;
        add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
        add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, X-Anonymous-User-Id' always;
        add_header 'Access-Control-Max-Age' 1728000;
        add_header 'Content-Type' 'text/plain; charset=utf-8';
        add_header 'Content-Length' 0;
        return 204;
    }

    location / {
        proxy_pass http://127.0.0.1:8081;
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
EOF

# Enable sites
echo ""
echo "🔗 Enabling Nginx sites..."
ln -sf /etc/nginx/sites-available/n8n /etc/nginx/sites-enabled/
ln -sf /etc/nginx/sites-available/taskboard-api /etc/nginx/sites-enabled/

# Remove default site if it exists
if [ -f /etc/nginx/sites-enabled/default ]; then
    rm /etc/nginx/sites-enabled/default
    print_info "Removed default Nginx site"
fi

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
echo ""

certbot --nginx -d $N8N_DOMAIN -d $TASKBOARD_API_DOMAIN \
    --non-interactive \
    --agree-tos \
    --email $EMAIL \
    --redirect

if [ $? -eq 0 ]; then
    print_success "SSL certificates installed successfully"
else
    print_error "Failed to install SSL certificates"
    print_warning "Make sure your domains point to this server's IP: $(curl -s ifconfig.me)"
    print_warning "You can run certbot manually later: sudo certbot --nginx"
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
    ufw delete allow 'Nginx HTTP' 2>/dev/null || true
    print_success "Firewall configured"
else
    print_warning "ufw not found, please configure firewall manually"
fi

# Display summary
echo ""
echo "========================================"
print_success "Setup Complete!"
echo "========================================"
echo ""
echo "🌐 Your applications are now accessible at:"
echo ""
echo "   📊 n8n:"
echo "      http://$N8N_DOMAIN"
echo "      https://$N8N_DOMAIN"
echo ""
echo "   🔌 TaskBoard Backend API:"
echo "      http://$TASKBOARD_API_DOMAIN"
echo "      https://$TASKBOARD_API_DOMAIN"
echo ""
echo "   🎨 TaskBoard Frontend (Vercel):"
echo "      https://$VERCEL_DOMAIN"
echo ""
echo "🔒 SSL certificates will auto-renew"
echo ""
echo "📝 IMPORTANT: Update your TaskBoard .env file with:"
echo "   CORS_ORIGIN=https://$VERCEL_DOMAIN"
echo ""
echo "And set in Vercel environment variables:"
echo "   REACT_APP_API_URL=https://$TASKBOARD_API_DOMAIN/api/v1"
echo ""
echo "Then restart TaskBoard backend:"
echo "   cd /path/to/task-board"
echo "   docker compose -f docker-compose.backend-only.yml restart"
echo ""
echo "📊 Check service status:"
echo "   sudo systemctl status nginx"
echo "   docker ps"
echo ""

