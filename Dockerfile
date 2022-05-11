FROM golang:1.17.9 AS deps

ARG KCTRL_VER=development
WORKDIR /workspace

# carvel
COPY ./hack/install-deps.sh .
COPY ./hack/dependencies.yml .
RUN ./install-deps.sh

RUN curl -sLo /helm https://get.helm.sh/helm-v3.8.0-linux-amd64.tar.gz && \
  echo "8408c91e846c5b9ba15eb6b1a5a79fc22dd4d33ac6ea63388e5698d1b2320c8b  /helm" | sha256sum -c - && \
  mkdir /helm-unpacked && tar -C /helm-unpacked -xzvf /helm

# sops
RUN curl -sLo /usr/local/bin/sops https://github.com/mozilla/sops/releases/download/v3.7.2/sops-v3.7.2.linux && \
  echo "0f54a5fc68f82d3dcb0d3310253f2259fef1902d48cfa0a8721b82803c575024  /usr/local/bin/sops" | sha256sum -c - && \
  chmod +x /usr/local/bin/sops && sops -v

# age (encryption for sops)
RUN curl -sLo age.tgz https://github.com/FiloSottile/age/releases/download/v1.0.0/age-v1.0.0-linux-amd64.tar.gz && \
  echo "6414f71ce947fbbea1314f6e9786c5d48436ebc76c3fd6167bf018e432b3b669  age.tgz" | sha256sum -c - && \
  tar -xzf age.tgz && cp age/age /usr/local/bin && \
  chmod +x /usr/local/bin/age && age --version

RUN curl -sLo cue.tgz https://github.com/cue-lang/cue/releases/download/v0.4.2/cue_v0.4.2_linux_amd64.tar.gz && \
  echo "d43cf77e54f42619d270b8e4c1836aec87304daf243449c503251e6943f7466a cue.tgz" | sha256sum -c - && \
  tar -xf cue.tgz -C /usr/local/bin cue && cue version

# kapp-controller
COPY . .
# helpful ldflags reference: https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-X 'main.Version=$KCTRL_VER'" -trimpath -o controller ./cmd/main.go

# --- run image ---
FROM photon:4.0

# Install openssh for git
# TODO(bmo): why do we need sed?
RUN tdnf install -y git openssh-clients sed

# Create the kapp-controller user in the root group, the home directory will be mounted as a volume
RUN echo "kapp-controller:x:1000:0:/home/kapp-controller:/usr/sbin/nologin" > /etc/passwd
# Give the root group write access to the openssh's root bundle directory
# so we can rename the certs file with our dynamic config, and append custom roots at runtime
RUN chmod g+w /etc/pki/tls/certs

# fetchers
COPY --from=deps /helm-unpacked/linux-amd64/helm .
COPY --from=deps /usr/local/bin/imgpkg .
COPY --from=deps /usr/local/bin/vendir .

# templaters
COPY --from=deps /usr/local/bin/ytt .
COPY --from=deps /usr/local/bin/kbld .
COPY --from=deps /usr/local/bin/sops .
COPY --from=deps /usr/local/bin/age .
COPY --from=deps /usr/local/bin/cue .

# deployers
COPY --from=deps /usr/local/bin/kapp .

# Name it kapp-controller to identify it easier in process tree
COPY --from=deps /workspace/controller kapp-controller

# Copy the ca-bundle so we have an original
RUN cp /etc/pki/tls/certs/ca-bundle.crt /etc/pki/tls/certs/ca-bundle.crt.orig

# Run as kapp-controller by default, will be overridden to a random uid on OpenShift
USER 1000
ENV PATH="/:${PATH}"
ENTRYPOINT ["/kapp-controller"]
