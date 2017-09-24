FROM golang:1.9-alpine

WORKDIR /root

RUN \
  apk update && \
  apk add git ca-certificates curl openssl && \
  update-ca-certificates

ENV SIA_COMMIT 74896a42aa656afba973fc70d062183b9211ce8b

RUN go get -u github.com/NebulousLabs/Sia/...

EXPOSE 9980 9981 9982

ADD dockerbin/*.sh /root/

ENTRYPOINT ["/root/bootstrap.sh"]
