FROM golang:alpine

RUN apk add --no-cache \
    git \
    gcc \
    libc-dev \
    readline-dev

RUN mkdir /app
WORKDIR /app

COPY go.mod /app
RUN go mod download
RUN go get -u github.com/nathany/looper
RUN go get -u golang.org/x/tools/cmd/goimports 
