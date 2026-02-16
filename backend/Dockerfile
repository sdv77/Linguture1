# backend/Dockerfile
# ============================================================================

# Этап 1: Сборка
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем ОДИН бинарник с именем 'app'
# Он будет поддерживать аргументы: ./app admin или ./app api
RUN go build -o app ./cmd/main.go

# Этап 2: Запуск
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/cmd/admin/admin_panel.html ./cmd/admin/

EXPOSE 8080 8081

CMD ["./app"]