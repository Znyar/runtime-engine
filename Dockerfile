FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/web-api ./cmd/web-api/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /bin/web-api /app/web-api

CMD ["/app/web-api"]