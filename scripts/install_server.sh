#!/bin/bash
source ./helpers.sh

### Important Directories ###
ETC_DIRECTORY=/etc/watchtower
OPT_DIRECTORY=/opt/watchtower
BINARY_DIRECTORY=$OPT_DIRECTORY/bin
TLS_DIRECTORY=$OPT_DIRECTORY/tls
LOG_DIRECTORY=/var/log/watchtower

### Important Filepaths ###
SRC_SERVER_CONFIG=../config/server.yaml
DST_SERVER_CONFIG=$ETC_DIRECTORY/server.yaml
SRC_SERVER_BINARY=../build/watchtower_server
DST_SERVER_BINARY=$BINARY_DIRECTORY/watchtower_server
SRC_CA_CERT=../certs/ca/ca.crt
DST_CA_CERT=$TLS_DIRECTORY/ca.crt
SRC_CA_KEY=../certs/ca/ca.key
DST_CA_KEY=$TLS_DIRECTORY/ca.key
SRC_SERVER_CERT=../certs/server/server.crt
DST_SERVER_CERT=$TLS_DIRECTORY/server.crt
SRC_SERVER_KEY=../certs/server/server.key
DST_SERVER_KEY=$TLS_DIRECTORY/server.key
SRC_SERVER_SERVICE=./watchtower_server.service
DST_SERVER_SERVICE=/etc/systemd/system/watchtower_server.service

### Core Functionality ###
function remove_previous_installation() {
    sudo rm -rf $DST_SERVER_CONFIG \
        $DST_SERVER_BINARY \
        $DST_CA_CERT \
        $DST_CA_KEY \
        $DST_SERVER_CERT \
        $DST_SERVER_KEY \
        $DST_SERVER_SERVICE
}

function prep_important_dirs() {
    start_step_message "Creating Important Directories"
    _create_dir $ETC_DIRECTORY
    _create_dir $OPT_DIRECTORY
    _create_dir $BINARY_DIRECTORY
    _create_dir $TLS_DIRECTORY
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

function build_server_binary() {
    cd .. 
    if ! make build_server; then 
        exit 1
    fi
    cd scripts
}

function move_important_files() {
    start_step_message "Moving Important Files"
    _move_file $SRC_SERVER_CONFIG $DST_SERVER_CONFIG
    _move_file $SRC_SERVER_BINARY $DST_SERVER_BINARY
    _move_file $SRC_CA_CERT $DST_CA_CERT
    _move_file $SRC_CA_KEY $DST_CA_KEY
    _move_file $SRC_SERVER_CERT $DST_SERVER_CERT
    _move_file $SRC_SERVER_KEY $DST_SERVER_KEY
    _move_file $SRC_SERVER_SERVICE $DST_SERVER_SERVICE
    sudo chown -R root:root $DST_SERVER_SERVICE
    successful
}

function _move_file() {
    start_step_message "$1 -> $2" "substep"
    if ! sudo cp $1 $2; then
        error_message "Failed to move $1 to $2"
    fi
}
function start_systemd_service() {
    start_step_message "Starting watchtower_server Service '${DST_SERVER_SERVICE}'"
    sudo systemctl daemon-reload
    if ! sudo systemctl enable watchtower_server.service; then
        error_message "Failed to enable watchtower_server.service"
    fi

    if ! sudo systemctl restart watchtower_server.service; then
        error_message "Failed to restart watchtower_server.service"
    fi
    successful
}

function recap() {
    return
}

function main() {
    remove_previous_installation
    prep_important_dirs
    build_server_binary
    move_important_files
    start_systemd_service
    recap
}

main