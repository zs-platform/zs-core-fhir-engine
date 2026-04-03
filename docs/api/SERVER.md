# FHIR Server

The zs-core-fhir-engine FHIR server is a lightweight, in-memory RESTful server that implements the FHIR R5 specification.

## Starting the Server

```bash
# Basic usage
./zs-core-fhir --server --port 8080

# With custom IG path
./zs-core-fhir --server --port 8080 --ig ./config
```

## Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `--server` | `false` | Start the full FHIR server |
| `--term-server` | `false` | Start the terminology server only |
| `--port` | `8080` | Port to listen on |
| `--ig` | `./config` | Path to FHIR Implementation Guide |

## Server Features

### In-Memory Storage

The server stores all resources in memory. This is suitable for:
- Development and testing
- Small-scale deployments
- Demo purposes

For production use, consider extending the server with a persistent database.

### Implementation Guide Support

The server can load CodeSystems and ValueSets from a FHIR Implementation Guide:

```bash
./zs-core-fhir --server --ig /path/to/your/IG
```

### Thread Safety

The server uses read-write mutexes for thread-safe operations, making it safe for concurrent access.

## Running Both Servers

You can run both the FHIR server and terminology server on different ports:

```bash
# Terminal 1: FHIR Server
./zs-core-fhir --server --port 8080

# Terminal 2: Standalone Terminology Server
./zs-core-fhir --term-server --port 8081
```

## Docker Deployment

```bash
# Build
docker build -t zs-core-fhir .

# Run
docker run -p 8080:8080 zs-core-fhir --server --port 8080
```
