ECHO_PB_PATH := "echo/echo.pb.go"
SERVER_BIN_PATH := "bin/server"
CLIENT_BIN_PATH := "bin/client"

.PHONY: all proto

all: build-server build-client

build-client: proto ## build client binary
	@echo "+ $@"
	@go build -i -o $(CLIENT_BIN_PATH) github.com/nathanows/ues/client
	@echo "    $(CLIENT_BIN_PATH)"

build-server: proto ## build server binary
	@echo "+ $@"
	@go build -i -o $(SERVER_BIN_PATH) github.com/nathanows/ues
	@echo "    $(SERVER_BIN_PATH)"

clean: ## remove all build artifacts
	@echo "+ $@"
	@rm -f $(ECHO_PB_PATH) $(SERVER_BIN_PATH) $(CLIENT_BIN_PATH)

proto: ## compile .proto files
	@echo "+ $@"
	@docker build -t protogen -f Dockerfile.protogen .
	@docker run --name protogen protogen
	@docker cp protogen:/proto/gen/echo.pb.go $(ECHO_PB_PATH)
	@docker rm protogen

server: build-server ## run local server
	@echo "+ $@"
	@clear
	@bin/server

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
