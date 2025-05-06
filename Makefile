.PHONY: build clean test install
.DEFAULT_GOAL := build

BIN_NAME=$(notdir $(CURDIR))
GOBASE=$(shell pwd)
GO_BIN=$(GOBASE)/bin

build: $(GO_BIN)/$(BIN_NAME)

$(GO_BIN)/$(BIN_NAME): go.mod $(shell find $(GOBASE) -name '*.go')
		@echo "  >  Building binary..."
		go build -o $(GO_BIN)/$(BIN_NAME) $(GOBASE)

clean:
		@echo "  >  Cleaning build cache"
		GOBIN=$(GO_BIN) go clean -i ./...

test: build
		@echo "  >  Testing..."
		go test $(GOBASE)/...

run: build
		@echo "  >  Running the project"
		go run $(GOBASE)/...

install: build
		@echo "  >  Installing the project"
		go install .