# Build Stages:
# build-base = downloads modules and prepares the directory for compilation. Based on the ego-dev image
# build-enclave = copies over the actual source code of the project and builds it using a compiler cache
# deploy = copies over only the enclave executable without the source
#          in a lightweight base image specialized for deployment and prepares the /data/ folder.

# Final container folder structure:
#   /home/obscuro/data                          contains working files for the enclave
#   /home/obscuro/go-obscuro/go/enclave/main    contains the executable for the enclave
#

FROM ghcr.io/edgelesssys/ego-dev:v1.5.3 AS build-base

# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN ego-go mod download


# Trigger new build stage for compiling the enclave
FROM build-base as build-enclave
COPY . .

WORKDIR /home/obscuro/go-obscuro/go/enclave/main

# Build the enclave using the cross image build cache.
RUN --mount=type=cache,target=/root/.cache/go-build \
    ego-go build

FROM build-enclave
# Sign the enclave executable
RUN ego sign enclave.json

# Trigger a new build stage and use the smaller ego version:
FROM ghcr.io/edgelesssys/ego-deploy:v1.5.3

# Copy just the binary for the enclave into this build stage
COPY --from=build-enclave \
    /home/obscuro/go-obscuro/go/enclave/main /home/obscuro/go-obscuro/go/enclave/main
    
WORKDIR /home/obscuro/go-obscuro/go/enclave/main

# simulation mode is ACTIVE by default
ENV OE_SIMULATION=1
EXPOSE 11000