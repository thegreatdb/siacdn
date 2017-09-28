FROM golang:1.9-alpine

WORKDIR /root

RUN \
  apk update && \
  apk add git ca-certificates curl openssl && \
  update-ca-certificates

ENV SIA_COMMIT 3c012755d14debc01bc0930143fb3cafc6a764ee

RUN go get -u github.com/ericflo/Sia/...

EXPOSE 9980 9981 9982

ADD dockerbin/*.sh /root/

ENTRYPOINT ["/root/bootstrap.sh"]
