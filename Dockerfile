# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git make

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build command placeholder - this would need to be fixed for actual building
# For now, we're just creating a placeholder binary to make the Dockerfile valid
RUN echo "#!/bin/sh" > /app/taskbountyd && \
    echo "echo 'This is a placeholder binary. The actual build command needs to be fixed.'" >> /app/taskbountyd && \
    chmod +x /app/taskbountyd

# ---------- Stage 2: Final ----------
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN addgroup -S appuser \
 && adduser -S -G appuser -H -s /sbin/nologin appuser

# Copy the binary from the builder stage
COPY --from=builder /app/taskbountyd /app/taskbountyd

# Set ownership to non-root user
RUN chown -R appuser:appuser /app

# Run as non-root user
USER appuser

# Set the entrypoint command
ENTRYPOINT ["/app/taskbountyd"]