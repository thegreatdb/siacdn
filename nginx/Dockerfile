FROM nginx:1.19

#############################################
# STAGE 0: Build the portal with nodejs etc #
#############################################

# For cache invalidation
ENV SIACDN_BUILD_VERSION 9

# Install base system
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get -y install gnupg2 ca-certificates curl git && \
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    apt-get update && apt-get install -y yarn

RUN mkdir /etc/skynet/

ENV GATSBY_API_URL "https://www.siacdn.com"

# Install and build portal and server
RUN git clone https://github.com/NebulousLabs/skynet-webportal.git && \
    cd skynet-webportal && \
    sed -i -- 's#'"siasky.net"'#'"www.siacdn.com"'#g' "gatsby-config.js" && \
    yarn && \
    yarn run build


######################################################################
# STAGE 1: Build the nginx image with custom config and start script #
######################################################################

FROM nginx:1.19

RUN apt-get update && apt-get -y install ca-certificates curl
COPY --from=0 /skynet-webportal/public /skynet-webportal/public
COPY . /etc/skynet/
RUN chmod +x /etc/skynet/*.sh