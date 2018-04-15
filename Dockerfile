FROM golang:alpine
MAINTAINER Eric Youngberg <eyoungberg@mapc.org>

ENV GOPATH /go
ENV ARCH=linux/amd64
WORKDIR /go/src/github.com/prql/prql

RUN set -ex; \
    \
    apk add --no-cache \
        git \
        make

CMD make static ARCH=${ARCH}
