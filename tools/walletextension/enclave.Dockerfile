# Build Stages:
# build-base = downloads modules and prepares the directory for compilation. Based on the ego-dev image
# build-enclave = copies over the actual source code of the project and builds it using a compiler cache
# deploy = copies over only the enclave executable without the source
#          in a lightweight base image specialized for deployment

# Final container folder structure:
#   /home/ten/go-ten/tools/walletextension/main    contains the executable for the enclave
#   /data                                          persistent volume mount point

# Trigger new build stage for compiling the enclave
FROM ghcr.io/edgelesssys/ego-dev:v1.8.0 AS build-base

# Install ca-certificates package and update it
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && update-ca-certificates

# Setup container data structure
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

# Run the complete Azure HSM setup (builds signer tool, signs binary, or skips if not needed)
RUN /home/ten/go-ten/tools/enclavesigner/AzureHSMSignatureScript.sh main /home/ten/go-ten/tools/enclavesigner/main

FROM ghcr.io/edgelesssys/ego-deploy:v1.8.0

# Create data directory that will be used for persistence
RUN mkdir -p /data && chmod 777 /data

# Copy just the binary for the enclave into this build stage
COPY --from=build-enclave \
    /home/ten/go-ten/tools/walletextension/main /home/ten/go-ten/tools/walletextension/main

# Copy the entry.sh script and make it executable
COPY tools/walletextension/main/entry.sh /home/ten/go-ten/tools/walletextension/main/entry.sh
RUN chmod +x /home/ten/go-ten/tools/walletextension/main/entry.sh

WORKDIR /home/ten/go-ten/tools/walletextension/main

# Add volume mount point
VOLUME ["/data"]

# simulation mode is ACTIVE by default
ENV OE_SIMULATION=1
EXPOSE 3000

# Set the entrypoint to entry.sh
ENTRYPOINT ["/home/ten/go-ten/tools/walletextension/main/entry.sh"]
