BINARY_NAME=dextrace

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod

BUILD_DIR=build


.PHONY: all build test clean run deps

all: deps build

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go
	./$(BUILD_DIR)/$(BINARY_NAME)

deps:
	$(GOMOD) download
	$(GOMOD) tidy


dev:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go
	./$(BUILD_DIR)/$(BINARY_NAME)

watch:
	@echo "Watching for changes..."
	@while true; do \
		make build; \
		sleep 2; \
	done

hot-reload:
	@echo "Running with hot reload..."
	@while true; do \
		make run; \
		sleep 2; \
	done