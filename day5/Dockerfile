# Build Stage
FROM rust:latest as builder
WORKDIR /app
COPY . .
RUN cargo build --release

# Runtime Stage
FROM debian:bookworm-slim
WORKDIR /app
# Install runtime dependencies
RUN apt-get update && apt-get install -y libssl3 ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/target/release/rust-k8s-api /app/rust-k8s-api
EXPOSE 8080
CMD ["./rust-k8s-api"]
