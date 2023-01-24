FROM ghcr.io/edgelesssys/ego-dev:latest AS build-base
# on the container:
#   /home/obscuro/data       contains working files for the enclave
#   /home/obscuro/go-obscuro contains the src
#

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

# Sign the enclave executable
RUN ego sign main

# Trigger a new build stage and use the smaller ego version
FROM ghcr.io/edgelesssys/ego-deploy:latest

# Copy just the binary for the enclave into this build stage
COPY --from=build-enclave \
    /home/obscuro/go-obscuro/go/enclave/main home/obscuro/go-obscuro/go/enclave/main
    
WORKDIR /home/obscuro/go-obscuro/go/enclave/main
RUN mkdir -p /home/obscuro/data

# simulation mode is ACTIVE by default
ENV OE_SIMULATION=1
EXPOSE 11000