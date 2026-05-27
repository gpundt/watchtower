#!/bin/bash
source ./helpers.sh

### Important Filepaths ###
DST_AGENT_SERVICE=/etc/systemd/system/watchtower_agent.service

### Core Functionality ###
function stop_systemd_service() {
    if systemctl list-unit-files watchtower_agent.service >/dev/null 2>&1; then
        start_step_message "Stopping Systemd Service"
        sudo systemctl disable --now watchtower_agent.service
        sudo systemctl stop watchtower_agent.service
        sudo rm -rf $DST_AGENT_SERVICE
        successful 
    fi
}

function main() {
    stop_systemd_service
}

main