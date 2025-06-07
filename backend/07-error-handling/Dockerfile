# syntax=docker/dockerfile:1

FROM golang:1.21-alpine3.18

# Install additional tools for debugging and health checks
RUN apk add --no-cache \
    curl \
    netcat-openbsd \
    ca-certificates

WORKDIR /app

# Copy go module files first for better caching
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . ./

# Build the application
RUN go build -o error-handling-server main.go

# Create non-root user for security
RUN adduser -D -s /bin/sh appuser
USER appuser

# Expose the application port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./error-handling-server"] 