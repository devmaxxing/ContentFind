# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Install required build tools
RUN apk add --no-cache git

# Copy source code (relative to api-server directory)
COPY . .

# Download dependencies
RUN go mod download

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-server .

# Final stage
FROM scratch
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api-server .

# Set up non-root user
USER 1000:1000

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api-server"]
