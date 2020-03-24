#!/bin/bash

READ_APIPASSWORD=`cat /root/.sia/apipassword`
WRITE_APIPASSWORD=`cat /root/.sia-upload/apipassword`
READ_BASE64_AUTHENTICATION=`echo -n ":$READ_APIPASSWORD" | base64 -`
WRITE_BASE64_AUTHENTICATION=`echo -n ":$WRITE_APIPASSWORD" | base64 -`
cat /etc/skynet/nginx.conf | \
    sed -- 's#'"READ_BASE64_AUTHENTICATION"'#'"$READ_BASE64_AUTHENTICATION"'#g' - | \
    sed -- 's#'"WRITE_BASE64_AUTHENTICATION"'#'"$WRITE_BASE64_AUTHENTICATION"'#g' - > /etc/nginx/conf.d/nginx.conf

nginx -g "daemon off;error_log /dev/stdout info;"