#!/bin/bash

# TaskBoard VPS Deployment Script
# This script automates deployment to a VPS/dedicated server

set -e

echo "🚀 TaskBoard VPS Deployment Script"
echo "=================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    echo -e "${RED}❌ Please don't run this script as root${NC}"
    exit 1
fi

# Function to print colored output
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "📦 Docker not found. Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
    print_success "Docker installed successfully"
    print_warning "Please log out and log back in for Docker permissions to take effect"
    print_warning "Then run this script again"
    exit 0
fi

# Check if Docker Compose is installed
if ! command -v docker compose &> /dev/null; then
    echo "📦 Docker Compose not found. Installing..."
    sudo apt-get update
    sudo apt-get install -y docker-compose-plugin
    print_success "Docker Compose installed successfully"
fi

# Check if .env file exists
if [ ! -f .env ]; then
    print_warning ".env file not found"
    if [ -f env.production.example ]; then
        echo "📝 Creating .env file from template..."
        cp env.production.example .env
        print_success ".env file created"
        echo ""
        echo "⚠️  IMPORTANT: Please edit .env file with your production values:"
        echo "   - DB_PASSWORD (strong password)"
        echo "   - JWT_SECRET (min 32 characters)"
        echo "   - CORS_ORIGIN (your domain)"
        echo "   - REACT_APP_API_URL (your API URL)"
        echo ""
        read -p "Press enter after you've updated the .env file..."
    else
        print_error "env.production.example not found"
        exit 1
    fi
fi

# Validate required environment variables
echo "🔍 Validating environment variables..."
source .env

required_vars=("DB_PASSWORD" "JWT_SECRET" "CORS_ORIGIN" "REACT_APP_API_URL")
missing_vars=()

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ] || [ "${!var}" == "changeme"* ]; then
        missing_vars+=("$var")
    fi
done

if [ ${#missing_vars[@]} -ne 0 ]; then
    print_error "The following environment variables need to be set:"
    for var in "${missing_vars[@]}"; do
        echo "   - $var"
    done
    exit 1
fi

print_success "Environment variables validated"

# Pull latest code
echo ""
echo "📥 Pulling latest code..."
if [ -d .git ]; then
    git pull origin main || print_warning "Could not pull latest code (continuing anyway)"
else
    print_warning "Not a git repository (skipping pull)"
fi

# Build and deploy
echo ""
echo "🏗️  Building and deploying TaskBoard..."
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d --build

# Wait for services to be healthy
echo ""
echo "⏳ Waiting for services to start..."
sleep 10

# Check service status
echo ""
echo "🔍 Checking service status..."
docker compose -f docker-compose.prod.yml ps

# Test backend health
echo ""
echo "🏥 Testing backend health..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if curl -f http://localhost:8080/api/v1/health &> /dev/null; then
        print_success "Backend is healthy"
        break
    else
        attempt=$((attempt+1))
        echo "Attempt $attempt/$max_attempts..."
        sleep 2
    fi
done

if [ $attempt -eq $max_attempts ]; then
    print_error "Backend failed to start"
    echo "Checking logs:"
    docker compose -f docker-compose.prod.yml logs backend
    exit 1
fi

# Test frontend
echo ""
echo "🌐 Testing frontend..."
if curl -f http://localhost:80 &> /dev/null; then
    print_success "Frontend is healthy"
else
    print_warning "Frontend may not be ready yet"
fi

# Display deployment information
echo ""
echo "=================================="
print_success "Deployment Complete!"
echo "=================================="
echo ""
echo "📊 Service URLs:"
echo "   Frontend: http://localhost:80"
echo "   Backend:  http://localhost:8080"
echo ""
echo "📝 Useful commands:"
echo "   View logs:    docker compose -f docker-compose.prod.yml logs -f"
echo "   Stop:         docker compose -f docker-compose.prod.yml down"
echo "   Restart:      docker compose -f docker-compose.prod.yml restart"
echo "   Status:       docker compose -f docker-compose.prod.yml ps"
echo ""
echo "🔒 Next steps:"
echo "   1. Set up Nginx reverse proxy with SSL"
echo "   2. Configure firewall (ufw)"
echo "   3. Set up automated backups"
echo "   4. Configure monitoring"
echo ""
print_warning "Don't forget to set up SSL certificates for production!"
echo ""

