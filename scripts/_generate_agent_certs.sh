#!/bin/bash
source ./_helpers.sh

if [ "$#" -eq 0 ]; then
    echo "Error: No agent hostname provided"
    echo "Usage: $0 [agent_hostname]"
    exit 1
fi

AGENT_HOSTNAME=$1
echo "AGENT_HOSTNAME=${AGENT_HOSTNAME}" > ../.env

CERTS_ROOT_DIR="../certs"
AGENT_CERTS_DIR="${CERTS_ROOT_DIR}/${AGENT_HOSTNAME}_agent"
CA_CERT="${CERTS_ROOT_DIR}/ca/ca.crt"
CA_KEY="${CERTS_ROOT_DIR}/ca/ca.key"
AGENT_CERT="${AGENT_CERTS_DIR}/agent.crt"
AGENT_KEY="${AGENT_CERTS_DIR}/agent.key"
AGENT_YAML_CONFIG="../config/agent.yaml"
DST_CA_CERT="/opt/watchtower/tls/ca/ca.crt"
DST_AGENT_CERT="/opt/watchtower/tls/${AGENT_HOSTNAME}_agent/agent.crt"
DST_AGENT_KEY="/opt/watchtower/tls/${AGENT_HOSTNAME}_agent/agent.key"


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
    start_step_message "Preparing Certificate Output Directory '${AGENT_CERTS_DIR}'"
    mkdir -p $AGENT_CERTS_DIR
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
    if [ -f "${CERTS_ROOT_DIR}/${AGENT_HOSTNAME}_agent/agent.csr" ]; then
        error_message "Agent CSR already exists: '${CERTS_ROOT_DIR}/${AGENT_HOSTNAME}_agent/agent.csr'"
    fi

    if ! openssl req -new -key "${AGENT_KEY}" -out "${CERTS_ROOT_DIR}/${AGENT_HOSTNAME}_agent/agent.csr" \
     -subj "/C=US/ST=HI/L=Honolulu/O=WatchtowerAgent/OU=Agent/CN=${AGENT_HOSTNAME}"; then
        error_message "Failed to generate Agent CSR"
    fi
     cat > /tmp/agent_ext.cnf << EOF
subjectAltName = DNS:localhost, DNS:${AGENT_HOSTNAME}, IP:127.0.0.1, IP:0.0.0.0
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth
EOF
    successful

    start_step_message "Generating CA Signed Agent Certificate '${CERTS_ROOT_DIR}/${AGENT_HOSTNAME}_agent/'"
    if [ -f "${AGENT_CERT}" ]; then
        error_message "Agent Cert already exists: ${AGENT_CERT}"
    fi

    if ! openssl x509 -req -in "${CERTS_ROOT_DIR}/${AGENT_HOSTNAME}_agent/agent.csr" \
     -CA "${CA_CERT}" -CAkey "${CA_KEY}" -CAcreateserial \
     -out "${AGENT_CERT}" -days 365 -sha256 \
     -extfile /tmp/agent_ext.cnf; then
        error_message "Failed to generate CA Signed agent cert: '${AGENT_CERT}'"
    fi

    rm -rf "${CERTS_ROOT_DIR}/${AGENT_HOSTNAME}_agent/agent.csr"
    rm -f /tmp/agent_ext.cnf
    successful
}

function modify_agent_yaml_config() {
    start_step_message "Modifying '${AGENT_YAML_CONFIG}'"
    
    if [ ! -f "${AGENT_YAML_CONFIG}" ]; then
        error_message "Agent YAML config not found : '${AGENT_YAML_CONFIG}'"
    fi

    if ! sed -i \
        -e "s|ca_cert:.*|ca_cert: \"${DST_CA_CERT}\"|" \
        -e "s|agent_cert:.*|agent_cert: \"${DST_AGENT_CERT}\"|" \
        -e "s|agent_key:.*|agent_key: \"${DST_AGENT_KEY}\"|" \
        "${AGENT_YAML_CONFIG}"; then
        error_message "Failed to modify '${AGENT_YAML_CONFIG}'"
    fi

    echo ""
    grep -A4 "^tls:" "${AGENT_YAML_CONFIG}"
    echo ""

    successful
}

function main() {
    apt_install_openssl
    prepare_certs_directory
    generate_agent_key
    generate_agent_cert
    modify_agent_yaml_config
}

main