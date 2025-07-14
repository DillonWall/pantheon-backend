# Use official Go image as builder
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o pantheon-backend ./server.go

# Use lightweight alpine image for runtime
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/pantheon-backend .

# Expose port (match your server's port)
EXPOSE 8080

# Run the application
CMD ["./pantheon-backend"]
