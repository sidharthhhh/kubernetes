# Multi-stage Dockerfile for Day 11 Autoscaling Demo
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod main.go ./

# Compile statically linked binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# Final minimal runtime image with curl for debugging
FROM alpine:latest
RUN apk add --no-cache curl ca-certificates

WORKDIR /
COPY --from=builder /app/server /server

EXPOSE 8080

CMD ["/server"]
