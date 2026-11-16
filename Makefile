.PHONY: build-all clean

DIST_DIR := dist
GO_ENV := GOOS=linux GOARCH=amd64 CGO_ENABLED=0

build-all: clean
	mkdir -p $(DIST_DIR)
	$(GO_ENV) go build -o $(DIST_DIR)/api-service ./cmd/api
	$(GO_ENV) go build -o $(DIST_DIR)/bot-service ./cmd/telegram
	$(GO_ENV) go build -o $(DIST_DIR)/sync-service ./cmd/sync

clean:
	rm -rf $(DIST_DIR)