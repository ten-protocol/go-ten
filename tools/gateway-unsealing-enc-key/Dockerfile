# Build Stages:
# build-base = downloads modules and prepares the directory for compilation. Based on the ego-dev image
# build-enclave = copies over the actual source code of the project and builds it using a compiler cache
# deploy = copies over only the enclave executable without the source
#          in a lightweight base image specialized for deployment and prepares the /data/ folder.

FROM ghcr.io/edgelesssys/ego-dev:v1.7.2 AS build-base

# setup container data structure
RUN mkdir -p /home/ten/go-ten

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/ten/go-ten
COPY go.mod .
COPY go.sum .
RUN ego-go mod download

# Trigger new build stage for compiling the enclave
FROM build-base as build-enclave
COPY . .

WORKDIR /home/ten/go-ten/tools/gateway-unsealing-enc-key/main

# Build the enclave using the cross image build cache.
RUN --mount=type=cache,target=/root/.cache/go-build \
    ego-go build

# New build stage for compiling the enclave with restricted flags mode
FROM build-enclave as sign-built-enclave
# Sign the enclave executable
RUN ego sign gw-unseal.json

# Trigger a new build stage and use the smaller ego version:
FROM ghcr.io/edgelesssys/ego-deploy:v1.7.2

# Copy the binary and the entrypoint script
COPY --from=sign-built-enclave \
    /home/ten/go-ten/tools/gateway-unsealing-enc-key/main /home/ten/go-ten/tools/gateway-unsealing-enc-key/main

WORKDIR /home/ten/go-ten/tools/gateway-unsealing-enc-key/main
