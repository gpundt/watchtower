# watchtower

Network observation tower - Golang Server/Agent Deployment

## Overview

This server is deployed onto your homelab to conduct the following activities:

- Network Host Discovery
- Log and Metric Aggregation

This server implements the following technologies / protocols:

- Server / Agent mTLS Communication
- Docker

## How to Deploy

### Server

#### Systemd
1) Run the local server install script: `cd scripts && ./install_server.sh`

This will:
- Remove any previous local Watchtower server installation
- Create important directories
- Build the watchtower_server Golang binary
- Generate server-side certs
- Relocate important configs
- Start the watchtower_server Systemd service

2) View logs with: `sudo systemctl status watchtower_server.service`

#### Docker

- IN PROGRESS

### Agent

#### Systemd
1) Run the local agent install script: `cd scripts && ./install_agent.sh AGENT_HOSTNAME`
- e.g.: `cd scripts && ./install_agent.sh test1`

This will:
- Create important directries (if they're not already present)
- Build the watchtower_agent Golang binary
- Generate agent-side certs
- Relocate important configs
- Start the watchtower_agent Systemd service

2) View logs with: `sudo systemctl status watchtower_agent.service`

#### Docker

- IN PROGRESS


## How to Uninstall

### Server

#### Systemd

1) Run the local server uninstall script: `cd scripts && ./uninstall_server.sh`

This will:
- Stop the watchtower_server Systemd service
- Remove the important directories created during installation
- Preserve `/var/log/watchtower` for debugging

#### Docker

- IN PROGRESS

### Agent

#### Systemd

1) Run the local agent uninstall script: `cd scripts && ./uninstall_agent.sh`

This will:
- Stop the watchtower_agent Systemd service

#### Docker

- IN PROGRESS