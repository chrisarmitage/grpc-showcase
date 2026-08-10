version:
    go version

build:
    go build -o ./bin/grpc-showcase ./cmd/grpc-showcase/main.go

run-server:
    go run ./cmd/grpc-showcase/main.go server start

run-server-tls:
    go run ./cmd/grpc-showcase/main.go server start --tls=tls

run-client:
    go run ./cmd/grpc-showcase/main.go client run

run-client-tls:
    go run ./cmd/grpc-showcase/main.go client run --tls=tls

run-client-stream:
    go run ./cmd/grpc-showcase/main.go client run --mode=stream

run-client-stream-tls:
    go run ./cmd/grpc-showcase/main.go client run --mode=stream --tls=tls

run-client-bidir:
    go run ./cmd/grpc-showcase/main.go client run --mode=bidir

run-client-bidir-tls:
    go run ./cmd/grpc-showcase/main.go client run --mode=bidir --tls=tls

proto:
	@echo "Generating protobuf code..."
	protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./proto/api.proto
	@echo "Protobuf code generated successfully"

pki:
    go run cmd/grpc-showcase/main.go ca generate
