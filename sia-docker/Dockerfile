FROM golang:1.9-alpine

WORKDIR /root

RUN \
  apk update && \
  apk add git ca-certificates curl openssl && \
  update-ca-certificates

ENV SIA_COMMIT b07148b60297f705495dae84e5e9de2fe889fb39

RUN go get -u github.com/ericflo/Sia/...

EXPOSE 9980 9981 9982

ADD dockerbin/*.sh /root/

ENTRYPOINT ["/root/bootstrap.sh"]
