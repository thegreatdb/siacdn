FROM golang:1.9-alpine

WORKDIR /root

RUN \
  apk update && \
  apk add git ca-certificates curl openssl util-linux && \
  update-ca-certificates

ENV SIA_COMMIT 561a85d9cde1772ab3b84dae76c13e19c1809ed3

RUN go get -u github.com/NebulousLabs/Sia/...

EXPOSE 9980 9981 9982

ADD dockerbin/*.sh /root/

ENTRYPOINT ["/root/bootstrap.sh"]
