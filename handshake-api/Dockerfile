FROM node:14.6.0

RUN mkdir /etc/handshake-api
WORKDIR /etc/handshake-api/

COPY index.js /etc/handshake-api/index.js

RUN yarn init -y && \
    yarn add express express-http-proxy hs-client

ENV HOSTNAME="127.0.0.1"
ENV PORT=3100
ENV HSD_NETWORK="main"
ENV HSD_HOST="10.152.183.250"
ENV HSD_PORT=12037
#ENV HSD_API_KEY="foo"

EXPOSE $PORT

ENTRYPOINT ["node", "--max-http-header-size=81000", "/etc/handshake-api/index.js"]