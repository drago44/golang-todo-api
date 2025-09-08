# Configuration Guide

This document describes all configuration options available for the Todo API application.

## Configuration Overview

The application uses **environment variables** for configuration, providing flexibility for different deployment environments. All configuration is loaded at application startup with validation and helpful error messages.

## Environment Files

### `.env.example`
Template file with all available configuration options and their default values.

### `.env`
Your local configuration file (copy from `.env.example`). This file is **not** version controlled.

### Usage
```bash
# Copy template to create your local config
cp .env.example .env

# Edit your local configuration
nano .env
```

## Configuration Sections

### Server Configuration

#### PORT
- **Default**: `8080`
- **Type**: Integer
- **Description**: HTTP server port
- **Example**: `PORT=3000`

#### HOST
- **Default**: `localhost`
- **Type**: String
- **Description**: Server host/bind address
- **Example**: `HOST=0.0.0.0` (bind to all interfaces)
- **Production**: Use `0.0.0.0` to accept external connections

#### PUBLIC_SCHEME
- **Default**: `http`
- **Type**: String
- **Values**: `http`, `https`
- **Description**: Public-facing scheme for URL generation
- **Example**: `PUBLIC_SCHEME=https`

### Application Features

#### ENABLE_SWAGGER
- **Default**: `false`
- **Type**: Boolean
- **Description**: Enable/disable Swagger UI for API documentation
- **Example**: `ENABLE_SWAGGER=true`
- **Access**: Available at `/swagger/index.html` when enabled
- **Production**: Usually disabled in production

#### ENABLE_LOGGER
- **Default**: `true`
- **Type**: Boolean
- **Description**: Enable/disable HTTP request logging
- **Example**: `ENABLE_LOGGER=false`
- **Output**: Logs requests with method, path, status, and duration

#### ENABLE_RATE_LIMIT
- **Default**: `false`
- **Type**: Boolean
- **Description**: Enable/disable rate limiting
- **Example**: `ENABLE_RATE_LIMIT=true`
- **Implementation**: Global rate limiting for all endpoints

### Framework Configuration

#### GIN_MODE
- **Default**: `release`
- **Type**: String
- **Values**: `debug`, `release`
- **Description**: Gin framework operation mode
- **Example**: `GIN_MODE=debug`
- **Debug Mode**: 
  - Verbose logging
  - Debug middleware
  - Development-friendly error messages
- **Release Mode**:
  - Minimal logging
  - Production optimizations
  - Better performance

#### TRUSTED_PROXIES
- **Default**: `` (empty)
- **Type**: Comma-separated list
- **Description**: List of trusted proxy IPs or CIDR blocks
- **Example**: `TRUSTED_PROXIES=192.168.1.0/24,10.0.0.1`
- **Security**: Important for proper IP address detection behind proxies

### CORS Configuration

#### ALLOWED_ORIGINS
- **Default**: `http://localhost:3000`
- **Type**: Comma-separated list
- **Description**: List of allowed CORS origins
- **Examples**:
  ```bash
  # Single origin
  ALLOWED_ORIGINS=http://localhost:3000
  
  # Multiple origins
  ALLOWED_ORIGINS=http://localhost:3000,https://example.com
  
  # All origins (development only)
  ALLOWED_ORIGINS=*
  ```
- **Production**: Never use `*` in production

#### ALLOW_CREDENTIALS
- **Default**: `true`
- **Type**: Boolean
- **Description**: Allow credentials in CORS requests
- **Example**: `ALLOW_CREDENTIALS=false`
- **Note**: Cannot be `true` when `ALLOWED_ORIGINS=*`

### Database Configuration

#### DATABASE_URL
- **Default**: `data/app.db`
- **Type**: String (file path)
- **Description**: SQLite database file path
- **Examples**:
  ```bash
  # Relative path
  DATABASE_URL=data/app.db
  
  # Absolute path
  DATABASE_URL=/var/lib/todo-api/app.db
  
  # In-memory (testing)
  DATABASE_URL=:memory:
  ```
- **Directory**: Ensure the directory exists and is writable

## Configuration Examples

### Development Configuration
```bash
# .env for development
PORT=8080
HOST=localhost
PUBLIC_SCHEME=http
ENABLE_SWAGGER=true
ENABLE_LOGGER=true
ENABLE_RATE_LIMIT=false
GIN_MODE=debug
TRUSTED_PROXIES=
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
ALLOW_CREDENTIALS=true
DATABASE_URL=data/dev.db
```

### Production Configuration
```bash
# .env for production
PORT=80
HOST=0.0.0.0
PUBLIC_SCHEME=https
ENABLE_SWAGGER=false
ENABLE_LOGGER=true
ENABLE_RATE_LIMIT=true
GIN_MODE=release
TRUSTED_PROXIES=192.168.1.0/24
ALLOWED_ORIGINS=https://yourapp.com
ALLOW_CREDENTIALS=true
DATABASE_URL=/var/lib/todo-api/app.db
```

### Docker Configuration
```bash
# .env for Docker
PORT=8080
HOST=0.0.0.0
PUBLIC_SCHEME=http
ENABLE_SWAGGER=false
ENABLE_LOGGER=true
ENABLE_RATE_LIMIT=false
GIN_MODE=release
TRUSTED_PROXIES=
ALLOWED_ORIGINS=http://localhost:3000
ALLOW_CREDENTIALS=true
DATABASE_URL=/data/app.db
```

### Testing Configuration
```bash
# .env for testing
PORT=8081
HOST=localhost
PUBLIC_SCHEME=http
ENABLE_SWAGGER=false
ENABLE_LOGGER=false
ENABLE_RATE_LIMIT=false
GIN_MODE=debug
TRUSTED_PROXIES=
ALLOWED_ORIGINS=*
ALLOW_CREDENTIALS=false
DATABASE_URL=:memory:
```

## Configuration Loading

### Load Order
1. **Default values** (hardcoded in application)
2. **Environment variables** (from system or .env file)
3. **Validation and parsing**
4. **Error handling** (application exits on invalid config)

### Environment Variable Priority
Environment variables take precedence over `.env` file values:

```bash
# .env file
PORT=8080

# Command line override
PORT=9000 ./golang-todo-api
# Application will use port 9000
```

### Validation
The application validates configuration at startup:

```bash
# Example validation errors
Invalid PORT: must be between 1 and 65535
Invalid GIN_MODE: must be 'debug' or 'release'
Invalid DATABASE_URL: directory does not exist
TRUSTED_PROXIES contains invalid CIDR: 192.168.1.300/24
```

## Dynamic Configuration

### Runtime Changes
Most configuration options require application restart. The following can be changed at runtime:

- **Log level** (through application signals)
- **Database connection pool** (through GORM settings)

### Configuration Reload
The application does not support configuration hot-reloading. To change configuration:

1. Update `.env` file
2. Restart application
3. Verify new configuration in startup logs

## Configuration in Different Environments

### Local Development
```bash
# Use .env file
cp .env.example .env
# Edit .env as needed
make run
```

### Docker Development
```bash
# Override in docker-compose.yml
services:
  api:
    environment:
      - PORT=8080
      - ENABLE_SWAGGER=true
```

### Kubernetes Deployment
```yaml
# ConfigMap
apiVersion: v1
kind: ConfigMap
metadata:
  name: todo-api-config
data:
  PORT: "80"
  GIN_MODE: "release"
  ENABLE_SWAGGER: "false"

# Pod spec
spec:
  containers:
  - name: api
    envFrom:
    - configMapRef:
        name: todo-api-config
```

### Docker Run
```bash
docker run \
  -e PORT=9000 \
  -e ENABLE_SWAGGER=true \
  -e DATABASE_URL=/data/app.db \
  -v ./data:/data \
  -p 9000:9000 \
  golang-todo-api:latest
```

## Security Considerations

### Sensitive Configuration
While this application doesn't currently have sensitive configuration (like API keys or passwords), follow these best practices:

#### Environment Variables
- **Never commit** `.env` files to version control
- **Use secrets management** in production (e.g., Kubernetes secrets)
- **Limit access** to environment configuration

#### Production Security
```bash
# Secure CORS configuration
ALLOWED_ORIGINS=https://yourapp.com  # Never use *

# Proper proxy configuration
TRUSTED_PROXIES=192.168.1.0/24  # Be specific

# Disable debug features
ENABLE_SWAGGER=false
GIN_MODE=release
```

## Configuration Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Error: address already in use
# Solution: Change port or kill existing process
PORT=8081 make run
```

#### 2. Database Directory Not Found
```bash
# Error: no such file or directory
# Solution: Create directory or use absolute path
mkdir -p data
# Or: DATABASE_URL=/absolute/path/to/app.db
```

#### 3. CORS Issues
```bash
# Error: CORS policy blocks request
# Solution: Add your frontend origin
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

#### 4. Invalid CIDR in TRUSTED_PROXIES
```bash
# Error: invalid CIDR address
# Solution: Use valid CIDR notation
TRUSTED_PROXIES=192.168.1.0/24  # Not 192.168.1.0/255
```

### Debugging Configuration

#### View Current Configuration
The application logs its configuration at startup (sensitive values are masked):

```bash
2023/01/01 12:00:00 Starting server with configuration:
2023/01/01 12:00:00 - Port: 8080
2023/01/01 12:00:00 - Host: 0.0.0.0
2023/01/01 12:00:00 - Mode: release
2023/01/01 12:00:00 - Swagger: disabled
2023/01/01 12:00:00 - Database: data/app.db
```

#### Configuration Validation
```bash
# Test configuration without starting server
go run cmd/server/main.go --validate-config
```

#### Environment Variable Debugging
```bash
# Print all environment variables
env | grep -E "(PORT|HOST|GIN_MODE|DATABASE_URL)"

# Test specific variable
echo $PORT
```

## Advanced Configuration

### Custom Configuration Path
```bash
# Use custom .env file location
ENV_FILE=/path/to/custom.env make run
```

### Configuration Profiles
```bash
# Development profile
ln -sf .env.development .env

# Production profile  
ln -sf .env.production .env

# Testing profile
ln -sf .env.testing .env
```

### Configuration Templating
For multiple environments, you can use environment variable substitution:

```bash
# .env.template
PORT=${TODO_API_PORT:-8080}
HOST=${TODO_API_HOST:-localhost}
DATABASE_URL=${TODO_API_DB_URL:-data/app.db}
```

### Monitoring Configuration Changes
```bash
# Watch for configuration file changes
fswatch .env | while read file; do echo "Config changed: $file"; done
```

## Best Practices

### Development
- Use `.env` file for local development
- Enable Swagger and logging for debugging
- Use debug mode for detailed error messages
- Keep database in local `data/` directory

### Testing
- Use in-memory database (`:memory:`)
- Disable logging to reduce noise
- Use different port to avoid conflicts
- Disable rate limiting for faster tests

### Production
- Use environment variables, not `.env` files
- Enable rate limiting for protection
- Use release mode for performance
- Store database in persistent volume
- Configure proper CORS and trusted proxies
- Disable Swagger for security