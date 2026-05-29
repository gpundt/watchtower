#!/bin/bash
source ./_helpers.sh

### Important Directories ###
ETC_DIRECTORY=/etc/watchtower
OPT_DIRECTORY=/opt/watchtower
BINARY_DIRECTORY=$OPT_DIRECTORY/bin
TLS_DIRECTORY=$OPT_DIRECTORY/tls
LOG_DIRECTORY=/var/log/watchtower

### Important Filepaths ###
DST_SERVER_SERVICE=/etc/systemd/system/watchtower_server.service

### Core Functionality ###
function stop_systemd_service() {
    start_step_message "Stopping Systemd Service"
    if systemctl list-unit-files watchtower_server.service >/dev/null 2>&1; then
        sudo systemctl disable --now watchtower_server.service
        sudo systemctl stop watchtower_server.service
        sudo rm -rf $DST_SERVER_SERVICE
        successful 
    else
        warning_message "watchtower_server.service not running... skipping"
    fi
}

function remove_important_directories() {
    start_step_message "Removing Important Directories"
    _remove_dir $ETC_DIRECTORY
    _remove_dir $OPT_DIRECTORY
    _remove_dir $BINARY_DIRECTORY
    _remove_dir $TLS_DIRECTORY
    successful
}

function _remove_dir() {
    start_step_message "$1" "substep"
    sudo rm -rf $1
}

function main() {
    stop_systemd_service
    remove_important_directories
}

main