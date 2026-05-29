#!/bin/bash
source ./_helpers.sh

### Important Filepaths ###
DST_AGENT_SERVICE=/etc/systemd/system/watchtower_agent.service

### Core Functionality ###
function stop_systemd_service() {
    start_step_message "Stopping Systemd Service"
    if systemctl list-unit-files watchtower_agent.service >/dev/null 2>&1; then
        sudo systemctl disable --now watchtower_agent.service
        sudo systemctl stop watchtower_agent.service
        sudo rm -rf $DST_AGENT_SERVICE
        successful 
    else
        warning_message "watchtower_agent.service not running... skipping"
    fi
}

function main() {
    stop_systemd_service
}

main