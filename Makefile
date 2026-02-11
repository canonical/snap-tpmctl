##### Remote VM testing Makefile #####

# Default VM connection (override with: make VM_HOST=...)
VM_USER ?= root
VM_HOST ?= 192.168.122.178
VM      := $(VM_USER)@$(VM_HOST)

# Remote project directory (override with: make REMOTE_DIR=...)
REMOTE_DIR ?= ~/snap-tpmctl

# SSH / RSYNC options
SSH_OPTS ?= -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR -tt
RSYNC_OPTS ?= -az --delete --exclude .git --exclude bin --exclude '*.snap' --exclude '*.swp'

# Binary name
BIN_NAME := snap-tpmctl
LOCAL_BIN := bin/$(BIN_NAME)

.PHONY: help build clean sync remote-build remote-test remote-status remote-run remote-shell remote-clean run test check snap snap-clean remote-snap

# Catch-all rule to prevent Make from treating arguments as targets
%:
	@:

build:
	@mkdir -p bin
	go build -o $(LOCAL_BIN) ./cmd/tpmctl

snap:
	@echo 'Building snap locally...'
	@snapcraft pack

snap-clean:
	@echo 'Cleaning snap artifacts...'
	@snapcraft clean
	@rm -f *.snap
	@echo 'Snap artifacts cleaned.'

remote-snap: sync
	@echo 'Building snap on remote VM...'
	@ssh $(SSH_OPTS) $(VM) 'cd $(REMOTE_DIR) && snapcraft pack'
	@echo 'Installing snap on remote VM...'
	@ssh $(SSH_OPTS) $(VM) 'cd $(REMOTE_DIR) && snap install *.snap --dangerous'
	@echo 'Connecting snap interfaces...'
	@ssh $(SSH_OPTS) $(VM) 'snap connect $(BIN_NAME):snapd-control'
	@ssh $(SSH_OPTS) $(VM) 'snap connect $(BIN_NAME):hardware-observe'
	@ssh $(SSH_OPTS) $(VM) 'snap connect $(BIN_NAME):mountctl'
	@ssh $(SSH_OPTS) $(VM) 'snap connect $(BIN_NAME):mount-observe'
	@ssh $(SSH_OPTS) $(VM) 'snap connect $(BIN_NAME):block-devices'
	@ssh $(SSH_OPTS) $(VM) 'snap connect $(BIN_NAME):dm-crypt'
	@ssh $(SSH_OPTS) $(VM) 'snap connect $(BIN_NAME):mkdir'
	@echo 'Snap installed and configured on remote VM.'

run:
	@go run ./cmd/tpmctl $(filter-out $@,$(MAKECMDGOALS))

clean:
	@rm -rf bin
	@echo 'Local artifacts cleaned.'

sync:
	@rsync $(RSYNC_OPTS) ./ $(VM):$(REMOTE_DIR)

remote-build: sync
	@ssh $(SSH_OPTS) $(VM) 'cd $(REMOTE_DIR) && mkdir -p bin && go build -o $(LOCAL_BIN) ./cmd/tpmctl'

remote-run: remote-build
	@ssh $(SSH_OPTS) $(VM) 'cd $(REMOTE_DIR) && $(LOCAL_BIN) $(filter-out $@,$(MAKECMDGOALS))'

remote-test: sync
	@echo 'Running tests on remote VM...'
	@ssh $(SSH_OPTS) $(VM) 'cd $(REMOTE_DIR) && go test -v ./...'

remote-shell:
	@ssh $(SSH_OPTS) $(VM)

remote-clean:
	@echo 'Cleaning remote artifacts...'
	@ssh $(SSH_OPTS) $(VM) 'cd $(REMOTE_DIR) && rm -rf bin'
	@echo 'Remote artifacts cleaned.'

check:
	@golangci-lint-v2 run ./... --fix
