#!/bin/bash

set -e

grep -H www.siacdn.com -R /skynet-webportal/ -lz | \
    xargs -t -l sed -i -e 's#'"www.siacdn.com"'#'"$SKYNET_HOSTNAME"'#g'

cat /etc/skynet/portal.conf | \
    sed -- 's#'"SKYNET_HOSTNAME"'#'"$SKYNET_HOSTNAME"'#g' - | \
    sed -- 's#'"SKYNET_HOSTNAME_ALT"'#'"$SKYNET_HOSTNAME_ALT"'#g' - > \
    /etc/nginx/conf.d/nginx.conf

nginx -g "daemon off;error_log /dev/stdout info;"