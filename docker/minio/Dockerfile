FROM golang:1.9.1-alpine

ENV BUILD_VERSION f4c3078d9c63d3500796f6d3acc2b557b01297d4

RUN \
  apk update && \
  apk add git ca-certificates curl openssl build-base && \
  update-ca-certificates && \
  mkdir -p /go/src/github.com/NebulousLabs && \
  cd /go/src/github.com/NebulousLabs && \
  git clone https://github.com/NebulousLabs/Sia.git && \
  cd Sia && \
  go get -u github.com/NebulousLabs/demotemutex && \
  go get -u github.com/NebulousLabs/fastrand && \
  go get -u github.com/NebulousLabs/merkletree && \
  go get -u github.com/NebulousLabs/bolt && \
  go get -u golang.org/x/crypto/blake2b && \
  go get -u golang.org/x/crypto/ed25519 && \
  go get -u github.com/NebulousLabs/entropy-mnemonics && \
  go get -u github.com/NebulousLabs/errors && \
  go get -u github.com/NebulousLabs/go-upnp && \
  go get -u github.com/NebulousLabs/muxado && \
  go get -u github.com/klauspost/reedsolomon && \
  go get -u github.com/julienschmidt/httprouter && \
  go get -u github.com/inconshreveable/go-update && \
  go get -u github.com/kardianos/osext && \
  go get -u github.com/inconshreveable/mousetrap && \
  go get -u golang.org/x/crypto/ssh/terminal && \
  go get -u github.com/spf13/cobra/... && \
  go get -u github.com/boltdb/bolt && \
  mkdir -p /go/src/github.com/minio && \
  cd /go/src/github.com/minio && \
  git clone https://github.com/ericflo/minio.git && \
  cd /go/src/github.com/minio/minio && \
  make && \
  mv /go/src/github.com/minio/minio/minio /go/bin/ && \
  cd /root && \
  rm -Rf /go/src/github.com && \
  apk del build-base

RUN mkdir /minio
WORKDIR /minio

EXPOSE 9000

ENTRYPOINT ["/go/bin/minio"]
CMD ["gateway", "sia"]
