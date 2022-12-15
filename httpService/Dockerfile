# Используем официальный образ Go для сборки
FROM golang:1.25-alpine AS builder

# Создаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы проекта в контейнер
COPY .. .

# Скачиваем зависимости и собираем бинарник
RUN go mod download
RUN go build -o myservice .

# Используем минимальный образ для финального контейнера
FROM alpine:latest

# Копируем бинарник из builder-образа
COPY --from=builder /app/myservice /usr/local/bin/myservice

# Указываем команду запуска сервиса
CMD ["myservice"]

# Открываем порт (если ваш сервис слушает, например, 8080)
EXPOSE 8080