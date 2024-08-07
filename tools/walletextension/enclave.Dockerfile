# Build Stages:
# build-base = downloads modules and prepares the directory for compilation. Based on the ego-dev image
# build-enclave = copies over the actual source code of the project and builds it using a compiler cache
# deploy = copies over only the enclave executable without the source
#          in a lightweight base image specialized for deployment and prepares the /data/ folder.

# Defaults to restricted flag mode
ARG TESTMODE=false

FROM ghcr.io/edgelesssys/ego-dev:v1.5.0 AS build-base

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

WORKDIR /home/ten/go-ten/tools/walletextension/main

# Build the enclave using the cross image build cache.
RUN --mount=type=cache,target=/root/.cache/go-build \
    ego-go build

# New build stage for compiling the enclave with restricted flags mode
FROM build-enclave as build-enclave-testmode-false
# Sign the enclave executable
RUN ego sign enclave.json

# New build stage for compiling the enclave without restricted flags mode
FROM build-enclave as build-enclave-testmode-true
# Sign the enclave executable
RUN ego sign enclave-test.json

# Tag the restricted mode as the current build
FROM build-enclave-testmode-${TESTMODE} as build-enclave

# Trigger a new build stage and use the smaller ego version:
FROM ghcr.io/edgelesssys/ego-deploy:v1.5.0

# Copy the binary and the entrypoint script
COPY --from=build-enclave \
    /home/ten/go-ten/tools/walletextension/main /home/ten/go-ten/tools/walletextension/main

# Copy migration files
COPY --from=build-enclave \
    /home/ten/go-ten/tools/walletextension/storage/database /home/ten/go-ten/tools/walletextension/storage/database

WORKDIR /home/ten/go-ten/tools/walletextension/main

# simulation mode is ACTIVE by default
ENV OE_SIMULATION=1

# Enable core dumps
RUN ulimit -c unlimited && mkdir -p /tmp/core-dump && chmod 777 /tmp/core-dump
ENV COREDUMP_LOCATION=/tmp/core-dump/core.%e.%p

EXPOSE 3000
