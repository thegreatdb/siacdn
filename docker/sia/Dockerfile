FROM golang:1.9-alpine

WORKDIR /root

RUN \
  apk update && \
  apk add git ca-certificates curl openssl util-linux && \
  update-ca-certificates

ENV SIA_COMMIT e1b6cca25e783ec64ba9cf370b6119cc112842de

RUN go get -u github.com/NebulousLabs/Sia/...

EXPOSE 9980 9981 9982

ADD dockerbin/*.sh /root/

ENTRYPOINT ["/root/bootstrap.sh"]
