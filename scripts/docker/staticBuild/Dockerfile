FROM docker.io/library/golang:1.23-alpine

LABEL org.opencontainers.image.authors="etho@ethohampton.com"

RUN apk add --no-cache build-base

COPY ./go.mod ./go.sum ./fakeDeps.go /tmp/install/

RUN cd /tmp/install && go mod download && go build -v fakeDeps.go

#remove the temporary go files
RUN rm -rf /tmp/install
