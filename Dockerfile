# Build stage
FROM golang:1.25.1-bookworm AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application using Makefile
RUN make build-production && mv bin/server main

# Final stage
FROM ubuntu:24.04

# Install ca-certificates and curl for health checks
RUN apt-get update && \
    apt-get install -y ca-certificates curl tzdata && \
    rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -r -s /bin/false appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Set default port
ENV PORT=8080

# Expose port
EXPOSE $PORT

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:${PORT}/api/health || exit 1

# Run the application
CMD ["./main"]