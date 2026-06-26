# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=1 is required for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Install certificates
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy environment example
COPY .env.example .env

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]