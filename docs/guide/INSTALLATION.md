# Installation Guide

> **Version:** 2.0.0 | **Standards:** Go 1.26.1 (ADR-0001), PostgreSQL 18 (ADR-0003)

This guide will help you install and set up the ZarishSphere FHIR Engine for development and production use.

---

## 🚀 Quick Install

### Prerequisites

- **Go 1.26.1+** (ADR-0001) - [Install Go](https://golang.org/dl/)
- **PostgreSQL 18.3+** (ADR-0003) - [Install PostgreSQL](https://www.postgresql.org/download/)
- **NATS 2.12.5+** (ADR-0004) - For event streaming
- **Git** - [Install Git](https://git-scm.com/downloads)

### One-Command Install

```bash
# Clone and build
git clone https://github.com/zarishsphere/zs-core-fhir-engine.git
cd zs-core-fhir-engine
go build -o fhir-engine ./cmd/fhir-engine

# Start the server
./fhir-engine serve --port 8080
```

## 📦 System Requirements

### Minimum Requirements
- **CPU**: 2 cores
- **Memory**: 4GB RAM
- **Storage**: 10GB available
- **OS**: Linux, macOS, or Windows

### Recommended for Production (per ADR standards)

- **CPU**: 4+ cores
- **Memory**: 8GB+ RAM
- **Storage**: 50GB+ SSD
- **Database**: PostgreSQL 18.3+ with RLS (ADR-0003)
- **Cache**: Valkey 9.0.3 (ADR-0005)
- **Messaging**: NATS 2.12.5 JetStream (ADR-0004)
- **OS**: Linux (Ubuntu 22.04 LTS recommended)

## 🐧 Linux Installation

### Ubuntu/Debian

```bash
# Install Go
wget https://go.dev/dl/go1.26.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Clone and build
git clone https://github.com/zarishsphere/zs-core-fhir-engine.git
cd zs-core-fhir-engine
go build -o fhir-engine ./cmd/fhir-engine

# Setup database (PostgreSQL with RLS per ADR-0003)
sudo -u postgres createdb fhir_db
sudo -u postgres psql -c "CREATE USER fhir_user WITH PASSWORD 'secure_password';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE fhir_db TO fhir_user;"
```

### CentOS/RHEL

```bash
# Install Go
wget https://go.dev/dl/go1.26.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install PostgreSQL
sudo yum install postgresql-server postgresql-contrib
sudo postgresql-setup initdb
sudo systemctl enable postgresql
sudo systemctl start postgresql

# Clone and build
git clone https://github.com/zarishsphere/zs-core-fhir-engine.git
cd zs-core-fhir-engine
go build -o zs-core-fhir-engine ./cmd/zs-core-fhir-engine
```

## 🍎 macOS Installation

```bash
# Install Go (using Homebrew)
brew install go

# Install PostgreSQL (optional)
brew install postgresql
brew services start postgresql

# Clone and build
git clone https://github.com/zarishsphere/zs-core-fhir-engine.git
cd zs-core-fhir-engine
go build -o zs-core-fhir-engine ./cmd/zs-core-fhir-engine

# Setup database (optional)
createdb fhir_db
createuser fhir_user
psql -c "ALTER USER fhir_user PASSWORD 'secure_password';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE fhir_db TO fhir_user;"
```

## 🪟 Windows Installation

### Using Chocolatey

```powershell
# Install Go
choco install golang

# Install PostgreSQL (optional)
choco install postgresql

# Clone and build
git clone https://github.com/zarishsphere/zs-core-fhir-engine.git
cd zs-core-fhir-engine
go build -o zs-core-fhir-engine.exe ./cmd/zs-core-fhir-engine

# Setup database (optional)
createdb fhir_db
createuser fhir_user
psql -U postgres -c "ALTER USER fhir_user PASSWORD 'secure_password';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE fhir_db TO fhir_user;"
```

### Manual Installation

1. Download and install [Go 1.26.1+](https://golang.org/dl/)
2. Install [PostgreSQL 18+](https://www.postgresql.org/download/windows/) (optional)
3. Install [Git](https://git-scm.com/download/win)
4. Open Command Prompt or PowerShell and run:

```powershell
git clone https://github.com/zarishsphere/zs-core-fhir-engine.git
cd zs-core-fhir-engine
go build -o zs-core-fhir-engine.exe ./cmd/zs-core-fhir-engine
```

## 🐳 Docker Installation

### Using Docker Hub

```bash
# Pull the image
docker pull zarishsphere/zs-core-fhir-engine:latest

# Run the server
docker run -p 8080:8080 zarishsphere/zs-core-fhir-engine:latest

# Run with PostgreSQL
docker run -p 8080:8080 \
  -e FHIR_DB_MODE=postgresql \
  -e FHIR_DB_HOST=postgres \
  --link postgres:postgres \
  zarishsphere/zs-core-fhir-engine:latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/zarishsphere/zs-core-fhir-engine.git
cd zs-core-fhir-engine

# Build Docker image
docker build -t zs-core-fhir-engine .

# Run the container
docker run -p 8080:8080 zs-core-fhir-engine
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'
services:
  fhir-server:
    build: .
    ports:
      - "8080:8080"
    environment:
      - FHIR_DB_MODE=postgresql
      - FHIR_DB_HOST=postgres
      - FHIR_DB_USER=fhir_user
      - FHIR_DB_PASSWORD=secure_password
      - FHIR_DB_NAME=fhir_db
    depends_on:
      - postgres
    volumes:
      - ./examples:/app/examples

  postgres:
    image: postgres:18
    environment:
      - POSTGRES_DB=fhir_db
      - POSTGRES_USER=fhir_user
      - POSTGRES_PASSWORD=secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  postgres_data:
```

Start the services:

```bash
docker-compose up -d
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file for configuration:

```bash
# Server Configuration
FHIR_SERVER_PORT=8080
FHIR_SERVER_HOST=0.0.0.0
FHIR_LOG_LEVEL=info

# Database Configuration (optional)
FHIR_DB_MODE=postgresql
FHIR_DB_HOST=localhost
FHIR_DB_PORT=5432
FHIR_DB_NAME=fhir_db
FHIR_DB_USER=fhir_user
FHIR_DB_PASSWORD=secure_password

# Multi-tenancy (optional)
FHIR_TENANT_ID=default

# Security (optional)
FHIR_TLS_ENABLED=false
FHIR_CERT_FILE=/path/to/cert.pem
FHIR_KEY_FILE=/path/to/key.pem
```

### Database Setup

#### PostgreSQL Setup

```bash
# Create database
sudo -u postgres createdb fhir_db

# Create user
sudo -u postgres createuser fhir_user

# Set password
sudo -u postgres psql -c "ALTER USER fhir_user PASSWORD 'secure_password';"

# Grant privileges
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE fhir_db TO fhir_user;"

# Run migrations
./zs-core-fhir-engine migrate --database fhir_db --user fhir_user --password secure_password
```

#### In-Memory Mode (for testing)

```bash
# Use in-memory storage (no database required)
./zs-core-fhir-engine serve --port 8080
```

## 🧪 Verification

### Test the Installation

```bash
# Check if binary works
./zs-core-fhir-engine --help

# Start the server
./zs-core-fhir-engine serve --port 8080

# Test health endpoint (in another terminal)
curl http://localhost:8080/healthz

# Test FHIR endpoint
curl http://localhost:8080/fhir/metadata
```

### Expected Output

Health check should return:
```json
{
  "status": "ok",
  "timestamp": "2026-04-03T10:00:01Z",
  "uptime_seconds": 1,
  "version": "0.4.0-alpha"
}
```

## 🔧 Development Setup

### Install Dependencies

```bash
# Download Go modules
go mod download

# Run tests
go test -v ./...

# Build with debug info
go build -gcflags="all=-N -l" -o zs-core-fhir-engine ./cmd/zs-core-fhir-engine

# Install development tools
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Development Server

```bash
# Start with hot reload (using air)
air -c .air.toml

# Or start manually with debug
./zs-core-fhir-engine serve --debug --log-level debug --port 8080
```

### IDE Setup

#### VS Code

Install these extensions:
- Go (golang.go)
- Docker (ms-vscode-remote.remote-containers)
- GitLens (eamodio.gitlens)

#### GoLand

1. Open the project directory
2. Go to File → Settings → Go
3. Set GOPATH and GOROOT
4. Enable Go Modules integration

## 🚀 Production Deployment

### Systemd Service (Linux)

Create `/etc/systemd/system/zs-fhir.service`:

```ini
[Unit]
Description=ZarishSphere FHIR Engine
After=network.target postgresql.service

[Service]
Type=simple
User=fhir
Group=fhir
WorkingDirectory=/opt/zs-core-fhir-engine
ExecStart=/opt/zs-core-fhir-engine/zs-core-fhir-engine serve --port 8080
Restart=always
RestartSec=5
Environment=FHIR_DB_MODE=postgresql
Environment=FHIR_DB_HOST=localhost
Environment=FHIR_DB_USER=fhir_user
Environment=FHIR_DB_PASSWORD=secure_password
Environment=FHIR_DB_NAME=fhir_db

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable zs-fhir
sudo systemctl start zs-fhir
sudo systemctl status zs-fhir
```

### Nginx Reverse Proxy

Create `/etc/nginx/sites-available/zs-fhir`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /healthz {
        proxy_pass http://localhost:8080/healthz;
        access_log off;
    }
}
```

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/zs-fhir /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 🔍 Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Find what's using port 8080
sudo lsof -i :8080

# Kill the process
sudo kill -9 <PID>

# Or use a different port
./zs-core-fhir-engine serve --port 8081
```

#### Database Connection Failed
```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Test connection
psql -h localhost -U fhir_user -d fhir_db

# Use in-memory mode for testing
./zs-core-fhir-engine serve --port 8080
```

#### Build Errors
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Rebuild
go build -o zs-core-fhir-engine ./cmd/zs-core-fhir-engine
```

#### Permission Denied
```bash
# Make binary executable
chmod +x zs-core-fhir-engine

# Or build with correct permissions
go build -o zs-core-fhir-engine -mod=mod ./cmd/zs-core-fhir-engine
```

### Logs and Debugging

```bash
# Run with debug logging
./zs-core-fhir-engine serve --debug --log-level debug

# Check system logs (Linux)
sudo journalctl -u zs-fhir -f

# Check application logs
tail -f /var/log/zs-fhir/app.log
```

## 📚 Next Steps

- [Quick Start Guide](/guide/quickstart) - Get running in 5 minutes
- [API Reference](/api/overview) - Complete API documentation
- [Production Deployment](/guide/complete-server-blueprint) - Production setup
- [Configuration Guide](/guide/configuration) - Advanced configuration

---

**Need help?** Check our [Troubleshooting Guide](/guide/troubleshooting) or open an issue on GitHub.
