#!/bin/bash
source ./_helpers.sh

if [ "$#" -eq 0 ]; then
    echo "Error: No agent hostname provided"
    echo "Usage: $0 [agent_hostname]"
    exit 1
fi

AGENT_HOSTNAME=$1

CERTS_ROOT_DIR="../certs"
CA_CERT="${CERTS_ROOT_DIR}/ca/ca.crt"
CA_KEY="${CERTS_ROOT_DIR}/ca/ca.key"
AGENT_CERT="${CERTS_ROOT_DIR}/agents/${AGENT_HOSTNAME}_agent.crt"
AGENT_KEY="${CERTS_ROOT_DIR}/agents/${AGENT_HOSTNAME}_agent.key"

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
    mkdir -p "${CERTS_ROOT_DIR}/agents"
    successful
}

function generate_agent_key() {
    start_step_message "Generating Agent Private Key '${CERTS_ROOT_DIR}/agents'"
    if [ -f "${AGENT_KEY}" ]; then
        error_message "Agent Key already exists: '${AGENT_KEY}'"
    fi

    if ! openssl genrsa -out "${AGENT_KEY}" 2048; then
        error_message "Failed to generate Agent Key"
    fi
    successful
}

function generate_agent_cert() {
    start_step_message "Generating Agent Certificate Signing Request (CSR)"
    if [ -f "${CERTS_ROOT_DIR}/agents/${AGENT_HOSTNAME}_agent.csr" ]; then
        error_message "Agent CSR already exists: '${CERTS_ROOT_DIR}/agents/${AGENT_HOSTNAME}_agent.csr'"
    fi

    if ! openssl req -new -key "${AGENT_KEY}" -out "${CERTS_ROOT_DIR}/agents/${AGENT_HOSTNAME}_agent.csr" \
     -subj "/C=US/ST=HI/L=Honolulu/O=WatchtowerAgent/OU=Agent/CN=${AGENT_HOSTNAME}"; then
        error_message "Failed to generate Agent CSR"
    fi
     cat > /tmp/agent_ext.cnf << EOF
subjectAltName = DNS:localhost, DNS:${AGENT_HOSTNAME}, IP:127.0.0.1, IP:0.0.0.0
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth
EOF
    successful

    start_step_message "Generating CA Signed Agent Certificate '${CERTS_ROOT_DIR}/agents/'"
    if [ -f "${AGENT_CERT}" ]; then
        error_message "Agent Cert already exists: ${AGENT_CERT}"
    fi

    if ! openssl x509 -req -in "${CERTS_ROOT_DIR}/agents/${AGENT_HOSTNAME}_agent.csr" \
     -CA "${CA_CERT}" -CAkey "${CA_KEY}" -CAcreateserial \
     -out "${AGENT_CERT}" -days 365 -sha256 \
     -extfile /tmp/agent_ext.cnf; then
        error_message "Failed to generate CA Signed agent cert: '${AGENT_CERT}'"
    fi

    rm -rf "${CERTS_ROOT_DIR}/agents/${AGENT_HOSTNAME}_agent.csr"
    rm -f /tmp/agent_ext.cnf
    successful
}

function main() {
    apt_install_openssl
    prepare_certs_directory
    generate_agent_key
    generate_agent_cert
}

main