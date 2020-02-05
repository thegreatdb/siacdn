FROM ubuntu:20.04

RUN apt-get update && apt-get upgrade -y && apt-get -y install gnupg2 ca-certificates curl && \
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    apt-get update && apt-get install -y yarn

RUN apt-get install -y git && \
    git clone https://github.com/NebulousLabs/skynet-webportal.git && \
    cd skynet-webportal && \
    yarn && \
    apt-get remove --purge -y git && \
    apt autoremove -y

WORKDIR /skynet-webportal