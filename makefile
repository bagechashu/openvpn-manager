# Makefile for compiling Golang Linux installation package and frontend package

# builder and cleaner
GO := go
GO_BUILD := $(GO) build
GO_CLEAN := $(GO) clean
GO_INSTALL := $(GO) install

# cross-compiling Go for Linux
GOOS := linux
GOARCH := amd64

# Variables
APP_NAME := ovpnmanager-$(GOOS)-$(GOARCH)
BIN_DIR := ./bin
GO_SRC := ./cmd/openvpn-manager/main.go
WEB_SRC_DIR := websrc

# compile
build:
	@echo "compile $(GOOS) $(GOARCH)..."
	env GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) -o $(BIN_DIR)/$(APP_NAME) $(GO_SRC)
	@echo "target: $(BIN_DIR)/$(APP_NAME)"

clean:
	$(GO_CLEAN)
	rm -f $(APP_NAME)

fe:
	@echo "compile webui..."
	cd $(WEB_SRC_DIR) && npm run build

# 默认目标
.DEFAULT_GOAL := build

# 设置伪目标
.PHONY: build clean fe
