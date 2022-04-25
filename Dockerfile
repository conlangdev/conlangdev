# syntax=docker/dockerfile:1
FROM golang:1.18.1-alpine

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /conlangdev github.com/conlangdev/conlangdev/cmd/conlangdev

EXPOSE 8000

ENTRYPOINT ["/conlangdev", "run"]