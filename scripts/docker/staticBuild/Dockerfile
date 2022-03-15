FROM golang:1.17-alpine

LABEL org.opencontainers.image.authors="etho@ethohampton.com"

RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gofiber/adaptor/v2
RUN go get github.com/gofiber/fiber/v2
RUN apk add --no-cache build-base
