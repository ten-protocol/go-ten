# Build Stages:
# build-base = downloads modules and prepares the directory for compilation. Based on the ego-dev image
# build-enclave = copies over the actual source code of the project and builds it using a compiler cache
# deploy = copies over only the enclave executable without the source
#          in a lightweight base image specialized for deployment

# Final container folder structure:
#   /home/ten/go-ten/tools/walletextension/main    contains the executable for the enclave


FROM ghcr.io/edgelesssys/ego-dev:v1.5.3 AS build-base

# setup container data structure
RUN mkdir -p /home/ten/go-ten

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/ten/go-ten
COPY go.mod .
COPY go.sum .
RUN ego-go mod download


# Trigger new build stage for compiling the enclave
FROM build-base AS build-enclave
COPY . .

WORKDIR /home/ten/go-ten/tools/walletextension/main

# Build the enclave using the cross image build cache.
RUN --mount=type=cache,target=/root/.cache/go-build \
    ego-go build

# Sign the enclave executable
RUN ego sign enclave.json


# Trigger a new build stage and use the smaller ego version:
FROM ghcr.io/edgelesssys/ego-deploy:v1.5.3

# Copy just the binary for the enclave into this build stage
COPY --from=build-enclave \
    /home/ten/go-ten/tools/walletextension/main /home/ten/go-ten/tools/walletextension/main

WORKDIR /home/ten/go-ten/tools/walletextension/main

# simulation mode is ACTIVE by default
ENV OE_SIMULATION=1
EXPOSE 3000