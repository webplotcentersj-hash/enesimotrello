#!/bin/bash

# Test script for the Task Board API

API_URL="http://localhost:8080/api/v1"

echo "Testing Task Board API..."

# Test health endpoint (if available)
echo "1. Testing API health..."
curl -s -o /dev/null -w "%{http_code}" "$API_URL/health" || echo "Health endpoint not available"

# Test user registration
echo -e "\n2. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }')

echo "Register response: $REGISTER_RESPONSE"

# Test user login
echo -e "\n3. Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')

echo "Login response: $LOGIN_RESPONSE"

# Extract token from login response
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
  echo "Token extracted: ${TOKEN:0:20}..."
  
  # Test getting boards
  echo -e "\n4. Testing get boards..."
  curl -s -X GET "$API_URL/boards" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" | jq '.' || echo "Failed to get boards"
  
  # Test creating a board
  echo -e "\n5. Testing create board..."
  BOARD_RESPONSE=$(curl -s -X POST "$API_URL/boards" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "title": "Test Board",
      "description": "A test board"
    }')
  
  echo "Board creation response: $BOARD_RESPONSE"
  
  # Extract board ID
  BOARD_ID=$(echo $BOARD_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
  
  if [ -n "$BOARD_ID" ]; then
    echo "Board ID: $BOARD_ID"
    
    # Test creating a task
    echo -e "\n6. Testing create task..."
    TASK_RESPONSE=$(curl -s -X POST "$API_URL/boards/$BOARD_ID/tasks" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "title": "Test Task",
        "description": "A test task",
        "priority": "high"
      }')
    
    echo "Task creation response: $TASK_RESPONSE"
  fi
else
  echo "Failed to extract token from login response"
fi

echo -e "\nAPI testing completed!"
