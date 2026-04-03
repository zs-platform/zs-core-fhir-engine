# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /zs-core-fhir ./cmd/zs-core-fhir-engine

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /zs-core-fhir .

# Copy the IG data
COPY --from=builder /app/BD-Core-FHIR-IG ./BD-Core-FHIR-IG

# Expose the server port
EXPOSE 8080

# Start the server by default
CMD ["./zs-core-fhir", "serve", "-port", "8080", "-ig", "./BD-Core-FHIR-IG"]
