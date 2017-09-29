#!/usr/bin/env sh

if ! type "mc" > /dev/null; then
  echo "Installing and configuring the minio client..."
  curl -O https://dl.minio.io/client/mc/release/linux-amd64/mc
  chmod +x mc
  mv mc /bin/
  mc config host add minio $MINIO_URL $MINIO_ACCESS_KEY $MINIO_SECRET_KEY
fi

echo "1/4) Deleting any leftover tempfiles from a previous run..."
rm -f /tmp/consensus.db \
  /tmp/transactionpool.db \
  /tmp/nodes.json \
  /tmp/host.db \
  /tmp/host.json \
  /tmp/contractmanager.json \
  /tmp/contractmanager.wal \
  /tmp/hostdb.json \
  /tmp/contractor.journal

echo "2/4) Making copies of the data files..."
cp /sia/consensus/consensus.db \
  /sia/transactionpool/transactionpool.db \
  /sia/gateway/nodes.json \
  /sia/host/host.db \
  /sia/host/host.json \
  /sia/host/contractmanager/contractmanager.json \
  /sia/host/contractmanager/contractmanager.wal \
  /sia/renter/hostdb.json \
  /sia/renter/contractor.journal /tmp/

echo "3/4) Uploading data file copies..."
mc cp /tmp/consensus.db \
  /tmp/transactionpool.db \
  /tmp/nodes.json \
  /tmp/host.db \
  /tmp/host.json \
  /tmp/contractmanager.json \
  /tmp/contractmanager.wal \
  /tmp/hostdb.json \
  /tmp/contractor.journal minio/sia/

echo "4/4) Cleaning up..."
rm -f /tmp/consensus.db \
  /tmp/transactionpool.db \
  /tmp/nodes.json \
  /tmp/host.db \
  /tmp/host.json \
  /tmp/contractmanager.json \
  /tmp/contractmanager.wal \
  /tmp/hostdb.json \
  /tmp/contractor.journal

echo "Backup complete."
