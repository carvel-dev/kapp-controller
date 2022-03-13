FROM photon:4.0

ARG KCTRL_VER=development
ARG TARGETARCH
ARG TARGETOS

# adapted from golang docker image
ENV PATH /usr/local/go/bin:$PATH
ENV GOLANG_VERSION 1.17.6

RUN set eux; \
    curl -sLo go.tgz "https://golang.org/dl/go${GOLANG_VERSION}.${TARGETOS}-${TARGETARCH}.tar.gz"; \
    if [ $TARGETARCH == "arm64" ] ; then export DOWNLOAD_SHA="82c1a033cce9bc1b47073fd6285233133040f0378439f3c4659fe77cc534622a" ; fi; \
    if [ $TARGETARCH == "amd64" ] ; then export DOWNLOAD_SHA="231654bbf2dab3d86c1619ce799e77b03d96f9b50770297c8f4dff8836fc8ca2" ; fi; \
    echo "${DOWNLOAD_SHA} go.tgz" | sha256sum -c -; \
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
#echo "f3bec3c7c55f6a9eb9e6586b8c503f370af92fe987fcbf741f37707606d70296  /helm" | sha256sum -c - && \
RUN curl -sLo /helm https://get.helm.sh/helm-v2.17.0-linux-${TARGETARCH}.tar.gz && \
  if [ $TARGETARCH == "arm64" ] ; then export DOWNLOAD_SHA="c3ebe8fa04b4e235eb7a9ab030a98d3002f93ecb842f0a8741f98383a9493d7f" ; fi && \
  if [ $TARGETARCH == "amd64" ] ; then export DOWNLOAD_SHA="f3bec3c7c55f6a9eb9e6586b8c503f370af92fe987fcbf741f37707606d70296" ; fi && \
  echo "${DOWNLOAD_SHA}  /helm" | sha256sum -c - && \
  mkdir /helm-v2-unpacked && tar -C /helm-v2-unpacked -xzvf /helm

RUN curl -sLo /helm https://get.helm.sh/helm-v3.7.1-linux-${TARGETARCH}.tar.gz && \
  if [ $TARGETARCH == "arm64" ] ; then export DOWNLOAD_SHA="57875be56f981d11957205986a57c07432e54d0b282624d68b1aeac16be70704" ; fi && \
  if [ $TARGETARCH == "amd64" ] ; then export DOWNLOAD_SHA="6cd6cad4b97e10c33c978ff3ac97bb42b68f79766f1d2284cfd62ec04cd177f4" ; fi && \
  echo "${DOWNLOAD_SHA}  /helm" | sha256sum -c - && \
  mkdir /helm-unpacked && tar -C /helm-unpacked -xzvf /helm

# sops
RUN curl -sLo /usr/local/bin/sops https://github.com/mozilla/sops/releases/download/v3.7.2/sops-v3.7.2.linux.${TARGETARCH} && \
  if [ $TARGETARCH == "arm64" ] ; then export DOWNLOAD_SHA="86a6c48ec64255bd317d7cd52c601dc62e81be68ca07cdeb21a1e0809763647f" ; fi && \
  if [ $TARGETARCH == "amd64" ] ; then export DOWNLOAD_SHA="0f54a5fc68f82d3dcb0d3310253f2259fef1902d48cfa0a8721b82803c575024" ; fi && \
  echo "${DOWNLOAD_SHA}  /usr/local/bin/sops" | sha256sum -c - && \
  chmod +x /usr/local/bin/sops && sops -v

# age (encryption for sops)
RUN curl -sLo age.tgz https://github.com/FiloSottile/age/releases/download/v1.0.0/age-v1.0.0-linux-${TARGETARCH}.tar.gz && \
  if [ $TARGETARCH == "arm64" ] ; then export DOWNLOAD_SHA="6c82aa1d406e5a401ec3bb344cd406626478be74d5ae628f192d907cd78af981" ; fi && \
  if [ $TARGETARCH == "amd64" ] ; then export DOWNLOAD_SHA="6414f71ce947fbbea1314f6e9786c5d48436ebc76c3fd6167bf018e432b3b669" ; fi && \
  echo "${DOWNLOAD_SHA}  age.tgz" | sha256sum -c - && \
  tar -xzf age.tgz && cp age/age /usr/local/bin && \
  chmod +x /usr/local/bin/age && age --version

RUN curl -sLo cue.tgz https://github.com/cue-lang/cue/releases/download/v0.4.2/cue_v0.4.2_linux_${TARGETARCH}.tar.gz && \
  if [ $TARGETARCH == "arm64" ] ; then export DOWNLOAD_SHA="6515c1f1b6fc09d083be533019416b28abd91e5cdd8ef53cd0719a4b4b0cd1c7" ; fi && \
  if [ $TARGETARCH == "amd64" ] ; then export DOWNLOAD_SHA="d43cf77e54f42619d270b8e4c1836aec87304daf243449c503251e6943f7466a" ; fi && \
  echo "${DOWNLOAD_SHA}  cue.tgz" | sha256sum -c - && \
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
ARG TARGETARCH

# Create the kapp-controller user in the root group, the home directory will be mounted as a volume
RUN echo "kapp-controller:x:1000:0:/home/kapp-controller:/usr/sbin/nologin" > /etc/passwd
# Give the root group write access to openssh's root bundle so we can append custom roots at runtime
RUN chmod g+w /etc/pki/tls/certs/ca-bundle.crt

# fetchers
COPY --from=0 /helm-v2-unpacked/linux-${TARGETARCH}/helm helmv2
COPY --from=0 /helm-unpacked/linux-${TARGETARCH}/helm .
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
