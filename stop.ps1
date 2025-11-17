# Stop script for Task Board Application

Write-Host "Stopping Task Board Application..." -ForegroundColor Green

# Stop and remove containers
docker-compose down

Write-Host "Application stopped!" -ForegroundColor Green
