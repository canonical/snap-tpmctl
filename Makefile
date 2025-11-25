##### Remote VM testing Makefile #####

# Default VM connection (override with: make VM_HOST=...)
VM_USER ?= ubuntu
VM_HOST ?= 192.168.122.105
VM      := $(VM_USER)@$(VM_HOST)

# Remote project directory (override with: make REMOTE_DIR=...)
REMOTE_DIR ?= ~/snap-tpmctl

# SSH / RSYNC options
SSH_OPTS ?= -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR
RSYNC_OPTS ?= -az --delete --exclude .git --exclude bin --exclude '*.snap' --exclude '*.swp'

# Binary name
BIN_NAME := tpmctl
LOCAL_BIN := bin/$(BIN_NAME)

.PHONY: help build clean sync remote-build remote-test remote-status remote-run remote-shell remote-clean run test

# Catch-all rule to prevent Make from treating arguments as targets
%:
	@:

build:
	@mkdir -p bin
	go build -o $(LOCAL_BIN) ./cmd/tpmctl

run:
	@go run ./cmd/tpmctl $(filter-out $@,$(MAKECMDGOALS))

clean:
	@rm -rf bin
	@echo 'Local artifacts cleaned.'

sync:
	@rsync $(RSYNC_OPTS) ./ $(VM):$(REMOTE_DIR)

remote-build: sync
	@ssh $(SSH_OPTS) $(VM) 'cd $(REMOTE_DIR) && mkdir -p bin && go build -o bin/$(BIN_NAME) ./cmd/tpmctl'

remote-run: remote-build
	@ssh $(SSH_OPTS) $(VM) 'cd $(REMOTE_DIR) && ./bin/$(BIN_NAME) $(filter-out $@,$(MAKECMDGOALS))'

remote-test: sync
	@echo 'Running tests on remote VM...'
	@ssh $(SSH_OPTS) $(VM) 'cd $(REMOTE_DIR) && go test -v ./...'

remote-clean:
	@echo 'Cleaning remote artifacts...'
	@ssh $(SSH_OPTS) $(VM) 'cd $(REMOTE_DIR) && rm -rf bin'
	@echo 'Remote artifacts cleaned.'
