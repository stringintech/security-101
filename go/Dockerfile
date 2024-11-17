# Build stage
FROM golang:1.23.3-alpine3.19 AS builder

# Install necessary build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server

# Final stage
FROM alpine:3.19

# Add non root user
RUN adduser -D -g '' appuser

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Use non root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
