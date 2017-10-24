FROM golang:1.9-alpine

WORKDIR /root

RUN \
  apk update && \
  apk add git ca-certificates curl openssl util-linux && \
  update-ca-certificates

ENV SIA_COMMIT 20403e3e26e47f82f9ac21fbb660cca1184b3a47

RUN go get -u github.com/ericflo/Sia/...

EXPOSE 9980 9981 9982

ADD dockerbin/*.sh /root/

ENTRYPOINT ["/root/bootstrap.sh"]
