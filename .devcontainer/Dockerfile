# syntax=docker/dockerfile:1
FROM alpine:3.20

# Install Go and basic build tools
ARG GO_VERSION=1.26.0
RUN apk add --no-cache curl bash git build-base aws-cli just protoc protobuf-dev \
  && curl -L https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
  | tar -C /usr/local -xz \
  && ln -s /usr/local/go/bin/* /usr/local/bin/ \
  && go version \
  && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Create dev user
RUN addgroup -S dev && adduser -S dev -G dev
USER dev
WORKDIR /workspaces/grpc-showcase
