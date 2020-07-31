#!/bin/bash

set -e

ORDINAL_ID=`echo -n $HOSTNAME | cut -d "-" -f3`
API_PASSWORD_ENVNAME="SIA_API_PASSWORD_$ORDINAL_ID"
echo "API_PASSWORD_ENVNAME: $API_PASSWORD_ENVNAME"
export SIA_API_AUTHORIZATION=`printf '%s' "${!API_PASSWORD_ENVNAME}"`
#echo "SIA_API_AUTHORIZATION: $SIA_API_AUTHORIZATION"
BASE64_AUTHENTICATION=`echo -n ":$SIA_API_AUTHORIZATION" | base64 -`
cat /etc/nginx/conf.d/include/sia-auth.template | \
    sed -- 's#'"BASE64_AUTHENTICATION"'#'"$BASE64_AUTHENTICATION"'#g' - > \
    /etc/nginx/conf.d/include/sia-auth

/usr/bin/openresty -g "daemon off;error_log /dev/stdout info;"