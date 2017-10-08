# SiaCDN

This repository contains a set of projects that can be combined to run a service
that:

* Provisions new full Sia nodes
* Provisions Sia-powered Minio instances
* Receives payments in U.S. Dollars ($)
* Provides user-facing administration via a hosted website

## Projects

* **/frontend (React/Next.js)**: A fairly standard React and next.js-based web
frontend that speaks to the SiaCDN HTTP API and helps users with management
tasks.
* **/backend (Go)**: Powers the SiaCDN HTTP API, keeps track of metadata like
customer login credentials, and controls the deployment of resources to
Kubernetes.
* **/docker (Dockerfile)**: Dockerfiles for projects that have been customized
to work in the SiaCDN environment, such as our patched/updated versions of Sia
and Minio.
* **/kube (YAML)**: This is the suite of YAML files that you will want to run
SiaCDN on your Kubernetes cluster.
* **/bin (Bash)**: Scripts that perform actions in the repo, like building and
uploading Docker artifacts, tunneling ports for development, etc.

## License

SiaCDN is [MIT licensed](./LICENSE).
