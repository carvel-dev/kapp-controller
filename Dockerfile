FROM golang:1.13
WORKDIR /go/src/github.com/k14s/kapp-controller/

RUN apt-get -y update && apt-get install -y ca-certificates && update-ca-certificates

# k14s
RUN bash -c "wget -O- https://k14s.io/install.sh | bash"

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
