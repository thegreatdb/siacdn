FROM golang:1.9-alpine

WORKDIR /root

RUN \
  apk update && \
  apk add git ca-certificates curl openssl && \
  update-ca-certificates

ENV SIA_COMMIT ebe4195160b756eb43283bf8572dd96fa5ab1f1c

RUN go get -u github.com/ericflo/Sia/...

EXPOSE 9980 9981 9982

ADD dockerbin/*.sh /root/

ENTRYPOINT ["/root/bootstrap.sh"]
