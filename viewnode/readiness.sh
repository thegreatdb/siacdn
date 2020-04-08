#!/usr/bin/env bash

set -e

cd /etc/sia
/go/bin/setup-env.sh
/go/bin/siac -d /etc/sia