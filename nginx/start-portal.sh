#!/bin/bash

set -e

# Wait for apipassword
while [ ! -f /etc/sia/apipassword ]; do
    echo "Waiting for API Password"
    sleep 1
done

APIPASSWORD=`cat /etc/sia/apipassword`
BASE64_AUTHENTICATION=`echo -n ":$APIPASSWORD" | base64 -`
cat /etc/skynet/portal.conf | \
    sed -- 's#'"SKYNET_HOSTNAME"'#'"$SKYNET_HOSTNAME"'#g' - | \
    sed -- 's#'"SKYNET_HOSTNAME_ALT"'#'"$SKYNET_HOSTNAME_ALT"'#g' - | \
    /etc/nginx/conf.d/nginx.conf

nginx -g "daemon off;error_log /dev/stdout info;"