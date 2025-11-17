# PowerShell script to push TaskBoard to GitHub
# Usage: .\push-to-github.ps1 YOUR_GITHUB_USERNAME

param(
    [Parameter(Mandatory=$true)]
    [string]$GitHubUsername
)

Write-Host "=======================================" -ForegroundColor Cyan
Write-Host "  Pushing TaskBoard to GitHub" -ForegroundColor Cyan
Write-Host "=======================================" -ForegroundColor Cyan
Write-Host ""

# Check if git is initialized
if (-not (Test-Path ".git")) {
    Write-Host "ERROR: Git repository not initialized!" -ForegroundColor Red
    exit 1
}

# Add remote
Write-Host "Step 1: Adding GitHub remote..." -ForegroundColor Yellow
$remoteUrl = "https://github.com/$GitHubUsername/task-board.git"
Write-Host "Remote URL: $remoteUrl" -ForegroundColor Gray

git remote add origin $remoteUrl 2>&1 | Out-Null

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Remote added successfully!" -ForegroundColor Green
} else {
    Write-Host "Note: Remote might already exist, continuing..." -ForegroundColor Yellow
    git remote set-url origin $remoteUrl
}

Write-Host ""

# Verify remote
Write-Host "Step 2: Verifying remote..." -ForegroundColor Yellow
git remote -v
Write-Host ""

# Push to GitHub
Write-Host "Step 3: Pushing to GitHub..." -ForegroundColor Yellow
Write-Host "This may take a moment..." -ForegroundColor Gray
Write-Host ""

git push -u origin main

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "=======================================" -ForegroundColor Green
    Write-Host "  ✓ Successfully pushed to GitHub!" -ForegroundColor Green
    Write-Host "=======================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Your repository is now live at:" -ForegroundColor Cyan
    Write-Host "https://github.com/$GitHubUsername/task-board" -ForegroundColor White
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "1. Add repository topics/tags on GitHub" -ForegroundColor White
    Write-Host "2. Update README.md with your personal info" -ForegroundColor White
    Write-Host "3. Star your own repository" -ForegroundColor White
    Write-Host "4. Share on LinkedIn!" -ForegroundColor White
    Write-Host ""
} else {
    Write-Host ""
    Write-Host "ERROR: Failed to push to GitHub!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Common issues:" -ForegroundColor Yellow
    Write-Host "1. Repository doesn't exist on GitHub yet" -ForegroundColor White
    Write-Host "   Solution: Create it at https://github.com/new" -ForegroundColor Gray
    Write-Host ""
    Write-Host "2. Authentication failed" -ForegroundColor White
    Write-Host "   Solution: Configure Git credentials or use SSH" -ForegroundColor Gray
    Write-Host ""
    Write-Host "3. Remote already exists with different URL" -ForegroundColor White
    Write-Host "   Solution: Run 'git remote remove origin' and try again" -ForegroundColor Gray
    Write-Host ""
    exit 1
}

