.PHONY: build install clean test

BINARY_NAME=twt
INSTALL_DIR=$(HOME)/.local/bin

build:
	go build -o $(BINARY_NAME) ./cmd/twt

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/
	chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Installed to $(INSTALL_DIR)/$(BINARY_NAME)"

clean:
	rm -f $(BINARY_NAME)
	go clean

test:
	go test -v ./...

run:
	go run ./cmd/twt

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  install  - Build and install to $(INSTALL_DIR)"
	@echo "  clean    - Remove built binaries"
	@echo "  test     - Run tests"
	@echo "  run      - Run without building"
