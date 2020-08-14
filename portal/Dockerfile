FROM openresty/openresty:1.17.8.2-buster

#############################################
# STAGE 0: Build the portal with nodejs etc #
#############################################

ENV SIACDN_BASE_BUILD_VERSION 1

# Install base system
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get -y install ca-certificates curl git build-essential dnsutils && \
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    apt-get update && apt-get install -y yarn

ENV SIACDN_CORE_BUILD_VERSION 2

ENV GATSBY_API_URL "https://www.siacdn.com"

# Install and build portal
RUN git clone https://github.com/NebulousLabs/skynet-webportal.git && \
    cd skynet-webportal && \
    sed -i -- 's#'"siasky.net"'#'"www.siacdn.com"'#g' "gatsby-config.js" && \
    yarn && \
    yarn run build


#######################################################################
# STAGE 1: Build the portal image with custom config and start script #
#######################################################################

FROM openresty/openresty:1.17.8.2-buster

RUN apt-get update && apt-get -y install ca-certificates curl
RUN mkdir /etc/skynet/ && mkdir /tmp/nginx-cache
COPY --from=0 /skynet-webportal/public /var/www/webportal
COPY ./nginx/* /etc/nginx/
COPY ./nginx/conf.d/* /etc/nginx/conf.d/
COPY ./nginx/conf.d/include/* /etc/nginx/conf.d/include/
COPY ./start-portal.sh /etc/skynet/
RUN rm /etc/nginx/conf.d/default.conf
RUN chmod +x /etc/skynet/*.sh
# Use SIGQUIT instead of default SIGTERM to cleanly drain requests
# See https://github.com/openresty/docker-openresty/blob/master/README.md#tips--pitfalls
STOPSIGNAL SIGQUIT