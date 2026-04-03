#!/bin/bash

# Build the CLI tool
echo "Building zs-core-fhir-engine CLI..."
go build -o fhir-engine ./cmd/fhir-engine

# Start the terminology server in the background
echo "Starting Terminology Server on port 8080..."
./fhir-engine terminology --port 8080 &

# Keep the script running
wait
