# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy everything including vendor
COPY . .

# Build the application with vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
