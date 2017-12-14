FROM golang:1.9-alpine

WORKDIR /root

RUN \
  apk update && \
  apk add git ca-certificates curl openssl util-linux && \
  update-ca-certificates

ENV SIA_COMMIT ef89f245eb80273d92dcd38b10af4b53bf80b52b

RUN go get -u github.com/NebulousLabs/Sia/...

EXPOSE 9980 9981 9982

ADD dockerbin/*.sh /root/

ENTRYPOINT ["/root/bootstrap.sh"]
