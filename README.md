# SiaCDN

SiaCDN is a high-quality hosted Skynet node.

This repository is the complete set of scripts that is used to run this high-quality hosted Skynet node.

## How can I use this repository?

This repository can be used to easily deploy a Skynet node of your own on Google Kubernetes Engine. Here's how:

1. Fork this repo on GitHub. It'll just take a second - I'll wait!
2. Navigate to your fork and continue reading this README there.
3. Clone your fork of the repo and make the following configuration changes, or do it using GitHub's built-in editor.


## Configuration changes before we can begin

1. Change the value in kube/siacdn-certificate-toplevel.yaml to match your domain name (line 7)
2. Change the value in kube/siacdn-certificate-www.yaml to match your domain name, keeping the www (line 7)
3. Ensure that these changes are pushed to your clone of this repo.


## Prerequisites

1. You should have Docker installed on your local machine.
2. You should have a Google Kubernetes Engine cluster, running, and have kubectl authenticated to connect to it. (gcloud container clusters create siacdn && gcloud container clusters get-credentials siacdn)
3. You should have a reserved external IP address named siacdn-ip-address (gcloud compute addresses create siacdn-ip-address --global)
4. Create two A records in your DNS that point to this reserved external IP address - one with a www prefix and another without.


## Installing SiaCDN

1. Clone your fork of this repository and cd to the directory.
2. Create the kubernetes resources from the local kube dir.

> cd siacdn

> kubectl create -f kube/


## Configuring your Sia node

1. Use kubectl to determine the pod where Sia is running:

> kubectl get pods -l app=siacdn -o name

2. Use the output of that to run siac commands on the running sia node:

> NOTE: In the following, siacdn-deployment-567cda4ea9-wx2sr is the output of #1, without pod/

> kubectl exec -it siacdn-deployment-567cda4ea9-wx2sr -c viewnode -- siac

3. Run siac wallet init and wait for consensus, checking using plain siac command.

> kubectl exec -it siacdn-deployment-567cda4ea9-wx2sr -c viewnode -- siac wallet init

> kubectl exec -it siacdn-deployment-567cda4ea9-wx2sr -c viewnode -- siac

4. Run siac wallet address to get an address, and send that address some siacoins, probably 25K or so.

> kubectl exec -it siacdn-deployment-567cda4ea9-wx2sr -c viewnode -- siac wallet address

5. Run siac renter setallowance twice, with parameters detailed by Nebulous at the following link, but use your judgement and adjust limits to your setup:

[https://github.com/NebulousLabs/skynet-webportal/tree/master/setup-scripts#viewnode-setup](https://github.com/NebulousLabs/skynet-webportal/tree/master/setup-scripts#viewnode-setup)

> kubectl exec -it siacdn-deployment-567cda4ea9-wx2sr -c viewnode -- siac renter setallowance

> kubectl exec -it siacdn-deployment-567cda4ea9-wx2sr -c viewnode -- siac renter setallowance --payment-contract-initial-price 10SC
