#!/bin/bash

# ZarishSphere FHIR Engine - Validation Script
# This script validates FHIR resource files

set -e

if [ $# -eq 0 ]; then
    echo "🏥 ZarishSphere FHIR Engine - Resource Validator"
    echo "==============================================="
    echo ""
    echo "Usage: $0 <fhir-resource-file.json>"
    echo ""
    echo "Example:"
    echo "  $0 examples/patient.json"
    echo "  $0 test-data/observation.json"
    echo ""
    exit 1
fi

RESOURCE_FILE="$1"

if [ ! -f "$RESOURCE_FILE" ]; then
    echo "❌ File not found: $RESOURCE_FILE"
    exit 1
fi

echo "🏥 ZarishSphere FHIR Engine - Resource Validator"
echo "==============================================="
echo "📄 Validating: $RESOURCE_FILE"
echo ""

# Build the CLI tool
echo "📦 Building zs-core-fhir-engine CLI..."
go build -o fhir-engine ./cmd/fhir-engine

# Validate the resource
echo "🔍 Validating FHIR resource..."
./fhir-engine validate --file="$RESOURCE_FILE"

echo ""
echo "✅ Validation complete!"
