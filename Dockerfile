#FROM ubuntu:latest
#LABEL authors="Ilonochka"
#
#ENTRYPOINT ["top", "-b"]

#FROM golang:1.25-alpine AS builder
#WORKDIR /app
#COPY go.mod go.sum ./
#RUN go mod download
#COPY . .
#
#RUN go build -o subscription-service ./cmd/api



# ------------------------
# СТАДИЯ 1: Сборка приложения
# ------------------------
FROM golang:1.25-alpine AS builder

# Рабочая директория
WORKDIR /app

# Копируем файлы с зависимостями для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник API
RUN go build -o subscription-service ./cmd/api

# ------------------------
# СТАДИЯ 2: Минимальный образ для запуска
# ------------------------
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates postgresql-client

# Копируем бинарник из стадии сборки
COPY --from=builder /app/subscription-service .

# Устанавливаем переменные окружения (по желанию)
ENV PORT=8080

# Открываем порт API
EXPOSE 8080

# Команда запуска контейнера
ENTRYPOINT ["./subscription-service"]


