#!/bin/bash
# Generate CA and server certificates for testing
# This creates a custom CA that simulates a corporate proxy CA

set -e

CERTS_DIR="$(dirname "$0")/certs"
mkdir -p "$CERTS_DIR"

echo "=== Generating Custom CA ==="

# Generate CA private key
openssl genrsa -out "$CERTS_DIR/ca-key.pem" 4096

# Generate CA certificate
openssl req -new -x509 -days 365 -key "$CERTS_DIR/ca-key.pem" \
    -out "$CERTS_DIR/ca.pem" \
    -subj "/C=US/ST=Test/L=Test/O=Trustica Test CA/CN=Trustica Test Root CA"

echo "=== Generating Server Certificate ==="

# Generate server private key
openssl genrsa -out "$CERTS_DIR/server-key.pem" 2048

# Generate server CSR
openssl req -new -key "$CERTS_DIR/server-key.pem" \
    -out "$CERTS_DIR/server.csr" \
    -subj "/C=US/ST=Test/L=Test/O=Test Server/CN=testserver.local"

# Create extensions file for SAN
cat > "$CERTS_DIR/server-ext.cnf" << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = testserver.local
DNS.2 = localhost
IP.1 = 127.0.0.1
EOF

# Sign server certificate with CA
openssl x509 -req -days 365 \
    -in "$CERTS_DIR/server.csr" \
    -CA "$CERTS_DIR/ca.pem" \
    -CAkey "$CERTS_DIR/ca-key.pem" \
    -CAcreateserial \
    -out "$CERTS_DIR/server.pem" \
    -extfile "$CERTS_DIR/server-ext.cnf"

# Cleanup CSR and temp files
rm -f "$CERTS_DIR/server.csr" "$CERTS_DIR/server-ext.cnf" "$CERTS_DIR/ca.srl"

# Set permissions
chmod 644 "$CERTS_DIR"/*.pem
chmod 600 "$CERTS_DIR"/*-key.pem

echo ""
echo "=== Certificates Generated ==="
echo "CA Certificate:     $CERTS_DIR/ca.pem"
echo "CA Private Key:     $CERTS_DIR/ca-key.pem"
echo "Server Certificate: $CERTS_DIR/server.pem"
echo "Server Private Key: $CERTS_DIR/server-key.pem"
echo ""
echo "To verify: openssl verify -CAfile $CERTS_DIR/ca.pem $CERTS_DIR/server.pem"
