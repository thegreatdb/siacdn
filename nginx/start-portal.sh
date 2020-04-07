#!/bin/bash

set -e

cat /etc/skynet/portal.conf | \
    sed -- 's#'"SKYNET_HOSTNAME"'#'"$SKYNET_HOSTNAME"'#g' - | \
    sed -- 's#'"SKYNET_HOSTNAME_ALT"'#'"$SKYNET_HOSTNAME_ALT"'#g' - > \
    /etc/nginx/conf.d/nginx.conf

nginx -g "daemon off;error_log /dev/stdout info;"