FROM ubuntu:20.04

# The domain we're hosting on
ENV SIACDN_DOMAIN siacdn.com

# For cache invalidation
ENV SIACDN_BUILD_VERSION 1

# Install base system
RUN apt-get update && apt-get upgrade -y && apt-get -y install gnupg2 ca-certificates curl && \
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    apt-get update && apt-get install -y yarn

# Install and build portal and server
RUN apt-get install -y git && \
    git clone https://github.com/NebulousLabs/skynet-webportal.git && \
    cd skynet-webportal && \
    sed -i -- 's#'"siasky.net"'#'"$SIACDN_DOMAIN"'#g' "packages/sia-skynet/package.json" && \
    yarn global add http-server && \
    yarn && \
    yarn workspace sia-skynet build && \
    yarn workspace siaviewnode-server build && \
    apt-get remove --purge -y git && \
    apt autoremove -y

WORKDIR /skynet-webportal