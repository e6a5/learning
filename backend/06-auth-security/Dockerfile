# syntax=docker/dockerfile:1

FROM golang:1.23.4-alpine3.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o app .

EXPOSE 8080

CMD ["./app"] 