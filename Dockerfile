FROM golang:alpine
MAINTAINER Eric Youngberg <eyoungberg@mapc.org>

ENV GOPATH /go
ENV ARCH=linux/amd64

WORKDIR /opt

RUN set -e; \
    \
    apk add --no-cache \
        git \
        curl \
        make \
    ; \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
    ; \
    echo "dep ensure" > entrypoint.sh \
    && echo 'make static ARCH=${ARCH}' >> entrypoint.sh \
    && chmod +x entrypoint.sh

WORKDIR /go/src/github.com/prql/prql

ENTRYPOINT ["sh", "/opt/entrypoint.sh"]
