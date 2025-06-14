#Build stage

FROM golang:1.23.10-alpine3.21 AS builder
#Установление рабочей директорий в контейнере
WORKDIR /app 
#Скопировать .(все файлы нынешней директорий) .(в главную директорию контейнера)
COPY . .
#Компиляция приложения
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz
        


#Run stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY starter.sh .
COPY db/migration ./migration
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/starter.sh" ]