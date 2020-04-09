#!/bin/bash

set -e

ORDINAL_ID=`echo -n $HOSTNAME | rev | cut -d "-" -f1 | rev`
API_PASSWORD_ENVNAME="SIA_API_PASSWORD_$ORDINAL_ID"
echo "API_PASSWORD_ENVNAME: $API_PASSWORD_ENVNAME"
export SIA_API_PASSWORD=`printf '%s' "${!API_PASSWORD_ENVNAME}"`
echo "SIA_API_PASSWORD: $SIA_API_PASSWORD"
BASE64_AUTHENTICATION=`echo -n ":$SIA_API_PASSWORD" | base64 -`
cat /etc/skynet/uploader.conf | \
    sed -- 's#'"BASE64_AUTHENTICATION"'#'"$BASE64_AUTHENTICATION"'#g' - | \
    sed -- 's#'"SKYNET_HOSTNAME"'#'"$SKYNET_HOSTNAME"'#g' - | \
    sed -- 's#'"SKYNET_HOSTNAME_ALT"'#'"$SKYNET_HOSTNAME_ALT"'#g' - > \
    /etc/nginx/conf.d/nginx.conf

nginx -g "daemon off;error_log /dev/stdout info;"