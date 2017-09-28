#!/usr/bin/env sh

set -e

cd /sia

if [ ! -d consensus ]; then
  echo "No consensus file found, downloading latest snapshot..."

  mkdir consensus && cd consensus
  curl -O http://minio.maxint.co/sia/consensus.db -H "Accept-Encoding: gzip, deflate, sdch" --compressed
  cd ..

  mkdir transactionpool && cd transactionpool
  curl -O http://minio.maxint.co/sia/transactionpool.db -H "Accept-Encoding: gzip, deflate, sdch" --compressed
  cd ..

  mkdir gateway && cd gateway
  curl -O http://minio.maxint.co/sia/nodes.json -H "Accept-Encoding: gzip, deflate, sdch" --compressed
  cd ..

  mkdir host && cd host
  curl -O http://minio.maxint.co/sia/host.db -H "Accept-Encoding: gzip, deflate, sdch" --compressed
  curl -O http://minio.maxint.co/sia/host.json -H "Accept-Encoding: gzip, deflate, sdch" --compressed
  mkdir contractmanager && cd contractmanager
  curl -O http://minio.maxint.co/sia/contractmanager.wal -H "Accept-Encoding: gzip, deflate, sdch" --compressed
  curl -O http://minio.maxint.co/sia/contractmanager.json -H "Accept-Encoding: gzip, deflate, sdch" --compressed
  cd ../..

  mkdir renter && cd renter
  curl -O http://minio.maxint.co/sia/hostdb.json -H "Accept-Encoding: gzip, deflate, sdch" --compressed
  curl -O http://minio.maxint.co/sia/contractor.journal -H "Accept-Encoding: gzip, deflate, sdch" --compressed
  cd ..
fi

siad -d /sia --authenticate-api --disable-api-security --api-addr :9980
