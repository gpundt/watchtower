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

prep_build_output_dirs:
	$(call start_step_message,"Prepping Build Output Dir")
	@mkdir -p $(BUILD_OUTPUT_DIR)
	$(call successful)

build_server_binary: prep_build_output_dirs					## Builds watchtower server binary
	$(call start_step_message,"Building Server Binary")
	@go mod vendor
	@go mod tidy
	@go build -ldflags="-s -w" -o $(SERVER_BINARY) ./cmd/server
	$(call successful)

build_agent_binary:											## Builds Watchtower agent binary
	$(call start_step_message,"Building Agent Binary")
	@go build -ldflags="-s -w" -o $(AGENT_BINARY) ./cmd/agent
	$(call successful)

help:						## Prints available make targets
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BLUE)  %-30s$(RESET) %s\n", $$1, $$2}'