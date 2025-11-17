#!/bin/bash

# TaskBoard Database Backup Script
# Run this script regularly to backup your PostgreSQL database

set -e

# Configuration
BACKUP_DIR="${BACKUP_DIR:-./backups}"
CONTAINER_NAME="${CONTAINER_NAME:-taskboard-postgres}"
DB_NAME="${DB_NAME:-taskboard}"
DB_USER="${DB_USER:-postgres}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

echo "💾 TaskBoard Database Backup"
echo "============================"
echo ""

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Generate backup filename with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/taskboard_${TIMESTAMP}.sql"
BACKUP_FILE_GZ="$BACKUP_FILE.gz"

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    print_error "Container $CONTAINER_NAME is not running"
    exit 1
fi

# Create backup
echo "📦 Creating backup..."
if docker exec "$CONTAINER_NAME" pg_dump -U "$DB_USER" "$DB_NAME" > "$BACKUP_FILE"; then
    print_success "Backup created: $BACKUP_FILE"
else
    print_error "Backup failed"
    exit 1
fi

# Compress backup
echo "🗜️  Compressing backup..."
gzip "$BACKUP_FILE"
print_success "Backup compressed: $BACKUP_FILE_GZ"

# Calculate backup size
BACKUP_SIZE=$(du -h "$BACKUP_FILE_GZ" | cut -f1)
echo "📊 Backup size: $BACKUP_SIZE"

# Remove old backups
echo ""
echo "🧹 Removing backups older than $RETENTION_DAYS days..."
DELETED_COUNT=$(find "$BACKUP_DIR" -name "taskboard_*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete -print | wc -l)
if [ "$DELETED_COUNT" -gt 0 ]; then
    print_success "Removed $DELETED_COUNT old backup(s)"
else
    echo "No old backups to remove"
fi

# Display backup info
echo ""
echo "=============================="
print_success "Backup Complete!"
echo "=============================="
echo ""
echo "📁 Backup location: $BACKUP_FILE_GZ"
echo "📊 Backup size: $BACKUP_SIZE"
echo ""
echo "💡 To restore this backup, run:"
echo "   gunzip -c $BACKUP_FILE_GZ | docker exec -i $CONTAINER_NAME psql -U $DB_USER $DB_NAME"
echo ""

# Optional: Upload to cloud storage
if [ -n "$AWS_S3_BUCKET" ]; then
    echo "☁️  Uploading to S3..."
    if command -v aws &> /dev/null; then
        aws s3 cp "$BACKUP_FILE_GZ" "s3://$AWS_S3_BUCKET/backups/"
        print_success "Uploaded to S3"
    else
        print_warning "AWS CLI not found, skipping S3 upload"
    fi
fi

