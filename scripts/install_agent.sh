#!/bin/bash
source ./_helpers.sh

if [ "$#" -eq 0 ]; then
    echo "Error: No agent hostname provided"
    echo "Usage: $0 [agent_hostname]"
    exit 1
fi

AGENT_HOSTNAME=$1

### Important Directories ###
ETC_DIRECTORY=/etc/watchtower
OPT_DIRECTORY=/opt/watchtower
BINARY_DIRECTORY=$OPT_DIRECTORY/bin
TLS_DIRECTORY=$OPT_DIRECTORY/tls
LOG_DIRECTORY=/var/log/watchtower

### Important Filepaths ###
SRC_AGENT_CONFIG=../config/agent.yaml
DST_AGENT_CONFIG=$ETC_DIRECTORY/agent.yaml
SRC_AGENT_BINARY=../build/watchtower_agent
DST_AGENT_BINARY=$BINARY_DIRECTORY/watchtower_agent
SRC_CA_CERT=../certs/ca/ca.crt
DST_CA_CERT=$TLS_DIRECTORY/ca/ca.crt
SRC_AGENT_CERT=../certs/agents/"$AGENT_HOSTNAME"_agent.crt
DST_AGENT_CERT=$TLS_DIRECTORY/agent/"$AGENT_HOSTNAME"_agent.crt
SRC_AGENT_KEY=../certs/agents/"$AGENT_HOSTNAME"_agent.key
DST_AGENT_KEY=$TLS_DIRECTORY/agent/"$AGENT_HOSTNAME"_agent.key
SRC_AGENT_SERVICE=./watchtower_agent.service
DST_AGENT_SERVICE=/etc/systemd/system/watchtower_agent.service

### Core Functionality ###
function remove_previous_installation() {
    if systemctl list-unit-files watchtower_agent.service >/dev/null 2>&1; then
        sudo systemctl stop watchtower_agent.service
        sudo systemctl disable watchtower_agent.service
    fi
}

function prep_important_dirs() {
    start_step_message "Creating Important Directories"
    _create_dir $ETC_DIRECTORY
    _create_dir $OPT_DIRECTORY
    _create_dir $BINARY_DIRECTORY
    _create_dir $TLS_DIRECTORY/ca
    _create_dir $TLS_DIRECTORY/agent
    _create_dir $LOG_DIRECTORY
    successful
}

function _create_dir() {
    if [ ! -d "$1" ]; then
        start_step_message "$1" "substep"
        if ! sudo mkdir -p "$1"; then
            error_message "Failed to create directory '$1'"
        fi
    fi
}

function build_agent_binary() {
    cd .. 
    if ! make build_agent; then 
        exit 1
    fi
    cd scripts
}

function move_important_files() {
    start_step_message "Moving Important Files"
    _move_file $SRC_AGENT_CONFIG $DST_AGENT_CONFIG
    _move_file $SRC_AGENT_BINARY $DST_AGENT_BINARY
    _move_file $SRC_CA_CERT $DST_CA_CERT
    _move_file $SRC_AGENT_CERT $DST_AGENT_CERT
    _move_file $SRC_AGENT_KEY $DST_AGENT_KEY
    _move_file $SRC_AGENT_SERVICE $DST_AGENT_SERVICE
    sudo chown -R root:root $DST_AGENT_SERVICE
    successful
}

function _move_file() {
    start_step_message "$1 -> $2" "substep"
    if ! sudo cp $1 $2; then
        error_message "Failed to move $1 to $2"
    fi
}
function start_systemd_service() {
    start_step_message "Starting watchtower_agent Service '${DST_AGENT_SERVICE}'"
    sudo systemctl daemon-reload
    if ! sudo systemctl enable watchtower_agent.service; then
        error_message "Failed to enable watchtower_agent.service"
    fi

    if ! sudo systemctl restart watchtower_agent.service; then
        error_message "Failed to restart watchtower_agent.service"
    fi
    successful
}

function recap() {
    return
}

function main() {
    remove_previous_installation
    prep_important_dirs
    build_agent_binary
    ./_generate_agent_certs.sh $AGENT_HOSTNAME
    move_important_files
    start_systemd_service
    recap
}

main