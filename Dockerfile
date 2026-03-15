# Build stage
FROM golang:1.24.0 AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o antibruteforce ./cmd/antibruteforce

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/antibruteforce .
# Copy migrations if they are needed by the app at startup
COPY --from=builder /app/migrations ./migrations

# Expose gRPC port
EXPOSE 50051

# Run the application
CMD ["./antibruteforce"]
