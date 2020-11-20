FROM golang:1.13
WORKDIR /go/src/github.com/vmware-tanzu/carvel-kapp-controller/

RUN apt-get -y update && apt-get install -y ca-certificates && update-ca-certificates

# k14s
RUN wget -O- https://github.com/k14s/ytt/releases/download/v0.30.0/ytt-linux-amd64 > /usr/local/bin/ytt && \
  echo "456e58c70aef5cd4946d29ed106c2b2acbb4d0d5e99129e526ecb4a859a36145  /usr/local/bin/ytt" | shasum -c - && \
  chmod +x /usr/local/bin/ytt && ytt version

RUN wget -O- https://github.com/k14s/kapp/releases/download/v0.34.0/kapp-linux-amd64 > /usr/local/bin/kapp && \
  echo "e170193c40ff5dff9f9274c25048de1f50e23c69e8406df274fbb416d5862d7f  /usr/local/bin/kapp" | shasum -c - && \
  chmod +x /usr/local/bin/kapp && kapp version

RUN wget -O- https://github.com/k14s/kbld/releases/download/v0.24.0/kbld-linux-amd64 > /usr/local/bin/kbld && \
  echo "63f06c428cacd66e4ebbd23df3f04214109bc44ee623c7c81ecb9aa35c192c65  /usr/local/bin/kbld" | shasum -c - && \
  chmod +x /usr/local/bin/kbld && kbld version

RUN wget -O- https://github.com/k14s/imgpkg/releases/download/v0.2.0/imgpkg-linux-amd64 > /usr/local/bin/imgpkg && \
  echo "57a73c4721c39f815408f486c1acfb720af82450996e2bfdf4c2c280d8a28dcc  /usr/local/bin/imgpkg" | shasum -c - && \
  chmod +x /usr/local/bin/imgpkg && imgpkg version

# helm
RUN wget -O- https://get.helm.sh/helm-v2.14.3-linux-amd64.tar.gz > /helm && \
  echo "38614a665859c0f01c9c1d84fa9a5027364f936814d1e47839b05327e400bf55  /helm" | shasum -c - && \
  mkdir /helm-unpacked && tar -C /helm-unpacked -xzvf /helm

# sops
RUN wget -O- https://github.com/mozilla/sops/releases/download/v3.6.1/sops-v3.6.1.linux > /usr/local/bin/sops && \
  echo "b2252aa00836c72534471e1099fa22fab2133329b62d7826b5ac49511fcc8997  /usr/local/bin/sops" | shasum -c - && \
  chmod +x /usr/local/bin/sops && sops -v

# kapp-controller
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags=-buildid= -trimpath -o controller ./cmd/controller/...

# ---
# Needs ubuntu for installing git/openssh
FROM ubuntu:bionic

RUN apt-get update && apt-get install -y git openssh-client dumb-init && rm -rf /var/lib/apt/lists/*

RUN groupadd -g 2000 kapp-controller && useradd -r -u 1000 --create-home -g kapp-controller kapp-controller
USER kapp-controller

# Name it kapp-controller to identify it easier in process tree
COPY --from=0 /go/src/github.com/vmware-tanzu/carvel-kapp-controller/controller kapp-controller

# fetchers
COPY --from=0 /helm-unpacked/linux-amd64/helm .
COPY --from=0 /usr/local/bin/imgpkg .

# templaters
COPY --from=0 /usr/local/bin/ytt .
COPY --from=0 /usr/local/bin/kbld .
COPY --from=0 /usr/local/bin/sops .

# deployers
COPY --from=0 /usr/local/bin/kapp .

ENV PATH="/:${PATH}"
ENTRYPOINT ["dumb-init", "--", "/kapp-controller"]
