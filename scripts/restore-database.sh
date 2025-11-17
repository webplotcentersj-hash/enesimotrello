#!/bin/bash

# TaskBoard Database Restore Script

set -e

# Configuration
CONTAINER_NAME="${CONTAINER_NAME:-taskboard-postgres}"
DB_NAME="${DB_NAME:-taskboard}"
DB_USER="${DB_USER:-postgres}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

echo "♻️  TaskBoard Database Restore"
echo "=============================="
echo ""

# Check if backup file is provided
if [ -z "$1" ]; then
    print_error "Please provide backup file path"
    echo "Usage: $0 <backup-file>"
    echo ""
    echo "Available backups:"
    ls -lh backups/ 2>/dev/null || echo "No backups found"
    exit 1
fi

BACKUP_FILE="$1"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    print_error "Backup file not found: $BACKUP_FILE"
    exit 1
fi

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    print_error "Container $CONTAINER_NAME is not running"
    exit 1
fi

# Confirm restore
print_warning "This will OVERWRITE the current database!"
echo "Database: $DB_NAME"
echo "Backup file: $BACKUP_FILE"
echo ""
read -p "Are you sure you want to continue? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "Restore cancelled"
    exit 0
fi

# Create a backup of current database before restore
echo ""
echo "📦 Creating safety backup of current database..."
SAFETY_BACKUP="./backups/pre_restore_$(date +%Y%m%d_%H%M%S).sql"
mkdir -p backups
docker exec "$CONTAINER_NAME" pg_dump -U "$DB_USER" "$DB_NAME" > "$SAFETY_BACKUP"
print_success "Safety backup created: $SAFETY_BACKUP"

# Restore database
echo ""
echo "♻️  Restoring database..."

if [[ "$BACKUP_FILE" == *.gz ]]; then
    echo "Decompressing and restoring..."
    gunzip -c "$BACKUP_FILE" | docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" "$DB_NAME"
else
    docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" "$DB_NAME" < "$BACKUP_FILE"
fi

if [ $? -eq 0 ]; then
    print_success "Database restored successfully"
else
    print_error "Restore failed"
    print_warning "Your original database is backed up at: $SAFETY_BACKUP"
    exit 1
fi

# Restart backend to clear any cached connections
echo ""
echo "🔄 Restarting backend service..."
docker compose -f docker-compose.prod.yml restart backend
print_success "Backend restarted"

echo ""
echo "=============================="
print_success "Restore Complete!"
echo "=============================="
echo ""

