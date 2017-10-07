export const displayStatus = {
  created: '2) Sending specifications to the deployment server...',
  deployed: '3) Waiting for resources from deployment server...',
  instanced:
    '4) Initialized Sia node, now downloading a recent blockchain snapshot...',
  snapshotted:
    '5) Finished snapshotting, now downloading the latest blockchain updates...',
  synchronized: '6) Blockchain fully synced. Initializing wallet...',
  initialized:
    '7) Unlocking your wallet for the first time. This has to scan the blockchain, so it can take up to 30 minutes...',
  unlocked:
    '8) Transferring enough funds to your Sia node to meet requested storage capacity...',
  funded: '9) Funds sent. Waiting for confirmation that funds were received...',
  confirmed: '10) Funding confirmed. Setting storage contract allowance now...',
  configured:
    '11) Negotiating storage contracts on your behalf. This can take a while (Sia team working to improve this)...',
  ready: 'Sia node is up and running!',
  stopping: 'Stopping and tearing down.',
  stopped: 'Stopped.',
  depleted: 'Insufficient funds to continue.',
  error: 'Errored.',
};
