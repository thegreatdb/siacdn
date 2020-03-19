#!/bin/bash

APIPASSWORD=`cat /root/.sia/apipassword`
BASE64_AUTHENTICATION=`echo -n ":$APIPASSWORD" | base64 -`
cat /etc/skynet/nginx.conf | sed -- 's#'"BASE64_AUTHENTICATION"'#'"$BASE64_AUTHENTICATION"'#g' - > /etc/nginx/conf.d/nginx.conf

nginx -g "daemon off;error_log /dev/stdout info;"