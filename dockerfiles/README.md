# Obscuro enclave service Docker image

To create the Docker image for an Obscuro enclave service running in SGX simulation mode:

    docker build -t obscuro_enclave - < dockerfiles/enclave

To run the image as a container (where `XXXX` is the local port on which to expose the enclave):

    docker run -p XXXX:11000/tcp obscuro_enclave