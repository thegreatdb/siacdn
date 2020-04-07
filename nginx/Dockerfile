FROM ubuntu:20.04

#############################################
# STAGE 0: Build the portal with nodejs etc #
#############################################

# For cache invalidation
ENV SIACDN_BUILD_VERSION 6

# The domain we're hosting on
ENV SIACDN_DOMAIN www.siacdn.com
ENV GATSBY_API_URL "https://${SIACDN_DOMAIN}"

# Install base system
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get -y install gnupg2 ca-certificates curl git && \
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    apt-get update && apt-get install -y yarn

# Install and build portal and server
RUN git clone https://github.com/NebulousLabs/skynet-webportal.git && \
    cd skynet-webportal && \
    sed -i -- 's#'"siasky.net"'#'"$SIACDN_DOMAIN"'#g' "gatsby-config.js" && \
    yarn && \
    yarn run build

######################################################################
# STAGE 1: Build the nginx image with custom config and start script #
######################################################################

FROM nginx:1.17

COPY --from=0 /skynet-webportal/public /skynet-webportal/public

# For cache invalidation
ENV SIACDN_BUILD_VERSION 6

RUN mkdir /etc/skynet/

COPY . /etc/skynet/
RUN chmod +x /etc/skynet/*.sh