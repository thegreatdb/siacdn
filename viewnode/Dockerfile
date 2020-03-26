FROM ubuntu:20.04

# For cache invalidation
ENV SIACDN_BUILD_VERSION 6

RUN apt-get update && apt-get upgrade -y && apt-get -y install ca-certificates curl

ENV GOLANG_GOOS linux
ENV GOLANG_GOARCH amd64
ENV GOLANG_VERSION 1.14
RUN curl -sSL https://golang.org/dl/go$GOLANG_VERSION.$GOLANG_GOOS-$GOLANG_GOARCH.tar.gz \
  | tar -v -C /usr/local -xz
ENV GOPATH $HOME/go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

RUN apt-get install -y git build-essential && \
    git clone https://gitlab.com/NebulousLabs/Sia.git && \
    cd Sia && \
    make && \
    apt-get remove --purge -y git build-essential && \
    apt autoremove -y

WORKDIR /root/.sia