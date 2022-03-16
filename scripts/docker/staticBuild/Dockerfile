FROM golang:1.18-alpine

LABEL org.opencontainers.image.authors="etho@ethohampton.com"

COPY ./go.mod /tmp/install/go.mod
COPY ./go.sum /tmp/install/go.sum
COPY ./fakeDeps.go /tmp/install/fakeDeps.go

RUN cd /tmp/install && go mod download && go build fakeDeps.go

#remove the temporary go files
RUN rm -rf /tmp/install

RUN apk add --no-cache build-base
