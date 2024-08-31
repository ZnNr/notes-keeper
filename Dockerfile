# Используем официальный образ Golang для базового образа
FROM golang:1.22.2

# Указываем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum, чтобы кешировать зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект в контейнер
COPY . .
COPY *.db ./
COPY cmd/*.go ./
# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o  /notes-keeper



# Открываем порт
EXPOSE 8080

# Команда для запуска приложения
CMD ["/notes-keeper"]