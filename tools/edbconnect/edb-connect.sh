#!/bin/bash

# Variables
IMAGE_NAME="testnetobscuronet.azurecr.io/obscuronet/edbconnect:latest"
CONTAINER_BASE_NAME="edb-connect"
UNIQUE_ID=$(date +%s%3N) # Using milliseconds for uniqueness
CONTAINER_NAME="${CONTAINER_BASE_NAME}-${UNIQUE_ID}"
VOLUME_NAME="obscuronode-enclave-volume"
NETWORK_NAME="node_network"
SGX_ENCLAVE_DEVICE="/dev/sgx_enclave"
SGX_PROVISION_DEVICE="/dev/sgx_provision"
COMMAND="ego run /home/ten/go-ten/tools/edbconnect/main/main"

# Function to destroy exited containers matching the base name
destroy_exited_containers() {
    exited_containers=$(sudo docker ps -a -q -f name=${CONTAINER_BASE_NAME} -f status=exited)
    if [ "$exited_containers" ];then
        echo "Removing exited containers matching ${CONTAINER_BASE_NAME}..."
        sudo docker rm $exited_containers || true
    else
        echo "No exited containers to remove."
    fi
}

# Destroy exited containers that match the base name
destroy_exited_containers

# Pull the latest image from Azure Docker repository
echo "Pulling the latest Docker image..."
sudo docker pull $IMAGE_NAME

# Run the container with the specified command
echo "Running the new container with name ${CONTAINER_NAME}..."
sudo docker run --name $CONTAINER_NAME \
  --network $NETWORK_NAME \
  -v $VOLUME_NAME:/enclavedata \
  --device $SGX_ENCLAVE_DEVICE:$SGX_ENCLAVE_DEVICE:rwm \
  --device $SGX_PROVISION_DEVICE:$SGX_PROVISION_DEVICE:rwm \
  -it $IMAGE_NAME $COMMAND

# After the REPL exits, destroy the container
echo "Destroying the container ${CONTAINER_NAME} after command exits..."
sudo docker rm $CONTAINER_NAME || true