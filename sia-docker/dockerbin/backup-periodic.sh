#!/usr/bin/env sh

set -e

while :
do
	sleep 1h
	/root/backup.sh
	sleep 8h
done
