#!/bin/bash

set -e -x -u

function dump {
  kubectl get apps -o yaml
}
trap dump ERR

kapp deploy -y -a rbac -f examples/rbac/

time kapp deploy -y -a simple-app -f examples/simple-app-git/1.yml
time kapp deploy -y -a simple-app -f examples/simple-app-git/2.yml
time kapp delete -y -a simple-app

## requires multiple node cluster
# time kapp deploy -y -a consul-image-helm -f examples/consul-image-helm.yml
# time kapp delete -y -a consul-image-helm

## requires multiple node cluster
# time kapp deploy -y -a istio-knative -f examples/istio-knative.yml
# time kapp delete -y -a istio-knative

time kapp deploy -y -a nginx-helm-git -f examples/nginx-helm-git.yml
time kapp delete -y -a nginx-helm-git

time kapp deploy -y -a redis-helm -f examples/redis-helm.yml
time kapp delete -y -a redis-helm

time kapp deploy -y -a simple-app-http -f examples/simple-app-http.yml
time kapp delete -y -a simple-app-http

time kapp deploy -y -a cue -f examples/cue.yml
time kapp delete -y -a cue

time kapp deploy -y -a step-values-and-config -f examples/pkgi-with-config-and-values-per-step.yaml
time kapp delete -y -a step-values-and-config

kapp delete -y -a rbac

echo EXTERNAL SUCCESS
