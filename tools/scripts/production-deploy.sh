#!/bin/bash

# ZarishSphere FHIR Engine - Production Deployment Script
# This script sets up the complete production-ready FHIR engine

set -e

echo "🏥 ZarishSphere FHIR Engine - Production Setup"
echo "======================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Step 1: Dependencies
echo -e "${BLUE}1. Dependencies${NC}"
echo "   Checking Go modules..."
if ! command -v go &> /dev/null; then
    echo -e "   ${RED}❌ Go not installed${NC}"
    exit 1
fi

echo "   ${GREEN}✅ Go is installed${NC}"
echo ""

# Step 2: Database Setup
echo -e "${BLUE}2. Database Setup${NC}"
echo "   Checking PostgreSQL connection..."
if command -v psql &> /dev/null; then
    echo -e "   ${YELLOW}⚠️  PostgreSQL not found - using in-memory mode${NC}"
    DB_MODE="memory"
else
    echo "   ${GREEN}✅ PostgreSQL is available${NC}"
    DB_MODE="postgresql"
fi
echo ""

# Step 3: Build Application
echo -e "${BLUE}3. Build Application${NC}"
echo "   Building FHIR engine..."
go mod tidy
go build -o zs-core-fhir-engine ./cmd/zs-core-fhir-engine

if [ ! -f "./zs-core-fhir-engine" ]; then
    echo -e "   ${RED}❌ Build failed${NC}"
    exit 1
fi

echo "   ${GREEN}✅ Build successful${NC}"
echo ""

# Step 4: Configuration
echo -e "${BLUE}4. Configuration${NC}"
echo "   Setting up production configuration..."

# Create production config
cat > production.env << EOF
# ZarishSphere FHIR Engine - Production Configuration
FHIR_SERVER_PORT=8080
FHIR_SERVER_HOST=0.0.0.0
FHIR_LOG_LEVEL=info
FHIR_DB_MODE=$DB_MODE
FHIR_TENANT_ID=default
FHIR_IG_PATH=./BD-Core-FHIR-IG
FHIR_TLS_ENABLED=false
FHIR_CERT_FILE=
FHIR_KEY_FILE=
FHIR_METRICS_ENABLED=true
FHIR_AUDIT_ENABLED=true
EOF

echo "   ${GREEN}✅ Configuration created${NC}"
echo ""

# Step 5: Start Services
echo -e "${BLUE}5. Start Services${NC}"
echo "   Starting FHIR engine with production configuration..."

# Set environment variables
export FHIR_SERVER_PORT=${FHIR_SERVER_PORT:-8080}
export FHIR_SERVER_HOST=${FHIR_SERVER_HOST:-0.0.0.0}
export FHIR_LOG_LEVEL FHIR_DB_MODE FHIR_TENANT_ID
export FHIR_IG_PATH FHIR_TLS_ENABLED FHIR_CERT_FILE FHIR_KEY_FILE
export FHIR_METRICS_ENABLED FHIR_AUDIT_ENABLED

# Start the server
if [ "$DB_MODE" = "postgresql" ]; then
    echo -e "   ${GREEN}🚀 Starting with PostgreSQL backend...${NC}"
    ./zs-core-fhir-engine serve --port "$FHIR_SERVER_PORT" --ig-path "$FHIR_IG_PATH"
else
    echo -e "   ${YELLOW}🚀 Starting with in-memory backend...${NC}"
    ./zs-core-fhir-engine serve --port "$FHIR_SERVER_PORT" --ig-path "$FHIR_IG_PATH"
fi

echo ""
echo -e "${GREEN}6. Health Check${NC}"
echo "   Verifying server health..."

# Wait for server to start
sleep 3

# Health check
if curl -s http://localhost:$FHIR_SERVER_PORT/healthz &> /dev/null; then
    echo -e "   ${GREEN}✅ Server is healthy and responding${NC}"
    echo -e "   ${BLUE}📋 Server running on: http://localhost:$FHIR_SERVER_PORT${NC}"
    echo -e "   ${BLUE}📋 FHIR endpoints available at: http://localhost:$FHIR_SERVER_PORT/fhir${NC}"
    echo -e "   ${BLUE}📋 Metadata endpoint: http://localhost:$FHIR_SERVER_PORT/fhir/metadata${NC}"
else
    echo -e "   ${RED}❌ Server health check failed${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 Production Deployment Complete!${NC}"
echo "======================================================="
echo -e "${BLUE}Summary:${NC}"
echo -e "   ${GREEN}✅${NC} Dependencies: Ready"
echo -e "   ${GREEN}✅${NC} Database: $([ "$DB_MODE" = "postgresql" ] && echo "PostgreSQL" || echo "In-Memory")"
echo -e "   ${GREEN}✅${NC} Build: Successful"
echo -e "   ${GREEN}✅${NC} Server: Running on port $FHIR_SERVER_PORT"
echo ""
echo -e "${BLUE}Next Steps:${NC}"
echo "   1. Configure PostgreSQL connection details"
echo "   2. Set up reverse proxy (nginx/traefik)"
echo "   3. Enable TLS certificates"
echo "   4. Configure monitoring (Prometheus/Grafana)"
echo ""
echo -e "${YELLOW}📚 Production Guide: See docs/production-deployment.md${NC}"
