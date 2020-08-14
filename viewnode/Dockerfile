FROM ubuntu:20.04

#########################################
# STAGE 0: Install Golang and build Sia #
#########################################

ENV SIACDN_BASE_BUILD_VERSION 1

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get -y install ca-certificates curl git build-essential

ENV GOLANG_GOOS linux
ENV GOLANG_GOARCH amd64
ENV GOLANG_VERSION 1.14.5
ENV GOPATH $HOME/go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

RUN curl -sSL https://golang.org/dl/go$GOLANG_VERSION.$GOLANG_GOOS-$GOLANG_GOARCH.tar.gz \
  | tar -v -C /usr/local -xz

ENV SIACDN_SIA_BUILD_VERSION 12

RUN git clone https://gitlab.com/NebulousLabs/Sia.git && \
    cd Sia && \
#    git fetch --all --tags && \
#    git checkout tags/v1.4.10 && \
    make

################################################
# STAGE 1: Copy the built binary to production #
################################################

FROM ubuntu:20.04

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get -y install ca-certificates curl unzip

ENV GOLANG_GOOS linux
ENV GOLANG_GOARCH amd64
ENV GOLANG_VERSION 1.14.5
ENV GOPATH $HOME/go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

COPY --from=0 /go/bin/ /go/bin/
COPY *.sh /go/bin/

# TODO: Remove this
WORKDIR /root/.sia