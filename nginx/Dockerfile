FROM nginx:1.17

# For cache invalidation
ENV SIACDN_BUILD_VERSION 6

RUN mkdir /etc/skynet/

COPY nginx.conf /etc/skynet/
COPY portals.json /etc/skynet/
COPY start-nginx.sh /etc/skynet/
RUN chmod +x /etc/skynet/start-nginx.sh