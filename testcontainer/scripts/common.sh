#!/bin/bash
# Common variables and functions for test scripts

# Test server URL (can be overridden by environment)
export TEST_SERVER_URL="${TEST_SERVER_URL:-https://testserver.local:8443}"

# CA certificate path
export CA_CERT_PATH="${CA_CERT_PATH:-/certs/ca.pem}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_failure() {
    echo -e "${RED}❌ $1${NC}"
}

log_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Wait for test server to be ready
wait_for_server() {
    local max_attempts=30
    local attempt=1
    
    echo "Waiting for test server at $TEST_SERVER_URL..."
    
    while [ $attempt -le $max_attempts ]; do
        # Use --insecure just to check if server is up
        if curl --silent --insecure --max-time 2 "$TEST_SERVER_URL/health" > /dev/null 2>&1; then
            echo "Test server is ready!"
            return 0
        fi
        echo "  Attempt $attempt/$max_attempts..."
        sleep 1
        ((attempt++))
    done
    
    echo "ERROR: Test server not ready after $max_attempts attempts"
    return 1
}
