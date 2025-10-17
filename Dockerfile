# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git make

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Copy the source code
COPY . .

# Create stub packages for missing imports with proper go.mod files
RUN mkdir -p /app/stub/docs && \
    echo 'module taskbounty/docs' > /app/stub/docs/go.mod && \
    echo 'go 1.25' >> /app/stub/docs/go.mod && \
    echo 'package docs' > /app/stub/docs/docs.go && \
    echo 'func RegisterOpenAPIService(name string, router interface{}) {}' >> /app/stub/docs/docs.go && \
    mkdir -p /app/stub/testutil/sample && \
    echo 'module taskbounty/testutil/sample' > /app/stub/testutil/sample/go.mod && \
    echo 'go 1.25' >> /app/stub/testutil/sample/go.mod && \
    echo 'package sample' > /app/stub/testutil/sample/sample.go && \
    echo '// AccAddress returns a sample account address' >> /app/stub/testutil/sample/sample.go && \
    echo 'func AccAddress() string {' >> /app/stub/testutil/sample/sample.go && \
    echo '    return "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"' >> /app/stub/testutil/sample/sample.go && \
    echo '}' >> /app/stub/testutil/sample/sample.go

# Update go.mod to use the stub packages
RUN go mod edit -replace taskbounty/docs=./stub/docs && \
    go mod edit -replace taskbounty/testutil/sample=./stub/testutil/sample && \
    go mod tidy

# Build the actual blockchain application
RUN go build -trimpath -ldflags="-s -w" -o taskbountyd ./cmd/taskbountyd

# ---------- Stage 2: Final ----------
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata bash

# Copy the binary from the builder stage
COPY --from=builder /app/taskbountyd /usr/local/bin/taskbountyd

# Copy the startup script
COPY start-taskbounty.sh /app/start-taskbounty.sh
RUN chmod +x /app/start-taskbounty.sh

# Create data directory for blockchain data
RUN mkdir -p /root/.taskbounty && \
    chmod 755 /root/.taskbounty

# Expose ports for blockchain operations
# 26656: P2P
# 26657: RPC
# 1317: REST API
# 9090: gRPC
# 9091: gRPC-web
EXPOSE 26656 26657 1317 9090 9091

# Default command - this will be overridden by docker-compose
ENTRYPOINT ["/app/start-taskbounty.sh"]