#!/usr/bin/env bash

set -e

cd /etc/sia

ORDINAL_ID=`$HOSTNAME | rev | cut -d "-" -f1 | rev`
CURRENT_PASSWORD_ENVNAME="SIA_WALLET_PASSWORD_$ORDINAL_ID"
export SIA_WALLET_PASSWORD=`printf '%s' "${!CURRENT_PASSWORD_ENVNAME}"`
echo "Bootstrapping wallet with password: $SIA_WALLET_PASSWORD"

if [ ! -f consensus/consensus.db ]; then
    echo "Found no consensus.db, downloading now..."
    curl -O https://siastats.info/bootstrap/bootstrap.zip
    echo "Finished downloading consensus.db, extracting now..."
    unzip bootstrap.zip
    echo "Finished bootstrapping consensus.db"
fi

if [ ! -f host/host.db ]; then
    echo "Found no host.db, downloading now..."
    curl -O https://siastats.info/bootstrap/hostdb.zip
    echo "Finished downloading host.db, extracting now..."
    unzip hostdb.zip
    echo "Finished bootstrapping host.db"
fi

/go/bin/siad -d /etc/sia "$@"