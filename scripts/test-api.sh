#!/bin/bash

# qkrn API test script

BASE_URL="http://localhost:8080"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}qkrn API Test Script${NC}"
echo "Testing API at $BASE_URL"
echo

# Function to run a test
run_test() {
    local description="$1"
    local command="$2"
    local expected_status="$3"
    
    echo -e "${YELLOW}Test:${NC} $description"
    echo "Command: $command"
    
    response=$(eval $command)
    status=$?
    
    if [ $status -eq 0 ]; then
        echo -e "${GREEN}✓ Success${NC}"
        echo "Response: $response"
    else
        echo -e "${RED}✗ Failed${NC}"
        echo "Response: $response"
    fi
    echo
}

# Test service info
run_test "Service Information" "curl -s $BASE_URL/"

# Test health check
run_test "Health Check" "curl -s $BASE_URL/health"

# Test storing a key-value pair
run_test "Store key-value" "curl -s -X PUT $BASE_URL/kv/test1 -H 'Content-Type: application/json' -d '{\"value\":\"hello\"}'"

# Test retrieving a value
run_test "Retrieve value" "curl -s $BASE_URL/kv/test1"

# Test storing another key-value pair
run_test "Store another key" "curl -s -X PUT $BASE_URL/kv/test2 -H 'Content-Type: application/json' -d '{\"value\":\"world\"}'"

# Test listing keys
run_test "List all keys" "curl -s $BASE_URL/keys"

# Test deleting a key
run_test "Delete key" "curl -s -X DELETE $BASE_URL/kv/test1"

# Test retrieving deleted key (should fail)
run_test "Retrieve deleted key (should fail)" "curl -s $BASE_URL/kv/test1"

# Test listing keys after deletion
run_test "List keys after deletion" "curl -s $BASE_URL/keys"

echo -e "${GREEN}Test script completed!${NC}"
