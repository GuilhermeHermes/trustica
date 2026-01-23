#!/bin/bash
# Test that tools SUCCEED after Trustica applies the CA
# This proves Trustica fixed the problem

set -e

SCRIPT_DIR="$(dirname "$0")"
source "$SCRIPT_DIR/common.sh"

echo "╔═══════════════════════════════════════════════════════════════════╗"
echo "║  AFTER TRUSTICA - Verifying SSL works                             ║"
echo "╚═══════════════════════════════════════════════════════════════════╝"
echo ""

SUCCESSES=0
EXPECTED_SUCCESSES=4

# Test 1: curl
echo "┌─────────────────────────────────────────────────────────────────────┐"
echo "│ Test 1: curl (should SUCCEED)                                       │"
echo "└─────────────────────────────────────────────────────────────────────┘"
if curl --silent --max-time 5 "$TEST_SERVER_URL/health" > /dev/null 2>&1; then
    echo "  ✅ SUCCESS: curl works with trusted CA"
    ((SUCCESSES++))
else
    echo "  ❌ FAILED: curl still failing"
fi
echo ""

# Test 2: git
echo "┌─────────────────────────────────────────────────────────────────────┐"
echo "│ Test 2: git (should SUCCEED)                                        │"
echo "└─────────────────────────────────────────────────────────────────────┘"
if git ls-remote "$TEST_SERVER_URL/repo.git" > /dev/null 2>&1; then
    echo "  ✅ SUCCESS: git works with trusted CA"
    ((SUCCESSES++))
else
    echo "  ❌ FAILED: git still failing"
fi
echo ""

# Test 3: pip
echo "┌─────────────────────────────────────────────────────────────────────┐"
echo "│ Test 3: pip (should SUCCEED)                                        │"
echo "└─────────────────────────────────────────────────────────────────────┘"
# pip will fail to find the package, but SSL should work
# We check for SSL error specifically
pip_output=$(pip3 install --quiet --index-url "$TEST_SERVER_URL/simple/" --dry-run nonexistent-package 2>&1 || true)
if echo "$pip_output" | grep -qi "ssl\|certificate"; then
    echo "  ❌ FAILED: pip still has SSL errors"
else
    echo "  ✅ SUCCESS: pip connects without SSL errors"
    ((SUCCESSES++))
fi
echo ""

# Test 4: Node.js
echo "┌─────────────────────────────────────────────────────────────────────┐"
echo "│ Test 4: node (should SUCCEED)                                       │"
echo "└─────────────────────────────────────────────────────────────────────┘"
if node -e "fetch('$TEST_SERVER_URL/health').then(() => process.exit(0)).catch(() => process.exit(1))" > /dev/null 2>&1; then
    echo "  ✅ SUCCESS: node works with trusted CA"
    ((SUCCESSES++))
else
    echo "  ❌ FAILED: node still failing"
fi
echo ""

# Summary
echo "═══════════════════════════════════════════════════════════════════════"
if [ "$SUCCESSES" -eq "$EXPECTED_SUCCESSES" ]; then
    echo "✅ All $EXPECTED_SUCCESSES tools working! Trustica fixed SSL trust."
    exit 0
else
    echo "❌ Only $SUCCESSES/$EXPECTED_SUCCESSES tools working."
    echo "   Some adapters may need fixing."
    exit 1
fi
