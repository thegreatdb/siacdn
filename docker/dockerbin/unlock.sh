#!/usr/bin/env sh

until siac wallet unlock > /dev/null 2>&1; do
  if [ ! -z "$SIA_WALLET_PASSWORD" ]; then
    echo "Could not unlock wallet, waiting 5 secs...";
  fi
  sleep 5;
done
