# Этап сборки
FROM golang:1.21 AS builder

WORKDIR /app

# Копируем файлы go.mod и go.sum отдельно
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копируем все исходники (каждую папку с исходниками отдельно, для кеширования)
COPY ./cmd ./cmd

# Сборка бинарного файла
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o myapp ./cmd/tictactoe/main.go


# Облегченный образ для запуска
FROM alpine:latest

RUN apk add tzdata

WORKDIR /app

COPY --from=builder /app/myapp .

EXPOSE 8080

CMD ["./myapp"]
