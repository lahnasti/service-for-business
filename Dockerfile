FROM golang:1.22.4

# Установка рабочей директории
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./src/cmd/main.go

COPY src/migrations /app/migrations

# Определяем команду для запуска приложения
CMD ["/wait-for-it.sh", "db:5432", "--", "./main"]