FROM ubuntu:20.04

####################################################
# STAGE 0: Install Golang and build statscollector #
####################################################

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

COPY . /etc/statscollector
WORKDIR /etc/statscollector
RUN go build -o /usr/local/go/bin/statscollector maxint.co/siacdn-statscollector

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

COPY --from=0 /usr/local/go/bin/statscollector /usr/local/go/bin/statscollector

CMD [ "/usr/local/go/bin/statscollector" ]