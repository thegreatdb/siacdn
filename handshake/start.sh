#!/bin/bash

set -e

/etc/handshake/hsd/bin/node \
    --prefix /var/lib/handshake \
    --rs-host 0.0.0.0 \
    --rs-port 53 \
    --http-host 0.0.0.0 \
    --http-port 12037 \
    --api-key="${HSD_API_KEY}"