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

FROM ghcr.io/edgelesssys/ego-dev:v1.7.0 AS build-base

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
COPY dockerfiles/AzureKeyVaultSignatureScript.sh /tmp/
RUN chmod +x /tmp/AzureKeyVaultSignatureScript.sh

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
    else \
        echo "Skipping Azure setup"; \
    fi

# Add this section before your Azure Key Vault signature replacement step
# Install Azure CLI
RUN if [ -n "$AZURE_TENANT_ID" ]; then \
        echo "====== INSTALLING AZURE CLI ======" && \
        # Update package lists
        apt-get update && \
        # Install required packages for Azure CLI installation
        apt-get install -y \
            ca-certificates \
            curl \
            apt-transport-https \
            lsb-release \
            gnupg && \
        # Add Microsoft GPG key
        mkdir -p /etc/apt/keyrings && \
        curl -sLS https://packages.microsoft.com/keys/microsoft.asc | \
            gpg --dearmor | \
            tee /etc/apt/keyrings/microsoft.gpg > /dev/null && \
        chmod go+r /etc/apt/keyrings/microsoft.gpg && \
        # Add Azure CLI repository
        AZ_REPO=$(lsb_release -cs) && \
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/microsoft.gpg] https://packages.microsoft.com/repos/azure-cli/ $AZ_REPO main" | \
            tee /etc/apt/sources.list.d/azure-cli.list && \
        # Update package lists again with Azure CLI repo
        apt-get update && \
        # Install Azure CLI
        apt-get install -y azure-cli && \
        # Verify installation
        az --version && \
        echo "Azure CLI installation completed successfully" && \
        # Clean up to reduce image size
        apt-get clean && \
        rm -rf /var/lib/apt/lists/*; \
    else \
        echo "Skipping Azure setup"; \
    fi

RUN if [ -n "$AZURE_TENANT_ID" ]; then \
        set -x && \
        echo "====== STARTING AZURE KEY VAULT SIGNATURE REPLACEMENT ======" && \
        echo "====== STEP 1: EXTRACT HASH FROM EGO-SIGNED ENCLAVE ======" && \
        # Copy the originally signed enclave for backup
        cp main main.original && \
        # Extract the signing hash using the enclavesigner tool
        /tmp/enclavesigner extract main > /tmp/signing_data.json && \
        cat /tmp/signing_data.json && \
        # Get the hash that needs to be signed by Azure Key Vault
        hash_to_sign=$(cat /tmp/signing_data.json | jq -r '.hash_to_sign') && \
        echo "$hash_to_sign" > /tmp/hash_to_sign.txt && \
        echo "Hash to sign extracted: $hash_to_sign" && \
        echo "====== STEP 2: AZURE KEY VAULT SIGNING ======" && \
        # Set environment variables for the Azure script
        export AZURE_TENANT_ID="$AZURE_TENANT_ID" && \
        export AZURE_SUBSCRIPTION_ID="$AZURE_SUBSCRIPTION_ID" && \
        # Get signature from Azure Key Vault using the extracted hash
        /tmp/AzureKeyVaultSignatureScript.sh "$hash_to_sign" && \
        # Verify that Azure signing succeeded
        if [ ! -f /tmp/signature.txt ] || [ ! -s /tmp/signature.txt ]; then \
            echo "ERROR: Azure Key Vault signing failed - no signature file generated" && \
            exit 1; \
        fi && \
        echo "Using signature: $(head -c 50 /tmp/signature.txt)..." && \
        echo "====== STEP 3: REPLACE SIGNATURE IN ENCLAVE BINARY ======" && \
        # Apply the Azure Key Vault signature to the enclave binary using enclavesigner
        /tmp/enclavesigner replace main /tmp/signature.txt main.azure_signed && \
        # Verify the signature replacement worked - FAIL BUILD IF THIS FAILS
        echo "====== VERIFYING AZURE-SIGNED ENCLAVE ======" && \
        /tmp/enclavesigner verify main.azure_signed || { \
            echo "ERROR: Azure-signed enclave verification FAILED!" && \
            echo "The signature replacement did not produce a valid signed enclave" && \
            exit 1; \
        } && \
        # Replace the original with the Azure-signed version
        mv main.azure_signed main && \
        echo "Enclave signature successfully replaced with Azure Key Vault signature" && \
        echo "====== SIGNATURE REPLACEMENT COMPLETED ======" && \
        echo "====== FINAL VERIFICATION ======" && \
        # Final verification of the signed enclave - FAIL BUILD IF THIS FAILS
        /tmp/enclavesigner verify main || { \
            echo "ERROR: Final enclave verification FAILED!" && \
            echo "The Azure-signed enclave is not valid" && \
            exit 1; \
        } && \
        ls -la main && \
        echo "Build completed successfully with Azure Key Vault signature integration"; \
    else \
        echo "Skipping Azure setup"; \
    fi

# Trigger a new build stage and use the smaller ego version:
FROM ghcr.io/edgelesssys/ego-deploy:v1.7.0

# Copy just the binary for the enclave into this build stage
COPY --from=build-enclave \
    /home/obscuro/go-obscuro/go/enclave/main /home/obscuro/go-obscuro/go/enclave/main
    
WORKDIR /home/obscuro/go-obscuro/go/enclave/main

# simulation mode is ACTIVE by default
ENV OE_SIMULATION=1
EXPOSE 11000