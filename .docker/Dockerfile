# syntax=docker/dockerfile:1
FROM golang:1.18.1-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /conlangdev github.com/conlangdev/conlangdev/cmd/conlangdev

EXPOSE 8000

ENTRYPOINT ["/conlangdev"]