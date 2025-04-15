FROM golang:1.23.4 AS builder
LABEL authors="Ionov Egor"

RUN apt-get update && \
    apt-get install -y sqlite3

WORKDIR /app

COPY go.mod go.sum ./
COPY ./pkg ./pkg
COPY main.go ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /to-do

COPY web /web

ENV TODO_PORT=7450
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=12345

EXPOSE 7540

CMD ["/to-do"]