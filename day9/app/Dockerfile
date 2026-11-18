# Multi-stage Dockerfile for Day 9 Observability & Prometheus Demo
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy module definition and application code
COPY go.mod main.go ./

# Tidy dependencies and build statically linked binary for Linux
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# Final minimal runtime image with curl for debugging
FROM alpine:latest
RUN apk add --no-cache curl ca-certificates

WORKDIR /
COPY --from=builder /app/server /server

EXPOSE 8080

CMD ["/server"]
