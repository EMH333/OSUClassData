FROM docker.io/library/golang:1.18-alpine

LABEL org.opencontainers.image.authors="etho@ethohampton.com"

COPY ./go.mod ./go.sum ./fakeDeps.go /tmp/install/

RUN cd /tmp/install && go mod download && go build fakeDeps.go

#remove the temporary go files
RUN rm -rf /tmp/install

RUN apk add --no-cache build-base
