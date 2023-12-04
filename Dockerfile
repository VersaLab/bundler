FROM golang:1.20-buster as builder

WORKDIR /app

COPY . .

RUN go mod download && go build -v -o bundler

FROM debian:buster-slim

WORKDIR /app

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bundler ./bundler

ENTRYPOINT ./bundler start --mode private
