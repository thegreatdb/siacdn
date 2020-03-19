# SiaCDN

SiaCDN is a high-quality hosted Skynet node.

This repository is the complete set of scripts that is used to run this high-quality hosted Skynet node.

## How can I use this repository?

This repository can be used to easily deploy a Skynet node of your own on Kubernetes. Here's how:

1. Fork this repo on GitHub. It'll just take a second - I'll wait!
2. Navigate to your fork and continue reading this README there.
3. Clone your fork of the repo and make the following configuration changes, or do it using GitHub's built-in editor.


## Prerequisites

1. You should have Docker installed on your local machine.
2. You should have access to a Kubernetes cluster, and have kubectl authenticated to connect to it.
3. The Kubernetes cluster should use nginx ingress an installation of cert-manager with a cluster-issuer named `letsencrypt-prod`.
4. You should have a domain configured to point to your Kubernetes cluster.


## Customization and configuration changes before we can begin

1. Copy `kube/siacdn-sia-upload-secret.yaml.template` to `kube/siacdn-sia-upload-secret.yaml` and fill in SIA_WALLET_PASSWORD with the result of `echo -n "YOUR SEED PHRASE HERE" | base64 -w0`
2. Copy `kube/siacdn-sia-secret.yaml.template` to `kube/siacdn-sia-secret.yaml` and fill in SIA_WALLET_PASSWORD with the result of `echo -n "YOUR SEED PHRASE HERE" | base64 -w0`
3. Change `ericflo/siacdn-nginx:latest` to `YOUR_DOCKERHUB_NAME/siacdn-nginx:latest` in bin/*, repeat for `siacdn-portal` and `siacdn-viewnode`.
4. Do the same change for the Docker image name values in `kube/siacdn-deployment.yaml`.
5. Create `siacdn-nginx`, `siacdn-portal`, and `siacdn-viewnode` projects on your Docker Hub account.
6. Edit the `server_name` field in `nginx/nginx.conf` to match your domain name.
7. Change `SIACDN_DOMAIN` in `portal/Dockerfile` to match your domain.
8. Run `bin/docker-build-nginx`, `bin/docker-build-portal`, and `bin/docker-build-viewnode`, which will build and upload the docker images to your account.
9. Commit and push all these changes to your fork of this repo.


## Installing SiaCDN

1. Clone your fork of this repository and cd to the directory.
2. Create the kubernetes resources from the local kube dir.

```
cd siacdn
kubectl create -f kube/
```


## Configuring your Sia node

1. Run siac wallet init and wait for consensus, checking using plain siac command.

```
kubectl exec -it deployment/siacdn-deployment -c sia -- siac wallet init
kubectl exec -it deployment/siacdn-deployment -c sia -- siac
```

2. Run siac wallet address to get an address, and send that address some siacoins, probably 25K or so.

```
kubectl exec -it deployment/siacdn-deployment -c sia -- siac wallet address
```

3. Run siac renter setallowance twice, with parameters detailed by Nebulous at the following link, but use your judgement and adjust limits to your setup:

[https://github.com/NebulousLabs/skynet-webportal/tree/master/setup-scripts#portal-setup](https://github.com/NebulousLabs/skynet-webportal/tree/master/setup-scripts#portal-setup)

```
kubectl exec -it deployment/siacdn-deployment -c sia -- siac renter setallowance
kubectl exec -it deployment/siacdn-deployment -c sia -- siac renter setallowance --payment-contract-initial-funding 10SC
```

4. Repeat steps 1 to 3, but instead of `-c sia`, use `-c sia-upload`, and instead of `siac` use `siac -addr localhost:9970`. Here's an example:

```
kubectl exec -it deployment/siacdn-deployment -c sia-upload -- siac -addr localhost:9970 wallet init
kubectl exec -it deployment/siacdn-deployment -c sia-upload -- siac -addr localhost:9970 wallet address
kubectl exec -it deployment/siacdn-deployment -c sia-upload -- siac -addr localhost:9970 renter setallowance --payment-contract-initial-funding 10SC
```