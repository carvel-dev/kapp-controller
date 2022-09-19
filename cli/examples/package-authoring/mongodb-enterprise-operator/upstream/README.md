# MongoDB Enterprise Kubernetes Operator Helm Chart

A Helm Chart for installing and upgrading the [MongoDB Enterprise
Kubernetes Operator](https://github.com/mongodb/mongodb-enterprise-kubernetes).

## Prerequisites

The installation of this Chart does not have prerequisites. However, in order to
create a Mongo Database in your Kubernetes cluster, you'll need a [Cloud
Manager](https://cloud.mongodb.com) account or an [Ops
Manager](https://www.mongodb.com/products/ops-manager) installation.

## Installing Enterprise Operator

You can install the MongoDB Enterprise Operator easily with:

``` shell
helm install enterprise-operator mongodb/enterprise-operator
```

This will install `CRD`s and the Enterprise Operator in the current namespace
(`default` by _default_). You can pass a different namespace with:

``` shell
helm install enterprise-operator mongodb/enterprise-operator --namespace mongodb [--create-namespace]
```

To install the Enterprise Operator in a namespace called `mongodb`; with the
optional `--create-namespace` Helm will create the Namespace if it does not exist.

## Configuring access to Cloud Manager

The Enteprise Operator can run against an account of [Cloud
Manager](https://cloud.mongodb.com) or [Ops
Manager](https://www.mongodb.com/products/ops-manager).

### Using Cloud Manager

Please visit "[Configure Kubernetes for Deploying MongoDB
Resource](https://docs.cloudmanager.mongodb.com/tutorial/nav/k8s-config-for-mdb-resource/)"
in order to create new `Secret`s and `ConfigMap`s with credentials for your
Cloud Manager project.

You will create a `ConfigMap` with name `my-project`, and a `Secret` with name
`my-credentials`.

## Deploying a MongoDB Replica Set

The Enterprise Operator will be watching for resources of type
`mongodb.mongodb.com` (among others); you can quickly install
a sample Mongo Database with:

``` shell
kubectl apply -f https://raw.githubusercontent.com/mongodb/mongodb-enterprise-kubernetes/master/samples/mongodb/minimal/replica-set.yaml [--namespace mongodb]
```

- _Note: Make sure you add the `--namespace` option when needed._

After a few minutes you will have a 3-member MongoDB Replica Set installed in
your cluster, that you can check with:

``` shell
$ kubectl get mdb
NAME              PHASE     VERSION
my-replica-set    Running   4.4.0-ent
```

## What to do next

Please follow the [Official MongoDB Enterprise Kubernetes
Operator](https://docs.mongodb.com/kubernetes-operator/stable/) for additional
deployment topologies and multitude of other MongoDB options.
