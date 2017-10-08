# siacdn-kube

A Suite of Kubernetes YAML files that are variously useful for running SiaCDN.

## Subprojects

* **/backend**: Used to deploy an instance of the backend server.
* **/external-dns**: Used to register the new minio nodes with DNS so that you
can access them externally.
* **/ingress**: We depend on this nginx ingress controller, so you may want to
install it.
* **/prime**: We deploy one prime SiaNode instance that acts as the bank for
now. These scripts deploy it.
* **/src**: This runs a node that just keeps consensus and backs up its state
periodically.
* **/worker**: Used to deploy an instance of the kube worker daemon, to
communicate with the backend and Kubernetes.
* **/www**: Used to deploy the frontend.
