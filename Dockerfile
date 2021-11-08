FROM golang:alpine

RUN apk add --no-cache \
    gcc \
    libc-dev \
    readline-dev

RUN mkdir /app
WORKDIR /app

COPY go.mod /app
RUN go mod download
RUN go get -u github.com/nathany/looper
