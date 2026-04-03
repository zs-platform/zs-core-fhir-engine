#!/bin/bash

# ZarishSphere FHIR Engine - Development Start Script
# This script starts the FHIR server for development

set -e

echo "🏥 ZarishSphere FHIR Engine - Development Server"
echo "================================================"

# Parse command line arguments
PORT=8080
DEBUG=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--port)
            PORT="$2"
            shift 2
            ;;
        --debug)
            DEBUG="--debug"
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Build the CLI tool
echo "📦 Building zs-core-fhir-engine CLI..."
go build -o zs-core-fhir-engine ./cmd/zs-core-fhir-engine

# Check if build was successful
if [ ! -f "./zs-core-fhir-engine" ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"

# Start the FHIR server
echo "🚀 Starting FHIR Server on port $PORT..."
echo "📋 Server will be available at: http://localhost:$PORT"
echo "📋 FHIR endpoints:"
echo "   - Health: http://localhost:$PORT/healthz"
echo "   - Metadata: http://localhost:$PORT/fhir/metadata"
echo "   - ValueSet Expand: http://localhost:$PORT/fhir/ValueSet/\$expand"
echo ""
echo "🛑 Press Ctrl+C to stop the server"
echo ""

# Start the server
./zs-core-fhir-engine serve -port $PORT $DEBUG
