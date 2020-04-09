#!/usr/bin/env bash

set -e

cd /etc/sia
source /go/bin/setup-env.sh
/go/bin/siac -d /etc/sia wallet | grep Unlocked