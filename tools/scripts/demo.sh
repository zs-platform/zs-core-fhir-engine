#!/bin/bash

# ZarishSphere FHIR Engine - Quick Demo
# This script demonstrates all the main features

echo "🏥 ZarishSphere FHIR Engine - Quick Demo"
echo "======================================="
echo ""

# Show version information
echo "📋 1. System Information:"
./zs-core-fhir-engine version
echo ""

# Test validation
echo "🔍 2. Patient Data Validation:"
./scripts/validate.sh examples/patient.json
echo ""

# Show available commands
echo "📚 3. Available Commands:"
echo "   • Development Server: ./scripts/dev_server.sh"
echo "   • Resource Validator: ./scripts/validate.sh <file.json>"
echo "   • Test Suite: ./scripts/test.sh"
echo "   • Terminology Server: ./scripts/run_fhir.sh"
echo ""

echo "🎯 Demo Complete!"
echo "📋 Your ZarishSphere FHIR Engine is ready for healthcare data management."
