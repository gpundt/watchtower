#!/bin/bash
source ./_helpers.sh

CERTS_ROOT_DIR="../certs"
CA_CERT="${CERTS_ROOT_DIR}/ca/ca.crt"
CA_KEY="${CERTS_ROOT_DIR}/ca/ca.key"
SERVER_CERT="${CERTS_ROOT_DIR}/server/server.crt"
SERVER_KEY="${CERTS_ROOT_DIR}/server/server.key"

function apt_install_openssl() {
    if sudo dpkg -s openssl >/dev/null 2>&1; then
        return
    fi
    start_step_message "Installing OpenSSL APT Package"
    if ! sudo apt install openssl; then
        error_message "Failed to 'sudo apt install openssl'"
    fi
    successful
}

function prepare_certs_directory() {
    start_step_message "Preparing Certificate Output Directory '${CERTS_DIR}'"
    mkdir -p "${CERTS_ROOT_DIR}/ca"
    mkdir -p "${CERTS_ROOT_DIR}/server"
    mkdir -p "${CERTS_ROOT_DIR}/agents"
    successful
}

function generate_ca() {
    start_step_message "Generating CA Cert and Key '${CERTS_ROOT_DIR}/ca'"
    if [ -f "${CA_CERT}" ]; then
        error_message "CA Cert already exists: '${CA_CERT}'"
    fi

    if [ -f "${CA_KEY}" ]; then
        error_message "CA Key already exists: '${CA_KEY}'"
    fi

    if ! openssl req -x509 -sha256 -nodes -days 3650 -newkey rsa:4096 \
  -keyout "${CA_KEY}" -out "${CA_CERT}" \
  -subj "/C=US/ST=HI/L=Honolulu/O=WatchtowerCA/OU=CertificateAuthority" \
  -addext "basicConstraints=critical,CA:TRUE"; then
        error_message "Failed to generate CA Cert and CA Key"
    fi
    successful
}

function generate_server_key() {
    start_step_message "Generating Server Private Key '${CERTS_ROOT_DIR}/server'"
    if [ -f "${SERVER_KEY}" ]; then
        error_message "Server Key already exists: '${SERVER_KEY}'"
    fi

    if ! openssl genrsa -out "${SERVER_KEY}" 2048; then
        error_message "Failed to generate Server Key"
    fi
    successful
}

function generate_server_cert() {
    start_step_message "Generating Server Certificate Signing Request (CSR)"
    if [ -f "${CERTS_ROOT_DIR}/server/server.csr" ]; then
        error_message "Server CSR already exists: '${CERTS_ROOT_DIR}/server/server.csr'"
    fi

    if ! openssl req -new -key "${SERVER_KEY}" -out "${CERTS_ROOT_DIR}/server/server.csr" \
     -subj "/C=US/ST=HI/L=Honolulu/O=Watchtower/OU=Server/"; then
        error_message "Failed to generate Server CSR"
    fi
    cat > /tmp/server_ext.cnf << EOF
subjectAltName = DNS:localhost, DNS:watchtower.home, IP:127.0.0.1, IP:0.0.0.0
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
EOF

    successful

    start_step_message "Generating CA Signed Server Certificate '${CERTS_ROOT_DIR}/server/'"
    if [ -f "${SERVER_CERT}" ]; then
        error_message "Server Cert already exists: '${SERVER_CERT}'"
    fi

    if ! openssl x509 -req -in "${CERTS_ROOT_DIR}/server/server.csr" \
     -CA "${CA_CERT}" -CAkey "${CA_KEY}" -CAcreateserial \
     -out "${SERVER_CERT}" -days 365 -sha256 \
     -extfile /tmp/server_ext.cnf; then
        error_message "Failed to generate CA Signed server cert: '${SERVER_CERT}'"
    fi

    rm -rf "${CERTS_ROOT_DIR}/server/server.csr"
    rm -f /tmp/server_ext.cnf
    successful
}

function main() {
    apt_install_openssl
    prepare_certs_directory
    generate_ca
    generate_server_key
    generate_server_cert
}

main