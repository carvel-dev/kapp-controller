FROM photon:3.0

RUN tdnf install -y tar wget gzip

# adapted from golang docker image
ENV PATH /usr/local/go/bin:$PATH
ENV GOLANG_VERSION 1.13.13
ENV GO_REL_ARCH linux-amd64
ENV GO_REL_SHA 0b8573c2335bebef53e819ab8d323456dc2b94838bebdbd8cc6623bb8a6d77b7

RUN set eux; \
    wget -O go.tgz "https://golang.org/dl/go${GOLANG_VERSION}.${GO_REL_ARCH}.tar.gz" --progress=dot:giga; \
    echo "${GO_REL_SHA} *go.tgz" | sha256sum -c -; \
    tar -C /usr/local -xzf go.tgz; \
    rm go.tgz; \
    go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR /go/src/github.com/vmware-tanzu/carvel-kapp-controller/

# carvel
RUN wget -O- https://github.com/k14s/ytt/releases/download/v0.30.0/ytt-linux-amd64 > /usr/local/bin/ytt && \
  echo "456e58c70aef5cd4946d29ed106c2b2acbb4d0d5e99129e526ecb4a859a36145  /usr/local/bin/ytt" | sha256sum -c - && \
  chmod +x /usr/local/bin/ytt && ytt version

RUN wget -O- https://github.com/k14s/kapp/releases/download/v0.34.0/kapp-linux-amd64 > /usr/local/bin/kapp && \
  echo "e170193c40ff5dff9f9274c25048de1f50e23c69e8406df274fbb416d5862d7f  /usr/local/bin/kapp" | sha256sum -c - && \
  chmod +x /usr/local/bin/kapp && kapp version

RUN wget -O- https://github.com/k14s/kbld/releases/download/v0.24.0/kbld-linux-amd64 > /usr/local/bin/kbld && \
  echo "63f06c428cacd66e4ebbd23df3f04214109bc44ee623c7c81ecb9aa35c192c65  /usr/local/bin/kbld" | sha256sum -c - && \
  chmod +x /usr/local/bin/kbld && kbld version

RUN wget -O- https://github.com/k14s/imgpkg/releases/download/v0.2.0/imgpkg-linux-amd64 > /usr/local/bin/imgpkg && \
  echo "57a73c4721c39f815408f486c1acfb720af82450996e2bfdf4c2c280d8a28dcc  /usr/local/bin/imgpkg" | sha256sum -c - && \
  chmod +x /usr/local/bin/imgpkg && imgpkg version

RUN wget -O- https://github.com/vmware-tanzu/carvel-vendir/releases/download/v0.14.0/vendir-linux-amd64 > /usr/local/bin/vendir && \
  echo "c224bdfe74df326d7e75b4c50669ec5976b95c0ff9a27d25c6e1833d0c781946  /usr/local/bin/vendir" | sha256sum -c - && \
  chmod +x /usr/local/bin/vendir && vendir version

# helm
RUN wget -O- https://get.helm.sh/helm-v2.17.0-linux-amd64.tar.gz > /helm && \
  echo "f3bec3c7c55f6a9eb9e6586b8c503f370af92fe987fcbf741f37707606d70296  /helm" | sha256sum -c - && \
  mkdir /helm-unpacked && tar -C /helm-unpacked -xzvf /helm

# sops
RUN wget -O- https://github.com/mozilla/sops/releases/download/v3.6.1/sops-v3.6.1.linux > /usr/local/bin/sops && \
  echo "b2252aa00836c72534471e1099fa22fab2133329b62d7826b5ac49511fcc8997  /usr/local/bin/sops" | sha256sum -c - && \
  chmod +x /usr/local/bin/sops && sops -v

# kapp-controller
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags=-buildid= -trimpath -o controller ./cmd/main.go

# --- run image ---
FROM photon:3.0

RUN tdnf install -y git openssh-clients shadow-tools sed

RUN groupadd -g 2000 kapp-controller && useradd -r -u 1000 --create-home -g kapp-controller kapp-controller
RUN chmod g+w /etc/pki/tls/certs/ca-bundle.crt && chgrp kapp-controller /etc/pki/tls/certs/ca-bundle.crt
USER kapp-controller

# Name it kapp-controller to identify it easier in process tree
COPY --from=0 /go/src/github.com/vmware-tanzu/carvel-kapp-controller/controller kapp-controller

# fetchers
COPY --from=0 /helm-unpacked/linux-amd64/helm .
COPY --from=0 /usr/local/bin/imgpkg .
COPY --from=0 /usr/local/bin/vendir .

# templaters
COPY --from=0 /usr/local/bin/ytt .
COPY --from=0 /usr/local/bin/kbld .
COPY --from=0 /usr/local/bin/sops .

# deployers
COPY --from=0 /usr/local/bin/kapp .

ENV PATH="/:${PATH}"
ENTRYPOINT ["/kapp-controller"]
