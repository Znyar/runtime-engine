FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/web-api ./cmd/web-api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/web-api /app/web-api

RUN apk add --no-cache curl

RUN chmod +x /app/web-api

CMD ["/app/web-api"]