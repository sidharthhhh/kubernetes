# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
# In a real scenario we would copy go.sum too, but keeping it simple for this lab
# COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run stage
FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
