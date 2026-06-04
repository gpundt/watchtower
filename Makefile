SHELL=/bin/bash
CURRENT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_OUTPUT_DIR := $(CURRENT_DIR)build/
SERVER_BINARY ?= $(BUILD_OUTPUT_DIR)watchtower_server
AGENT_BINARY ?= $(BUILD_OUTPUT_DIR)watchtower_agent


## Colors ##
RED     := \033[0;31m
GREEN   := \033[0;32m
YELLOW  := \033[0;33m
BLUE    := \033[0;34m
CYAN	:= \033[0;36m
RESET   := \033[0m
define start_step_message
	@echo -e "\n$(CYAN)[*] $(1) [*]$(RESET)"
endef
define error_message
	@echo "$(RED)ERROR$(RESET): $(1)"
endef
define successful
	@echo -e "\t - $(GREEN)*Successful*$(RESET)\n"
endef

all: build_server build_agent deploy_docker				## Builds everything

prep_build_output_dirs:
	@mkdir -p $(BUILD_OUTPUT_DIR)

build_server: prep_build_output_dirs		## Builds watchtower server binary
	$(call start_step_message,"Building Server Binary")
	@cd src && \
# 	go mod vendor && \
# 	go mod tidy && \
	go build -mod=vendor -ldflags="-s -w" -o $(SERVER_BINARY) ./cmd/server
	$(call successful)

build_agent: prep_build_output_dirs			## Builds Watchtower agent binary
	$(call start_step_message,"Building Agent Binary")
	@cd src && go build -mod=vendor -ldflags="-s -w" -o $(AGENT_BINARY) ./cmd/agent
	$(call successful)

deploy_docker:								## Deploys docker-compose stack
	$(call start_step_message,"Deploying Docker Stack")
	@sudo docker-compose up -d --build
	$(call successful)

shutdown_docker:
	$(call start_step_message,"Shutting Down Docker Stack")
	@sudo docker-compose down
	$(call successful)

help:										## Prints available make targets
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BLUE)  %-30s$(RESET) %s\n", $$1, $$2}'