#!/bin/bash
# Test that tools FAIL to connect without trusting the CA
# This proves the problem exists before Trustica fixes it

set -e

SCRIPT_DIR="$(dirname "$0")"
source "$SCRIPT_DIR/common.sh"

echo "╔═══════════════════════════════════════════════════════════════════╗"
echo "║  BEFORE TRUSTICA - Verifying SSL failures                         ║"
echo "╚═══════════════════════════════════════════════════════════════════╝"
echo ""

FAILURES=0
EXPECTED_FAILURES=4

# Test 1: curl
echo "┌─────────────────────────────────────────────────────────────────────┐"
echo "│ Test 1: curl (should FAIL)                                          │"
echo "└─────────────────────────────────────────────────────────────────────┘"
if curl --silent --max-time 5 "$TEST_SERVER_URL/health" > /dev/null 2>&1; then
    echo "  ❌ UNEXPECTED: curl succeeded (should have failed)"
else
    echo "  ✅ EXPECTED: curl failed with SSL error"
    ((FAILURES++))
fi
echo ""

# Test 2: git
echo "┌─────────────────────────────────────────────────────────────────────┐"
echo "│ Test 2: git (should FAIL)                                           │"
echo "└─────────────────────────────────────────────────────────────────────┘"
if git ls-remote "$TEST_SERVER_URL/repo.git" > /dev/null 2>&1; then
    echo "  ❌ UNEXPECTED: git succeeded (should have failed)"
else
    echo "  ✅ EXPECTED: git failed with SSL error"
    ((FAILURES++))
fi
echo ""

# Test 3: pip
echo "┌─────────────────────────────────────────────────────────────────────┐"
echo "│ Test 3: pip (should FAIL)                                           │"
echo "└─────────────────────────────────────────────────────────────────────┘"
if pip3 install --quiet --index-url "$TEST_SERVER_URL/simple/" --dry-run nonexistent-package > /dev/null 2>&1; then
    echo "  ❌ UNEXPECTED: pip succeeded (should have failed)"
else
    echo "  ✅ EXPECTED: pip failed with SSL error"
    ((FAILURES++))
fi
echo ""

# Test 4: Node.js
echo "┌─────────────────────────────────────────────────────────────────────┐"
echo "│ Test 4: node (should FAIL)                                          │"
echo "└─────────────────────────────────────────────────────────────────────┘"
if node -e "fetch('$TEST_SERVER_URL/health').then(() => process.exit(0)).catch(() => process.exit(1))" > /dev/null 2>&1; then
    echo "  ❌ UNEXPECTED: node succeeded (should have failed)"
else
    echo "  ✅ EXPECTED: node failed with SSL error"
    ((FAILURES++))
fi
echo ""

# Summary
echo "═══════════════════════════════════════════════════════════════════════"
if [ "$FAILURES" -eq "$EXPECTED_FAILURES" ]; then
    echo "✅ All $EXPECTED_FAILURES tools failed as expected (SSL not trusted)"
    echo "   This proves the problem exists. Trustica should fix this."
    exit 0
else
    echo "❌ Expected $EXPECTED_FAILURES failures, got $FAILURES"
    echo "   Some tools didn't fail as expected."
    exit 1
fi
