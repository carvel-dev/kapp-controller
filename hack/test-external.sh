#!/bin/bash

set -e -x -u

kapp deploy -y -a rbac -f examples/rbac/

time kapp deploy -y -a concourse-helm -f examples/concourse-helm.yml
time kapp delete -y -a concourse-helm

time kapp deploy -y -a consul-image-helm -f examples/consul-image-helm.yml
time kapp delete -y -a consul-image-helm

time kapp deploy -y -a istio-knative -f examples/istio-knative.yml
time kapp delete -y -a istio-knative

time kapp deploy -y -a nginx-helm-git -f examples/nginx-helm-git.yml
time kapp delete -y -a nginx-helm-git

time kapp deploy -y -a redis-helm -f examples/redis-helm.yml
time kapp delete -y -a redis-helm

time kapp deploy -y -a simple-app-http -f examples/simple-app-http.yml
time kapp delete -y -a simple-app-http

kapp delete -y -a rbac

echo EXTERNAL SUCCESS
