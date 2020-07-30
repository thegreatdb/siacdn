# SiaCDN

SiaCDN is a high-quality hosted Skynet node.

This repository is the complete set of scripts that is used to run this high-quality hosted Skynet node in your kubernetes cluster.

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

### NOTE: This section is very out of date. New documentation TBA

1. Make a copy of `kube/secrets.yaml.template` and save it as `kube/secrets.yaml`, then fill the values of `SIA_WALLET_PASSWORD_N` with the result of `echo -n "YOUR SEED PHRASE HERE" | base64 -w0` for as many viewers and uploaders as you want.
2. Choose API passwords for each of your nodes (e.g. random UUID) and fill in `SIA_API_PASSWORD_N` with those values. Use `echo -n "YOUR API PASSWORD HERE" | base64 -w0` to base64 encode it for the yaml file.
3. Edit `SKYNET_HOSTNAME` and `SKYNET_HOSTNAME_ALT` in `kube/config.yaml` to set it to your domain instead of the defaults.
4. Edit the number of replicas in `kube/uploader.yaml` and `kube/viewer.yaml` to match the number of secrets you filled in for step 2.
5. Edit the domains in `kube/ingress.yaml` to point to your domain instead of the defaults.
6. Change the `storageClassName` in both `kube/uploader.yaml` and `kube/viewer.yaml` to match the storage class in your kubernetes cluster. It may be that you want to simply delete the `storageClassName` lines, so that your default storage class will provision the volume claims.
7. Commit and push all these changes to your fork of this repo.


## Installing SiaCDN

1. Clone your fork of this repository and cd to the directory.
2. Create the kubernetes resources from the local kube dir.

```
cd siacdn
kubectl create -f kube/
```


## Configuring your Sia node

1. Run siac wallet init and wait for consensus for the uploaders, checking using siac.sh:

```
kubectl exec -it siacdn-uploader-0 -c sia -- siac.sh wallet init
kubectl exec -it siacdn-uploader-0 -c sia -- siac.sh
```

Do the same for the viewers:

```
kubectl exec -it siacdn-viewer-0 -c sia -- siac.sh wallet init
kubectl exec -it siacdn-viewer-0 -c sia -- siac.sh
```

2. Run siac wallet address to get an address for each node, and send that address some siacoins, probably 25K or so.

```
kubectl exec -it siacdn-uploader-0 -c sia -- siac.sh wallet address
```

3. Run siac renter setallowance twice per viewer and uploader, with parameters detailed by Nebulous at the following link, but use your judgement and adjust limits to your setup:

[https://github.com/NebulousLabs/skynet-webportal/tree/master/setup-scripts#portal-setup](https://github.com/NebulousLabs/skynet-webportal/tree/master/setup-scripts#portal-setup)

```
kubectl exec -it siacdn-uploader-0 -c sia -- siac.sh renter setallowance
kubectl exec -it siacdn-uploader-0 -c sia -- siac.sh renter setallowance --payment-contract-initial-funding 10SC
```