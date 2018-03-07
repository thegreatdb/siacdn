# SiaCDN

SiaCDN is a protocol for Sia hosts to pool their resources to create a
for-profit CDN service to the internet. It's like a protocol for a mining pool,
but for Sia hosts instead of miners.

## How it works

There are three main components that you need to run a CDN on Sia:

* Edge Nodes: HTTP servers that serve files to users in exchange for SC.
* DNS Clusters: DNS routers that can fairly route requests to the nearest
  trusted edge node that has the file.
* Validating Clients: Clients that periodically, randomly request parts of
  files from edge nodes and submit proofs of them, in exchange for a period of
  lowered fees for fetching their other files.

Running an edge node is extremely easy, you simply run a SiaCDN daemon
alongside Sia:

```
$ siacdn edgenode -join=cdn.example.com=5000,cdn.example2.com=1000
```

This will start an edge node and put up 5000 SC in collateral to cdn.example.com
and 1000 SC in collateral to cdn.example2.com, which are SiaCDN-compatible DNS
nodes. The more collateral your edge node puts up, the easier time your node
will have of earning trust, thus receiving more requests and earning more
income.

Running a DNS cluster node is similarly easy:

```
$ siacdn dnsnode -name=cdn.example.com -pass=116d8d1021dd11e8af15636c85f905ed
```

You have to have ICANN control over the domain name and have the domain pointed
at the IP addresses of your DNS cluster nodes. Make sure you run each node with
the same password.

## Isn't that centralized?

SiaCDN edge nodes are totally decentralized. Anyone who runs a Sia node and has
some SC for collateral can take part, earning income based on file serving
performance and collateralized trust.

However, there's no getting around the centralization of DNS for a CDN right
now. So the idea is to build this on top of the most decentralized base storage
service out there, Sia, and then make it a protocol so that if one DNS cluster
starts acting badly, it's easier for the community to punish them by switching
to another compatible one.

This incentive structure, designed from the start, should create pressure to
act fairly even though there are some counteracting centralization pressures.
