version:
    go version

build:
    go build -o ./bin/grpc-showcase ./cmd/grpc-showcase/main.go

run-server:
    go run ./cmd/grpc-showcase/main.go server start

run-client:
    go run ./cmd/grpc-showcase/main.go client run

proto:
	@echo "Generating protobuf code..."
	protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./proto/api.proto
	@echo "Protobuf code generated successfully"
