#!/usr/bin/env sh

TMPDIR=/sia/tmp-`uuidgen`
mkdir $TMPDIR

if ! type "mc" > /dev/null; then
  echo "Installing and configuring the minio client..."
  curl -O https://dl.minio.io/client/mc/release/linux-amd64/mc
  chmod +x mc
  mv mc /bin/
  mc config host add minio $MINIO_URL $MINIO_ACCESS_KEY $MINIO_SECRET_KEY
fi

echo "1/4) Deleting any leftover tempfiles from a previous run..."
rm -f $TMPDIR/consensus.db \
  $TMPDIR/transactionpool.db \
  $TMPDIR/nodes.json \
  $TMPDIR/host.db \
  $TMPDIR/host.json \
  $TMPDIR/contractmanager.json \
  $TMPDIR/contractmanager.wal \
  $TMPDIR/hostdb.json \
  $TMPDIR/contractor.json

echo "2/4) Making copies of the data files..."
cp /sia/consensus/consensus.db \
  /sia/transactionpool/transactionpool.db \
  /sia/gateway/nodes.json \
  /sia/host/host.db \
  /sia/host/host.json \
  /sia/host/contractmanager/contractmanager.json \
  /sia/host/contractmanager/contractmanager.wal \
  /sia/renter/hostdb.json \
  /sia/renter/contractor.json $TMPDIR/

echo "3/4) Uploading data file copies..."
mc cp $TMPDIR/consensus.db \
  $TMPDIR/transactionpool.db \
  $TMPDIR/nodes.json \
  $TMPDIR/host.db \
  $TMPDIR/host.json \
  $TMPDIR/contractmanager.json \
  $TMPDIR/contractmanager.wal \
  $TMPDIR/hostdb.json \
  $TMPDIR/contractor.json minio/sia/

echo "4/4) Cleaning up..."
rm -f $TMPDIR/consensus.db \
  $TMPDIR/transactionpool.db \
  $TMPDIR/nodes.json \
  $TMPDIR/host.db \
  $TMPDIR/host.json \
  $TMPDIR/contractmanager.json \
  $TMPDIR/contractmanager.wal \
  $TMPDIR/hostdb.json \
  $TMPDIR/contractor.json

rm -Rf $TMPDIR

echo "Backup complete."
