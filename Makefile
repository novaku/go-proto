.PHONY: all build test clean generate tools

PROTOC_GEN_GO := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(shell go env GOPATH)/bin/protoc-gen-go-grpc

all: generate build

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

proto-deps:
	@echo "Downloading proto dependencies..."
	@mkdir -p proto/protoc-gen-openapiv2/options
	@curl -sSL https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/annotations.proto \
		-o proto/protoc-gen-openapiv2/options/annotations.proto
	@curl -sSL https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/openapiv2.proto \
		-o proto/protoc-gen-openapiv2/options/openapiv2.proto

generate: tools proto-deps
	mkdir -p gen
	protoc --proto_path=proto \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		--plugin=protoc-gen-grpc-gateway=$(shell go env GOPATH)/bin/protoc-gen-grpc-gateway \
		--plugin=protoc-gen-openapiv2=$(shell go env GOPATH)/bin/protoc-gen-openapiv2 \
		--go_out=gen --go_opt=paths=source_relative \
		--go-grpc_out=gen --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=gen --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=gen --openapiv2_opt=logtostderr=true \
		proto/guestbook/v1/guestbook.proto

build:
	go build -o bin/server cmd/server/main.go

test:
	go test ./...

run:
	go run cmd/server/main.go

swagger: generate
	@echo "Preparing Swagger JSON..."
	@mkdir -p swagger
	@cp gen/guestbook/v1/guestbook.swagger.json swagger/guestbook.swagger.json
	@echo "Killing any existing Swagger server on :8081 (if running)..."
	@lsof -ti:8081 | xargs kill -9 2>/dev/null || true
	@echo "Serving Swagger UI and API docs on http://localhost:8081 (Ctrl+C to stop) ..."
	@open "http://localhost:8081/swagger-ui.html"
	@python3 -m http.server 8081

clean:
	rm -rf gen bin swagger
