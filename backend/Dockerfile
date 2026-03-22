# backend/Dockerfile
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

# 🔹 ВАЖНО: Второй этап с postgresql-client
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/cmd/admin/admin_panel.html ./cmd/admin/

EXPOSE 8080 8081

CMD ["./app"]