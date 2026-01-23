# Trustica Integration Test Container

This directory contains the integration test setup for Trustica.

## What It Does

1. **Creates a custom CA** - Simulates a corporate proxy CA
2. **Runs an HTTPS server** - Uses a cert signed by the custom CA
3. **Tests tools WITHOUT trust** - Proves they fail (curl, git, pip, node)
4. **Runs Trustica** - Applies the CA certificate
5. **Tests tools WITH trust** - Proves they now work

## Quick Start

```bash
# Generate certificates, build containers, run tests
make integration-test
```

## Individual Commands

```bash
# Generate test certificates
make integration-certs

# Build test containers
make integration-build

# Run full test suite
make integration-test

# Open shell in test container (for debugging)
make integration-shell

# Start test server locally (for development)
make integration-server

# Clean up
make integration-clean
```

## Directory Structure

```
testcontainer/
├── generate_certs.sh       # Creates CA + server certificates
├── Dockerfile              # Multi-stage build for test container
├── docker-compose.yml      # Orchestrates test server + test runner
├── testserver/
│   └── main.go             # Simple HTTPS server
├── scripts/
│   ├── common.sh           # Shared variables and functions
│   ├── test_before.sh      # Verify tools FAIL without CA
│   ├── test_after.sh       # Verify tools SUCCEED with CA
│   └── run_all_tests.sh    # Full test suite
└── certs/                  # Generated certificates (gitignored)
    ├── ca.pem              # Custom CA certificate
    ├── ca-key.pem          # CA private key
    ├── server.pem          # Server certificate
    └── server-key.pem      # Server private key
```

## Test Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    INTEGRATION TEST FLOW                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. BEFORE TRUSTICA                                             │
│     curl https://testserver.local    → ❌ SSL error             │
│     git ls-remote https://...        → ❌ SSL error             │
│     pip install --index-url ...      → ❌ SSL error             │
│     node fetch(...)                  → ❌ SSL error             │
│                                                                 │
│  2. RUN TRUSTICA                                                │
│     trustica apply --cert /certs/ca.pem                         │
│                                                                 │
│  3. AFTER TRUSTICA                                              │
│     curl https://testserver.local    → ✅ works                 │
│     git ls-remote https://...        → ✅ works                 │
│     pip install --index-url ...      → ✅ works                 │
│     node fetch(...)                  → ✅ works                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```
