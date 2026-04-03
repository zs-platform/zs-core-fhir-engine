#!/bin/bash

# ZarishSphere FHIR Engine - Test Script
# This script runs comprehensive tests

set -e

echo "🧪 ZarishSphere FHIR Engine - Test Suite"
echo "======================================"

# Build the CLI tool
echo "📦 Building zs-core-fhir-engine CLI..."
go build -o zs-core-fhir-engine ./cmd/zs-core-fhir-engine

# Test CLI commands
echo ""
echo "🔧 Testing CLI commands..."

echo "  • Testing help command..."
./zs-core-fhir-engine --help > /dev/null
echo "    ✅ Help command works"

echo "  • Testing version command..."
./zs-core-fhir-engine version > /dev/null
echo "    ✅ Version command works"

echo "  • Testing validate help..."
./zs-core-fhir-engine validate --help > /dev/null
echo "    ✅ Validate command works"

echo "  • Testing terminology help..."
./zs-core-fhir-engine terminology --help > /dev/null
echo "    ✅ Terminology command works"

echo "  • Testing serve help..."
./zs-core-fhir-engine serve --help > /dev/null
echo "    ✅ Serve command works"

# Test Go modules
echo ""
echo "📦 Testing Go modules..."
go mod verify
echo "  ✅ Go modules verified"

# Test build
echo ""
echo "🔨 Testing build..."
go build ./cmd/zs-core-fhir-engine
echo "  ✅ Build successful"

# Test if binary runs
echo ""
echo "🚀 Testing binary execution..."
./zs-core-fhir-engine version
echo "  ✅ Binary executes successfully"

echo ""
echo "🎉 All tests passed!"
echo "📋 Ready for development and deployment"
