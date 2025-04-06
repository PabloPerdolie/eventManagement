FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod ./
COPY go.sum ./

# Загружаем все зависимости
RUN go mod download

# Копируем исходники
COPY . .

# Компилируем приложение 
RUN CGO_ENABLED=0 GOOS=linux go build -a -o core-service .

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Устанавливаем зависимости
RUN apk --no-cache add ca-certificates tzdata

# Копируем исполняемый файл из builder
COPY --from=builder /app/core-service .
COPY --from=builder /app/.env .

# Экспортируем порт
EXPOSE 8080

# Команда запуска приложения
CMD ["./core-service"]
