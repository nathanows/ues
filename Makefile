ECHO_PB_PATH := "echo/echo.pb.go"
SERVER_BIN_PATH := "bin/server"
CLIENT_BIN_PATH := "bin/client"
AUTH_TOKEN := "super-secret"
SERVER_ADDR := "localhost:6000"

.PHONY: build-client build-server clean container gen-certs protogen server test help

build-client: ## build client binary
	@echo "+ $@"
	@go build -i -o $(CLIENT_BIN_PATH) github.com/nathanows/ues/client

build-server: ## build server binary
	@echo "+ $@"
	@go build -i -o $(SERVER_BIN_PATH) github.com/nathanows/ues/server

clean: ## remove all build artifacts
	@echo "+ $@"
	@rm -f $(ECHO_PB_PATH) $(SERVER_BIN_PATH) $(CLIENT_BIN_PATH)

client-request: ## make sample client request
	@echo "+ $@"
	@TOKEN=$(AUTH_TOKEN) SERVER_ADDR=$(SERVER_ADDR) bin/client test again "something else" 123

container: ## build Docker image
	@echo "+ $@"
	@docker build -t ues:latest .

gen-certs: ## generate self-signed SSL cert
	@echo "+ $@"
	@openssl req -x509 -newkey rsa:4096 -keyout certs/server-key.pem -out certs/server-cert.pem -days 365 -nodes -subj '/CN=localhost'

protogen: ## compile .proto files
	@echo "+ $@"
	@docker build -t protogen -f Dockerfile.protogen .
	@docker run --name protogen protogen
	@docker cp protogen:/proto/gen/echo.pb.go $(ECHO_PB_PATH)
	@docker rm protogen

server: gen-certs build-server ## run local server
	@echo "+ $@"
	@TOKEN=$(AUTH_TOKEN) SERVER_ADDR=$(SERVER_ADDR) bin/server

test: ## run full test suite
	@echo "+ $@"
	@TOKEN=$(AUTH_TOKEN) go test ./...

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
