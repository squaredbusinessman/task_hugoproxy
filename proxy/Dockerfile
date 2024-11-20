# Используем официальный образ Go как базовый
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем оставшиеся исходники
COPY . .

# Собираем приложение
RUN go build -o main

FROM alpine:latest
COPY --from=builder /app/main /main

EXPOSE 8080
CMD ["/main"]
