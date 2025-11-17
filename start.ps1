# Startup script for Task Board Application

Write-Host "Starting Task Board Application..." -ForegroundColor Green

# Check if Docker is running
try {
    docker version | Out-Null
    Write-Host "Docker is running" -ForegroundColor Green
} catch {
    Write-Host "Docker is not running. Please start Docker Desktop and try again." -ForegroundColor Red
    exit 1
}

# Start the application with Docker Compose
Write-Host "`nStarting services with Docker Compose..." -ForegroundColor Yellow
docker-compose up -d

# Wait for services to be ready
Write-Host "`nWaiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Check if services are running
Write-Host "`nChecking service status..." -ForegroundColor Yellow
docker-compose ps

Write-Host "`nApplication is starting up!" -ForegroundColor Green
Write-Host "Frontend: http://localhost:3000" -ForegroundColor Cyan
Write-Host "Backend API: http://localhost:8080/api/v1" -ForegroundColor Cyan
Write-Host "`nTo test the API, run: .\scripts\test-api.ps1" -ForegroundColor Yellow
Write-Host "To stop the application, run: docker-compose down" -ForegroundColor Yellow
