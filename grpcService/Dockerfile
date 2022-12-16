# Используем официальный образ Go для сборки
FROM golang:1.25-alpine AS builder

# Создаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы проекта в контейнер
COPY .. .

# Скачиваем зависимости и собираем бинарник
RUN go mod download
RUN go build -o grpcservice .

# Используем минимальный образ для финального контейнера
FROM alpine:latest

# Копируем бинарник из builder-образа
COPY --from=builder /app/grpcservice /usr/local/bin/grpcservice

# Указываем команду запуска сервиса
CMD ["grpcservice"]

# Открываем порт (если ваш сервис слушает, например, 8080)
EXPOSE 50051