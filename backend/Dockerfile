FROM golang:1.9-alpine

RUN \
  apk update && \
  apk add git ca-certificates curl openssl && \
  update-ca-certificates

RUN mkdir /siacdn
WORKDIR /siacdn

ADD . /go/src/github.com/thegreatdb/siacdn/backend

RUN cd /go/src/github.com/thegreatdb/siacdn/backend && \
  go build -o /siacdn/siacdn main.go && \
  rm -Rf /go/src/*

ENTRYPOINT ["/siacdn/siacdn"]
