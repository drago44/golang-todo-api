# Multi-stage build for a small final image

# Builder stage (Debian-based for CGO)
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Install build tools for CGO
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  build-essential \
  ca-certificates \
  git \
  && rm -rf /var/lib/apt/lists/* \
  && update-ca-certificates

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary with CGO for sqlite
RUN CGO_ENABLED=1 go build -tags 'sqlite_omit_load_extension' -ldflags='-s -w' -o /bin/golang-todo-api ./cmd/server

# Runtime stage (non-static, includes glibc)
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /

# Create app directories
USER nonroot:nonroot

# Copy binary from builder
COPY --from=builder /bin/golang-todo-api /golang-todo-api

# Expose default port
EXPOSE 8080

# Default envs (can be overridden)
ENV PORT=8080 \
  HOST=0.0.0.0 \
  DATABASE_URL=/data/app.db

# Create data volume for sqlite file
VOLUME ["/data"]

# Run the server
ENTRYPOINT ["/golang-todo-api"]

