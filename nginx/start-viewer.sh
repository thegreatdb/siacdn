#!/bin/bash

set -e

APIPASSWORD=`cat /etc/sia/apipassword`
BASE64_AUTHENTICATION=`echo -n ":$APIPASSWORD" | base64 -`
cat /etc/skynet/viewer.conf | \
    sed -- 's#'"BASE64_AUTHENTICATION"'#'"$BASE64_AUTHENTICATION"'#g' - | \
    sed -- 's#'"SKYNET_HOSTNAME"'#'"$SKYNET_HOSTNAME"'#g' - | \
    sed -- 's#'"SKYNET_HOSTNAME_ALT"'#'"$SKYNET_HOSTNAME_ALT"'#g' - | \
    /etc/nginx/conf.d/nginx.conf

nginx -g "daemon off;error_log /dev/stdout info;"