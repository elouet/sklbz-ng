#!/bin/bash

# Script to generate test certificates for mTLS
# This script creates:
# - CA certificate and key
# - Server certificate and key
# - Client certificate and key

set -e

echo "Generating mTLS test certificates..."

# Create certificates directory
mkdir -p certs
cd certs

# 1. Generate CA private key and certificate
echo "Generating CA certificate..."
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt -subj "/CN=sklbz-ng CA"

# 2. Generate server private key and CSR
echo "Generating server certificate..."
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/CN=localhost"

# 3. Sign server certificate with CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256

# 4. Generate client private key and CSR
echo "Generating client certificate..."
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "/CN=test-client"

# 5. Sign client certificate with CA
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365 -sha256

# 6. Create client PKCS#12 file (optional, for browser testing)
echo "Creating client PKCS#12 file..."
openssl pkcs12 -export -out client.p12 -inkey client.key -in client.crt -certfile ca.crt -passout pass:test123

# 7. Clean up CSR files
rm -f *.csr

# 8. Create a combined client certificate+key file for curl testing
cat client.crt client.key > client-combined.pem

echo ""
echo "✅ Certificates generated successfully in the 'certs/' directory:"
echo ""
echo "  ca.crt          - CA certificate"
echo "  ca.key          - CA private key"
echo "  server.crt      - Server certificate"
echo "  server.key      - Server private key"
echo "  client.crt      - Client certificate"
echo "  client.key      - Client private key"
echo "  client.p12      - Client PKCS#12 file (password: test123)"
echo "  client-combined.pem - Combined client cert+key for curl"
echo ""
echo "To start the server with mTLS:"
echo "  ./sklbz-ng -mtls -cert certs/server.crt -key certs/server.key -ca-cert certs/ca.crt"
echo ""
echo "To test with curl:"
echo "  curl --cert certs/client-combined.pem --cacert certs/ca.crt https://localhost:8080/api/articles"
