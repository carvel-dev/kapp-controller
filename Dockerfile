FROM photon:4.0

ARG KCTRL_VER=development

# adapted from golang docker image
ENV PATH /usr/local/go/bin:$PATH
ENV GOLANG_VERSION 1.19.5
ENV GO_REL_ARCH linux-amd64
ENV GO_REL_SHA 7629a36ea5ee00c30df8aef33a954012ab3884265af95dda08ada393f435f340

RUN set eux; \
    curl -sLo go.tgz "https://golang.org/dl/go${GOLANG_VERSION}.${GO_REL_ARCH}.tar.gz"; \
    echo "${GO_REL_SHA} go.tgz" | sha256sum -c -; \
    tar -C /usr/local -xzf go.tgz; \
    rm go.tgz; \
    go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR /go/src/github.com/vmware-tanzu/carvel-kapp-controller/

# carvel
COPY ./hack/install-deps.sh .
COPY ./hack/dependencies.yml .
RUN ./install-deps.sh

# [DEPRECATED] Helm V2
# Maintaining two versions of helm until we drop support in a future release
RUN curl -sLo /helm https://get.helm.sh/helm-v2.17.0-linux-amd64.tar.gz && \
  echo "f3bec3c7c55f6a9eb9e6586b8c503f370af92fe987fcbf741f37707606d70296  /helm" | sha256sum -c - && \
  mkdir /helm-v2-unpacked && tar -C /helm-v2-unpacked -xzvf /helm

RUN curl -sLo /helm https://get.helm.sh/helm-v3.10.3-linux-amd64.tar.gz && \
  echo "950439759ece902157cf915b209b8d694e6f675eaab5099fb7894f30eeaee9a2  /helm" | sha256sum -c - && \
  mkdir /helm-unpacked && tar -C /helm-unpacked -xzvf /helm

# sops
RUN curl -sLo /usr/local/bin/sops https://github.com/mozilla/sops/releases/download/v3.7.3/sops-v3.7.3.linux && \
  echo "53aec65e45f62a769ff24b7e5384f0c82d62668dd96ed56685f649da114b4dbb  /usr/local/bin/sops" | sha256sum -c - && \
  chmod +x /usr/local/bin/sops && sops -v

# age (encryption for sops)
RUN curl -sLo age.tgz https://github.com/FiloSottile/age/releases/download/v1.0.0/age-v1.0.0-linux-amd64.tar.gz && \
  echo "6414f71ce947fbbea1314f6e9786c5d48436ebc76c3fd6167bf018e432b3b669  age.tgz" | sha256sum -c - && \
  tar -xzf age.tgz && cp age/age /usr/local/bin && \
  chmod +x /usr/local/bin/age && age --version

RUN curl -sLo cue.tgz https://github.com/cue-lang/cue/releases/download/v0.4.3/cue_v0.4.3_linux_amd64.tar.gz && \
  echo "5e7ecb614b5926acfc36eb1258800391ab7c6e6e026fa7cacbfe92006bac895c cue.tgz" | sha256sum -c - && \
  tar -xf cue.tgz -C /usr/local/bin cue && cue version

# kapp-controller
COPY . .
# helpful ldflags reference: https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-X 'main.Version=$KCTRL_VER' -buildid=" -trimpath -o controller ./cmd/main.go

# --- run image ---
FROM photon:4.0

# Install openssh for git
# TODO(bmo): why do we need sed?
RUN tdnf install -y git openssh-clients sed

# Create the kapp-controller user in the root group, the home directory will be mounted as a volume
RUN echo "kapp-controller:x:1000:0:/home/kapp-controller:/usr/sbin/nologin" > /etc/passwd
# Give the root group write access to openssh's root bundle so we can append custom roots at runtime
RUN chmod g+w /etc/pki/tls/certs/ca-bundle.crt

# fetchers
COPY --from=0 /helm-v2-unpacked/linux-amd64/helm helmv2
COPY --from=0 /helm-unpacked/linux-amd64/helm .
COPY --from=0 /usr/local/bin/imgpkg .
COPY --from=0 /usr/local/bin/vendir .

# templaters
COPY --from=0 /usr/local/bin/ytt .
COPY --from=0 /usr/local/bin/kbld .
COPY --from=0 /usr/local/bin/sops .
COPY --from=0 /usr/local/bin/age .
COPY --from=0 /usr/local/bin/cue .

# deployers
COPY --from=0 /usr/local/bin/kapp .

# Name it kapp-controller to identify it easier in process tree
COPY --from=0 /go/src/github.com/vmware-tanzu/carvel-kapp-controller/controller kapp-controller

# Run as kapp-controller by default, will be overridden to a random uid on OpenShift
USER 1000
ENV PATH="/:${PATH}"
ENTRYPOINT ["/kapp-controller"]
