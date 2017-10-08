# siacdn-backend

This is the backend server for SiaCDN.

## Commands

Currently this produces a binary which has two commands: serve, and kube.

* **serve**: Starts the API server, which starts an HTTP server on port 9095.
Check out [makeRouter in server.go](./server/server.go) for a list of routes.
* **kube**: Coordinates between the API server and a Kubernetes cluster to
perform the actual deployment and provisioning of Sia and Minio instances.
