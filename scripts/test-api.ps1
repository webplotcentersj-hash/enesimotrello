# Test script for the Task Board API

$API_URL = "http://localhost:8080/api/v1"

Write-Host "Testing Task Board API..." -ForegroundColor Green

# Test user registration
Write-Host "`n1. Testing user registration..." -ForegroundColor Yellow
$registerBody = @{
    email = "test@example.com"
    username = "testuser"
    password = "password123"
    first_name = "Test"
    last_name = "User"
} | ConvertTo-Json

try {
    $registerResponse = Invoke-RestMethod -Uri "$API_URL/auth/register" -Method POST -Body $registerBody -ContentType "application/json"
    Write-Host "Register response: $($registerResponse | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "Register failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test user login
Write-Host "`n2. Testing user login..." -ForegroundColor Yellow
$loginBody = @{
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$API_URL/auth/login" -Method POST -Body $loginBody -ContentType "application/json"
    Write-Host "Login response: $($loginResponse | ConvertTo-Json)" -ForegroundColor Green
    
    $token = $loginResponse.token
    if ($token) {
        Write-Host "Token extracted: $($token.Substring(0, 20))..." -ForegroundColor Green
        
        # Test getting boards
        Write-Host "`n3. Testing get boards..." -ForegroundColor Yellow
        $headers = @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        }
        
        try {
            $boardsResponse = Invoke-RestMethod -Uri "$API_URL/boards" -Method GET -Headers $headers
            Write-Host "Boards response: $($boardsResponse | ConvertTo-Json)" -ForegroundColor Green
        } catch {
            Write-Host "Get boards failed: $($_.Exception.Message)" -ForegroundColor Red
        }
        
        # Test creating a board
        Write-Host "`n4. Testing create board..." -ForegroundColor Yellow
        $boardBody = @{
            title = "Test Board"
            description = "A test board"
        } | ConvertTo-Json
        
        try {
            $boardResponse = Invoke-RestMethod -Uri "$API_URL/boards" -Method POST -Body $boardBody -Headers $headers
            Write-Host "Board creation response: $($boardResponse | ConvertTo-Json)" -ForegroundColor Green
            
            $boardId = $boardResponse.board.id
            if ($boardId) {
                Write-Host "Board ID: $boardId" -ForegroundColor Green
                
                # Test creating a task
                Write-Host "`n5. Testing create task..." -ForegroundColor Yellow
                $taskBody = @{
                    title = "Test Task"
                    description = "A test task"
                    priority = "high"
                } | ConvertTo-Json
                
                try {
                    $taskResponse = Invoke-RestMethod -Uri "$API_URL/boards/$boardId/tasks" -Method POST -Body $taskBody -Headers $headers
                    Write-Host "Task creation response: $($taskResponse | ConvertTo-Json)" -ForegroundColor Green
                } catch {
                    Write-Host "Create task failed: $($_.Exception.Message)" -ForegroundColor Red
                }
            }
        } catch {
            Write-Host "Create board failed: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
} catch {
    Write-Host "Login failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nAPI testing completed!" -ForegroundColor Green
