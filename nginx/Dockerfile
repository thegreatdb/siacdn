FROM nginx:1.17

RUN mkdir /etc/skynet/

COPY nginx.conf /etc/skynet/
COPY portals.json /etc/skynet/
COPY start-nginx.sh /etc/skynet/
RUN chmod +x /etc/skynet/start-nginx.sh