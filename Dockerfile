FROM golang:1.12
WORKDIR /go/src/github.com/k14s/kapp-controller/

RUN mkdir -p /tmp
RUN apt-get -y update && apt-get install -y git ca-certificates && update-ca-certificates
RUN adduser --disabled-login kapp-controller

RUN bash -c "wget -O- https://k14s.io/install.sh | bash"

RUN wget -O- https://get.helm.sh/helm-v2.14.3-linux-amd64.tar.gz > /helm && echo "38614a665859c0f01c9c1d84fa9a5027364f936814d1e47839b05327e400bf55  /helm" | shasum -c - && mkdir /helm-unpacked && tar -C /helm-unpacked -xzvf /helm

RUN wget -O- https://github.com/k14s/imgpkg/releases/download/v0.1.0/imgpkg-linux-amd64 > /imgpkg && echo "a9d0ba0edaa792d0aaab2af812fda85ca31eca81079505a8a5705e8ee1d8be93  /imgpkg" | shasum -c - && chmod +x /imgpkg

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o controller ./cmd/controller/...

RUN chmod 700 /tmp
RUN rm -rf /tmp/*

# Needs ubuntu for installing git/openssh
FROM ubuntu:bionic
RUN apt-get update && apt-get install -y git openssh-client && rm -rf /var/lib/apt/lists/*

# Needed for scratch but using ubuntu now
# COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/passwd /etc/passwd

COPY --from=0 /go/src/github.com/k14s/kapp-controller/controller .
COPY --from=0 /tmp /tmp

# fetchers
COPY --from=0 /helm-unpacked/linux-amd64/helm .
COPY --from=0 /imgpkg .

# templaters
COPY --from=0 /usr/local/bin/ytt .
COPY --from=0 /usr/local/bin/kbld .

# deployers
COPY --from=0 /usr/local/bin/kapp .

ENV PATH="/:${PATH}"
ENTRYPOINT ["/controller"] # TODO USER kapp-controller
