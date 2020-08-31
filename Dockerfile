FROM golang:1.13
WORKDIR /go/src/github.com/k14s/kapp-controller/

RUN apt-get -y update && apt-get install -y ca-certificates && update-ca-certificates

# k14s
RUN wget -O- https://github.com/k14s/ytt/releases/download/v0.30.0/ytt-linux-amd64 > /usr/local/bin/ytt && \
  echo "456e58c70aef5cd4946d29ed106c2b2acbb4d0d5e99129e526ecb4a859a36145  /usr/local/bin/ytt" | shasum -c - && \
  chmod +x /usr/local/bin/ytt && ytt version

RUN wget -O- https://github.com/k14s/kapp/releases/download/v0.33.0/kapp-linux-amd64 > /usr/local/bin/kapp && \
  echo "2a3328c9eca9f43fe639afb524501d9d119feeea52c8a913639cfb96e38e93d1  /usr/local/bin/kapp" | shasum -c - && \
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

# kapp-controller
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags=-buildid= -trimpath -o controller ./cmd/controller/...

# ---
# Needs ubuntu for installing git/openssh
FROM ubuntu:bionic

RUN apt-get update && apt-get install -y git openssh-client && rm -rf /var/lib/apt/lists/*

RUN groupadd -g 2000 kapp-controller && useradd -r -u 1000 --create-home -g kapp-controller kapp-controller
USER kapp-controller

# Name it kapp-controller to identify it easier in process tree
COPY --from=0 /go/src/github.com/k14s/kapp-controller/controller kapp-controller

# fetchers
COPY --from=0 /helm-unpacked/linux-amd64/helm .
COPY --from=0 /usr/local/bin/imgpkg .

# templaters
COPY --from=0 /usr/local/bin/ytt .
COPY --from=0 /usr/local/bin/kbld .

# deployers
COPY --from=0 /usr/local/bin/kapp .

ENV PATH="/:${PATH}"
ENTRYPOINT ["/kapp-controller"]
