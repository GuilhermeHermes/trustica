#!/bin/bash
# Full integration test suite
# 1. Verify tools FAIL without CA trust
# 2. Run Trustica to apply CA
# 3. Verify tools SUCCEED after Trustica

set -e

SCRIPT_DIR="$(dirname "$0")"
source "$SCRIPT_DIR/common.sh"

echo ""
echo "╔═══════════════════════════════════════════════════════════════════╗"
echo "║                                                                   ║"
echo "║             TRUSTICA INTEGRATION TEST SUITE                       ║"
echo "║                                                                   ║"
echo "╚═══════════════════════════════════════════════════════════════════╝"
echo ""

# Wait for test server
wait_for_server

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  PHASE 1: Verify SSL failures (before Trustica)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if ! "$SCRIPT_DIR/test_before.sh"; then
    echo ""
    echo "ERROR: Pre-test failed. Tools should be failing without CA trust."
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  PHASE 2: Run Trustica"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

echo "Running: trustica apply --cert $CA_CERT_PATH"
echo ""

if ! trustica apply --cert "$CA_CERT_PATH"; then
    echo ""
    echo "ERROR: Trustica failed to apply certificate"
    exit 1
fi

echo ""
echo "Trustica completed successfully!"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  PHASE 3: Verify SSL works (after Trustica)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if ! "$SCRIPT_DIR/test_after.sh"; then
    echo ""
    echo "ERROR: Post-test failed. Some tools still not working."
    exit 1
fi

echo ""
echo "╔═══════════════════════════════════════════════════════════════════╗"
echo "║                                                                   ║"
echo "║  ✅ ALL TESTS PASSED - Trustica is working correctly!             ║"
echo "║                                                                   ║"
echo "╚═══════════════════════════════════════════════════════════════════╝"
echo ""
