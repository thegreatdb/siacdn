#!/usr/bin/env bash

set -e

cd /etc/sia

/go/bin/setup-env.sh

#echo "Bootstrapping wallet with password: $SIA_WALLET_PASSWORD"

if [ ! -f consensus/consensus.db ]; then
    echo "Found no consensus.db, downloading now..."
    curl https://siastats.info/bootstrap/bootstrap.zip -o /tmp/bootstrap.zip
    echo "Finished downloading consensus.db, extracting now..."
    unzip /tmp/bootstrap.zip -d /etc/sia
    echo "Finished bootstrapping consensus.db"
fi

if [ ! -f renter/hostdb.json ]; then
    echo "Found no hostdb.json, downloading now..."
    curl https://siastats.info/bootstrap/hostdb.zip -o /tmp/hostdb.zip
    echo "Finished downloading hostdb.json, extracting now..."
    unzip /tmp/hostdb.zip -d /etc/sia
    echo "Finished bootstrapping hostdb.json"
fi

/go/bin/siad -d /etc/sia "$@"