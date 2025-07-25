# Multi-stage build
FROM golang:1.23-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git ca-certificates tzdata

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение из cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Финальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл из builder stage
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"] 