# watchtower

Network observation tower - Golang Server/Agent Deployment

## Overview

This server is deployed onto your homelab to conduct the following activities:

- Network Host Discovery
- Log and Metric Aggregation

This server implements the following technologies / protocols:

- Server / Agent mTLS Communication
- Docker
- Grafana
- PostgreSQL

## Viewing the Network Dashboard

When the server is successfully deployed, Grafana will be running on port 3000.

Navigate to `http://<SERVER_IP>:3000` and login with the admin credentials listed in `./grafana/grafana.ini`

## How to Deploy

### Server

#### Docker (recommended)

1) Build the server and agent binaries via the Makefile: `make`

2) Generate the server-side certificates: `cd scripts && ./_generate_server_certs.sh`

3) Generate the agent-side certificates: `cd scripts && ./_generate_agent_certs.sh server`

4) Deploy the stack via docker-compose: `docker-compose up -d`

- This will:
    1) Start the Watchtower Database (TimescaleDB PostgreSQL)
    2) Start the Watchtower Dashboard (Grafana)
    3) Start the Watchtower Server (Busybox running the watchtower_server binary)
    4) Start a Watchtower Agent (Busybox rnning the watchtower_agent binary)

5) View your running containers with: `docker ps -a`

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

## How to Uninstall

### Server

#### Systemd

1) Run the local server uninstall script: `cd scripts && ./uninstall_server.sh`

This will:
- Stop the watchtower_server Systemd service
- Remove the important directories created during installation
- Preserve `/var/log/watchtower` for debugging

#### Docker

1) Shutdown the stack with: `docker-compose -f docker-compose.server.yml down`
    - *Can wipe database entries with: `docker-compose -f docker-compose.server.yml down -v`*

2) Confirm there are no running containers with: `docker ps -a`

### Agent

#### Systemd

1) Run the local agent uninstall script: `cd scripts && ./uninstall_agent.sh`

This will:
- Stop the watchtower_agent Systemd service
