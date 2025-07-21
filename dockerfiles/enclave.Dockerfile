# for integration with the Azure Key Vault, it needs an azure "tenant id" and "subscription id" to be passed as build arguments
# The build process when using these arguments is interactive. It requires logging in on Azure and approving the login request.
ARG AZURE_TENANT_ID
ARG AZURE_SUBSCRIPTION_ID

# Build Stages:
# build-base = downloads modules and prepares the directory for compilation. Based on the ego-dev image
# build-enclave = copies over the actual source code of the project and builds it using a compiler cache
# deploy = copies over only the enclave executable without the source
#          in a lightweight base image specialized for deployment and prepares the /data/ folder.

# Final container folder structure:
#   /home/obscuro/data                          contains working files for the enclave
#   /home/obscuro/go-obscuro/go/enclave/main    contains the executable for the enclave
#

FROM ghcr.io/edgelesssys/ego-dev:v1.7.2 AS build-base

# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN ego-go mod download

# Trigger new build stage for compiling the enclave
FROM build-base as build-enclave

ARG AZURE_TENANT_ID
ARG AZURE_SUBSCRIPTION_ID

COPY . .

# Copy the Azure signing script
COPY dockerfiles/AzureHSMSignatureScript.sh /tmp/
RUN chmod +x /tmp/AzureHSMSignatureScript.sh

WORKDIR /home/obscuro/go-obscuro/go/enclave/main

# Build the enclave using the cross image build cache.
RUN --mount=type=cache,target=/root/.cache/go-build \
    ego-go build

# Build the enclavesigner tool
WORKDIR /home/obscuro/go-obscuro/tools/enclavesigner/main
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /tmp/enclavesigner .

# Return to enclave directory for signing process
WORKDIR /home/obscuro/go-obscuro/go/enclave/main

# Sign the enclave executable
RUN ego sign enclave.json

RUN if [ -n "$AZURE_TENANT_ID" ]; then \
        apt-get update && apt-get install -y jq; \
        echo "====== INSTALLING AZURE CLI ======" && \
        apt-get update && \
        apt-get install -y \
            ca-certificates \
            curl \
            apt-transport-https \
            lsb-release \
            gnupg && \
        mkdir -p /etc/apt/keyrings && \
        curl -sLS https://packages.microsoft.com/keys/microsoft.asc | \
            gpg --dearmor | \
            tee /etc/apt/keyrings/microsoft.gpg > /dev/null && \
        chmod go+r /etc/apt/keyrings/microsoft.gpg && \
        AZ_REPO=$(lsb_release -cs) && \
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/microsoft.gpg] https://packages.microsoft.com/repos/azure-cli/ $AZ_REPO main" | \
            tee /etc/apt/sources.list.d/azure-cli.list && \
        apt-get update && \
        apt-get install -y azure-cli && \
        az --version && \
        echo "Azure CLI installation completed successfully" && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/*; \
        set -x && \
        echo "====== STARTING AZURE KEY VAULT SIGNATURE REPLACEMENT ======" && \
        echo "====== STEP 1: EXTRACT HASH FROM EGO-SIGNED ENCLAVE ======" && \
        cp main main.original && \
        /tmp/enclavesigner extract_hash main > /tmp/hash.hex && \
        hash_to_sign=$(cat /tmp/hash.hex) && \
        echo "Hash to sign extracted: $hash_to_sign" && \
        echo "====== STEP 2: AZURE KEY VAULT SIGNING ======" && \
        export AZURE_TENANT_ID="$AZURE_TENANT_ID" && \
        export AZURE_SUBSCRIPTION_ID="$AZURE_SUBSCRIPTION_ID" && \
        /tmp/AzureHSMSignatureScript.sh "$hash_to_sign" && \
        if [ ! -f /tmp/signature.b64 ] || [ ! -s /tmp/signature.b64 ]; then \
            echo "ERROR: Azure Key Vault signing failed - no signature file generated" && \
            exit 1; \
        fi && \
        if [ ! -f /tmp/modulus.b64 ] || [ ! -s /tmp/modulus.b64 ]; then \
            echo "ERROR: Azure Key Vault signing failed - no modulus file generated" && \
            exit 1; \
        fi && \
        signature_b64=$(cat /tmp/signature.b64) && \
        modulus_b64=$(cat /tmp/modulus.b64) && \
        echo "Using signature: $(echo "$signature_b64" | head -c 50)..." && \
        echo "Using modulus: $(echo "$modulus_b64" | head -c 50)..." && \
        echo "====== STEP 3: REPLACE SIGNATURE IN ENCLAVE BINARY ======" && \
        /tmp/enclavesigner replace main "$signature_b64" "$modulus_b64" main.azure_signed  2>&1 && \
        echo "====== VERIFYING AZURE-SIGNED ENCLAVE ======" && \
        /tmp/enclavesigner verify main.azure_signed  2>&1 || { \
            echo "ERROR: Azure-signed enclave verification FAILED!" && \
            exit 1; \
        } && \
        mv main.azure_signed main && \
        echo "Enclave signature successfully replaced with Azure Key Vault signature" && \
        echo "====== SIGNATURE REPLACEMENT COMPLETED ======" && \
        echo "====== FINAL VERIFICATION ======" && \
        /tmp/enclavesigner verify main 2>&1 || { \
            echo "ERROR: Final enclave verification FAILED!" && \
            exit 1; \
        } && \
        ls -la main && \
        echo "Build completed successfully with Azure Key Vault signature integration"; \
    else \
        echo "Skipping Azure setup"; \
    fi

# Trigger a new build stage and use the smaller ego version:
FROM ghcr.io/edgelesssys/ego-deploy:v1.7.2

# Copy just the binary for the enclave into this build stage
COPY --from=build-enclave \
    /home/obscuro/go-obscuro/go/enclave/main /home/obscuro/go-obscuro/go/enclave/main
    
WORKDIR /home/obscuro/go-obscuro/go/enclave/main

# simulation mode is ACTIVE by default
ENV OE_SIMULATION=1
EXPOSE 11000