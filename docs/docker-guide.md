# Docker Guide

This guide covers building, running, and deploying the Todo API using Docker and Docker Compose.

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Start the application
docker compose up

# Start in background
docker compose up -d

# Stop the application
docker compose down
```

The API will be available at `http://localhost:8080`.

### Using Docker Run

```bash
# Build the image
docker build -t golang-todo-api .

# Run the container
docker run -p 8080:8080 -v $(pwd)/data:/data golang-todo-api
```

## Docker Configuration

### Dockerfile Analysis

The project uses a **multi-stage build** for optimal image size and security:

#### Stage 1: Builder (golang:1.23-bookworm)
```dockerfile
FROM golang:1.23-bookworm AS builder
```

**Purpose**: Compile the Go application with all build tools
- **Base Image**: Debian-based Go image with CGO support
- **Build Tools**: gcc, build-essential for SQLite CGO
- **Dependencies**: Downloads and caches Go modules
- **Build Process**: Compiles with CGO enabled and optimizations

#### Stage 2: Runtime (gcr.io/distroless/base-debian12:nonroot)
```dockerfile
FROM gcr.io/distroless/base-debian12:nonroot
```

**Purpose**: Minimal runtime environment
- **Size**: ~20MB (vs ~300MB+ with full Go image)
- **Security**: No shell, package manager, or unnecessary tools
- **User**: Runs as non-root user for security
- **Only Contains**: glibc and minimal runtime dependencies

### Build Arguments and Optimizations

#### CGO Configuration
```dockerfile
CGO_ENABLED=1 go build -tags 'sqlite_omit_load_extension' -ldflags='-s -w'
```

- **CGO_ENABLED=1**: Required for SQLite performance
- **sqlite_omit_load_extension**: Security - disables extension loading
- **-ldflags='-s -w'**: Strip debug info and symbol table (smaller binary)

#### Build Tags
- **sqlite_omit_load_extension**: Prevents loading of SQLite extensions

## Docker Compose Configuration

### Service Configuration

```yaml
services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    image: golang-todo-api:local
    container_name: golang-todo-api
    environment:
      - PORT=8080
      - HOST=0.0.0.0
      - DATABASE_URL=/data/app.db
    ports:
      - '8080:8080'
    volumes:
      - ./data:/data
    restart: unless-stopped
```

### Key Features

#### Volume Mounting
- **Host Path**: `./data` (relative to compose file)
- **Container Path**: `/data`
- **Purpose**: Persist SQLite database between container restarts

#### Port Mapping
- **Host Port**: 8080
- **Container Port**: 8080
- **Access**: `http://localhost:8080`

#### Environment Variables
- **PORT**: Container internal port
- **HOST**: Bind to all interfaces (0.0.0.0)
- **DATABASE_URL**: Path to SQLite file inside container

#### Restart Policy
- **unless-stopped**: Auto-restart on failure, but not after manual stop

## Building Images

### Local Development Build

```bash
# Build development image
docker build -t golang-todo-api:dev .

# Build with specific tag
docker build -t golang-todo-api:v1.0.0 .

# Build with build args
docker build --build-arg GO_VERSION=1.23 -t golang-todo-api .
```

### Production Build

```bash
# Build optimized production image
docker build -t golang-todo-api:prod --target runtime .

# Multi-platform build
docker buildx build --platform linux/amd64,linux/arm64 -t golang-todo-api:multi .
```

### Build Optimization

#### Layer Caching
The Dockerfile is optimized for layer caching:

```dockerfile
# 1. Copy dependency files first (changes rarely)
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy source code last (changes frequently)
COPY . .
```

#### .dockerignore
```
.git
.env
.env.local
data/
coverage.out
coverage.html
bin/
*.log
README.md
docs/
```

Benefits:
- **Faster builds** - exclude unnecessary files
- **Smaller context** - faster upload to build daemon
- **Security** - don't include sensitive files

## Running Containers

### Development Mode

```bash
# Run with Swagger enabled
docker run -p 8080:8080 \
  -e ENABLE_SWAGGER=true \
  -e GIN_MODE=debug \
  -e ENABLE_LOGGER=true \
  -v $(pwd)/data:/data \
  golang-todo-api

# Run with custom configuration
docker run -p 9000:9000 \
  -e PORT=9000 \
  -e ALLOWED_ORIGINS="http://localhost:3000,http://localhost:3001" \
  -v $(pwd)/data:/data \
  golang-todo-api
```

### Production Mode

```bash
# Run in production mode
docker run -d \
  --name todo-api-prod \
  -p 80:8080 \
  -e GIN_MODE=release \
  -e ENABLE_SWAGGER=false \
  -e ENABLE_RATE_LIMIT=true \
  -v /var/lib/todo-api:/data \
  --restart unless-stopped \
  golang-todo-api:prod

# With resource limits
docker run -d \
  --name todo-api-prod \
  -p 80:8080 \
  -e GIN_MODE=release \
  -v /var/lib/todo-api:/data \
  --restart unless-stopped \
  --memory=512m \
  --cpus=1.0 \
  golang-todo-api:prod
```

### Health Checks

Add health check to Dockerfile:
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/todos || exit 1
```

Or in docker-compose.yml:
```yaml
services:
  api:
    # ... other config
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/v1/todos"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

## Environment Configuration

### Docker Environment Variables

```bash
# Server Configuration
-e PORT=8080
-e HOST=0.0.0.0
-e PUBLIC_SCHEME=https

# Application Features
-e ENABLE_SWAGGER=false
-e ENABLE_LOGGER=true
-e ENABLE_RATE_LIMIT=true

# Framework Configuration
-e GIN_MODE=release

# CORS Configuration
-e ALLOWED_ORIGINS="https://example.com"
-e ALLOW_CREDENTIALS=true

# Database Configuration
-e DATABASE_URL=/data/app.db
```

### Using .env Files

```bash
# Create environment file
cat > docker.env << EOF
PORT=8080
HOST=0.0.0.0
GIN_MODE=release
ENABLE_SWAGGER=false
DATABASE_URL=/data/app.db
EOF

# Run with env file
docker run --env-file docker.env -p 8080:8080 -v $(pwd)/data:/data golang-todo-api
```

### Docker Compose Environment

```yaml
# docker-compose.yml
services:
  api:
    # ... other config
    env_file:
      - docker.env
    # or inline
    environment:
      PORT: 8080
      GIN_MODE: release
      DATABASE_URL: /data/app.db
```

## Volume Management

### Database Persistence

```bash
# Named volume (recommended for production)
docker volume create todo-api-data
docker run -v todo-api-data:/data golang-todo-api

# Host bind mount (good for development)
docker run -v $(pwd)/data:/data golang-todo-api

# Temporary volume (testing)
docker run -v /data golang-todo-api
```

### Volume Backup

```bash
# Backup named volume
docker run --rm -v todo-api-data:/data -v $(pwd):/backup alpine \
  tar czf /backup/todo-backup-$(date +%Y%m%d).tar.gz -C /data .

# Restore from backup
docker run --rm -v todo-api-data:/data -v $(pwd):/backup alpine \
  tar xzf /backup/todo-backup-20230101.tar.gz -C /data
```

## Networking

### Default Bridge Network

```bash
# Containers can communicate by container name
docker run --name todo-api golang-todo-api
docker run --link todo-api nginx  # Can reach via 'todo-api' hostname
```

### Custom Networks

```bash
# Create custom network
docker network create todo-network

# Run containers in custom network
docker run --network todo-network --name api golang-todo-api
docker run --network todo-network --name proxy nginx
```

### Docker Compose Networking

```yaml
# docker-compose.yml
services:
  api:
    networks:
      - backend
  
  proxy:
    networks:
      - frontend
      - backend

networks:
  frontend:
  backend:
```

## Development Workflow

### Live Reload Setup

```bash
# Mount source code for development
docker run -p 8080:8080 \
  -v $(pwd):/app \
  -v $(pwd)/data:/data \
  -w /app \
  golang:1.23 \
  go run cmd/server/main.go
```

### Development Compose

```yaml
# docker-compose.dev.yml
services:
  api:
    build:
      context: .
      dockerfile: Dockerfile.dev  # Development Dockerfile
    volumes:
      - .:/app
      - ./data:/data
    environment:
      - GIN_MODE=debug
      - ENABLE_SWAGGER=true
      - ENABLE_LOGGER=true
    ports:
      - "8080:8080"
    command: go run cmd/server/main.go
```

```bash
# Use development compose
docker compose -f docker-compose.dev.yml up
```

## Production Deployment

### Production Best Practices

#### 1. Multi-stage Build
```dockerfile
# Use multi-stage build for smaller images
FROM golang:1.23-bookworm AS builder
# ... build stage

FROM gcr.io/distroless/base-debian12:nonroot
# ... runtime stage
```

#### 2. Security Hardening
```dockerfile
# Run as non-root user
USER nonroot:nonroot

# Use distroless base image
FROM gcr.io/distroless/base-debian12:nonroot

# Set security labels
LABEL security.non-root=true
```

#### 3. Resource Limits
```yaml
# docker-compose.prod.yml
services:
  api:
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '1.0'
        reservations:
          memory: 256M
          cpus: '0.5'
```

### Container Orchestration

#### Docker Swarm
```bash
# Deploy to swarm
docker stack deploy -c docker-compose.prod.yml todo-api

# Scale service
docker service scale todo-api_api=3
```

#### Kubernetes
```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todo-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: todo-api
  template:
    metadata:
      labels:
        app: todo-api
    spec:
      containers:
      - name: api
        image: golang-todo-api:prod
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: GIN_MODE
          value: "release"
        volumeMounts:
        - name: data
          mountPath: /data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: todo-api-pvc
```

## Monitoring and Logging

### Container Logs

```bash
# View logs
docker compose logs api

# Follow logs
docker compose logs -f api

# Last 100 lines
docker compose logs --tail=100 api

# Logs with timestamps
docker compose logs -t api
```

### Log Drivers

```yaml
# docker-compose.yml
services:
  api:
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"
```

### Monitoring

```bash
# Container stats
docker stats todo-api

# Resource usage
docker compose exec api sh -c "cat /proc/meminfo | head -5"

# Process list
docker compose exec api ps aux
```

## Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Error: port is already allocated
# Solution: Use different port
docker run -p 8081:8080 golang-todo-api
```

#### 2. Volume Permissions
```bash
# Error: permission denied
# Solution: Fix ownership
sudo chown -R 1000:1000 ./data
# Or run with user flag
docker run --user $(id -u):$(id -g) golang-todo-api
```

#### 3. Database Lock
```bash
# Error: database is locked
# Solution: Stop all containers using the database
docker compose down
```

#### 4. Out of Disk Space
```bash
# Clean up Docker
docker system prune -a

# Remove unused volumes
docker volume prune

# Check disk usage
docker system df
```

### Debugging

#### Container Shell Access
```bash
# Access running container (if shell available)
docker compose exec api sh

# Debug with different image
docker run -it --rm -v $(pwd)/data:/data alpine sh
```

#### Image Inspection
```bash
# Inspect image layers
docker history golang-todo-api

# Inspect image metadata
docker inspect golang-todo-api

# Check image size
docker images golang-todo-api
```

#### Network Debugging
```bash
# Test network connectivity
docker compose exec api wget -O- http://localhost:8080/api/v1/todos

# Check port binding
docker compose ps
```

## Performance Optimization

### Image Size Optimization

```dockerfile
# Use minimal base image
FROM gcr.io/distroless/static-debian12

# Multi-stage build
FROM golang:alpine AS builder
# ... build
FROM alpine
# ... runtime

# Remove unnecessary files
RUN rm -rf /var/cache/apk/*
```

### Build Performance

```bash
# Use BuildKit for faster builds
DOCKER_BUILDKIT=1 docker build .

# Parallel builds
docker build --parallel .

# Build cache mounting
RUN --mount=type=cache,target=/go/pkg/mod go mod download
```

### Runtime Performance

```yaml
# docker-compose.yml
services:
  api:
    ulimits:
      nofile:
        soft: 65536
        hard: 65536
    deploy:
      resources:
        limits:
          memory: 1G
        reservations:
          memory: 512M
```

## Security Considerations

### Image Security
- Use **distroless** base images
- Run as **non-root** user
- Keep base images **updated**
- Scan images for **vulnerabilities**

### Runtime Security
- Use **read-only** root filesystem where possible
- Limit **capabilities**
- Use **secrets** for sensitive data
- Enable **AppArmor/SELinux** if available

### Network Security
- Use **custom networks**
- Limit **exposed ports**
- Enable **TLS** for external communication
- Use **firewall** rules

```bash
# Security scan
docker scan golang-todo-api

# Run with security options
docker run --read-only --tmpfs /tmp golang-todo-api
```