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
1) Generate server-side certificates: `cd scripts && ./generate_server_certs.sh`
2) Run the local server install script: `cd scripts && ./install_server.sh`

#### Docker

- IN PROGRESS

### Agent

#### Systemd
1) Generate agent-side certificates: `cd scripts && ./generate_agent_certs.sh`
2) Run the local agent install script: `cd scripts && ./install_agent.sh`


#### Docker

- IN PROGRESS


## How to Uninstall

### Server

#### Systemd

1) Run the local server uninstall script: `cd scripts && ./uninstall_server.sh`

#### Docker

- IN PROGRESS

### Agent

#### Systemd

1) Run the local agent uninstall script: `cd scripts && ./uninstall_agent.sh`

#### Docker

- IN PROGRESS