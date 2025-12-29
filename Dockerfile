# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary main file executable
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Runtime stage (smaller image)
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migrations
COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080

CMD ["./main"]