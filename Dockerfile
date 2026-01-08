# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-server ./cmd/api-server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/api-server .
COPY --from=builder /app/config ./config

# Expose port
EXPOSE 8080

# Set environment
ENV API_ENVIRONMENT=production

# Run the application
CMD ["./api-server"]