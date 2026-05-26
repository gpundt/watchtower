# watchtower
Network observation tower - Golang

## How to Deploy

### Locally

1) Generate server-side certificates: `cd scripts && ./generate_server_certs.sh`
2) Generate agent-side certificates: `cd scripts && ./generate_agent_certs.sh`
3) Build the binaries: `make build_server_binary`
4) Run the server binary: `cd build && ./watchtower_server`
5) Run the agent binary: `cd build && ./watchtower_agent`
