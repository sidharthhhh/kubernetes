FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
# Static build for scratch container
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Production vs Development Base Images:
# - FROM scratch (0 MB, no shell, max security) - Use for production
# - FROM alpine:latest (~5 MB, has /bin/sh) - Use for debugging with kubectl exec
# Changed from 'scratch' to 'alpine' to enable shell access for learning/debugging
FROM alpine:latest
COPY --from=builder /app/server /server
CMD ["/server"]
