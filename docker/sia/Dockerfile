FROM golang:1.9-alpine

WORKDIR /root

RUN \
  apk update && \
  apk add git ca-certificates curl openssl util-linux && \
  update-ca-certificates

ENV SIA_COMMIT 6da299c28746970d444baeb1a3bdcc90f4e8b6db-v3

RUN go get -u github.com/ericflo/Sia/...

EXPOSE 9980 9981 9982

ADD dockerbin/*.sh /root/

ENTRYPOINT ["/root/bootstrap.sh"]
