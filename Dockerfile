# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /zs-core-fhir-engine ./cmd/fhir-engine

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /zs-core-fhir-engine ./zs-core-fhir-engine

# Copy the runtime config
COPY --from=builder /app/config ./config

# Expose the server port
EXPOSE 8080

# Start the server by default
CMD ["./zs-core-fhir-engine", "serve", "-p", "8080", "-i", "./config"]
