APP_NAME   := go-echo-api
BUILD_DIR  := build

.PHONY: install dev clean run build build-all

install:
	go mod tidy

dev:
	go run main.go

run:
	.\$(APP_NAME).exe

build:
	go build -o $(APP_NAME)$(SUFFIX) .

# ── Cross-platform builds ──────────────────────────────────────────

build-all: build-windows-amd64 build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .

build-linux-amd64:
	GOOS=linux   GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .

build-linux-arm64:
	GOOS=linux   GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .

build-darwin-amd64:
	GOOS=darwin  GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .

build-darwin-arm64:
	GOOS=darwin  GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .

# ── Utilities ──────────────────────────────────────────────────────

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(APP_NAME) $(APP_NAME).exe
