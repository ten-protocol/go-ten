# Obscuro enclave service Docker image

To create the Obscuro enclave service Docker image:

    docker build -t obscuro_enclave - < enclave

To run the image as a container (where `XXXX` is the local port on which to expose the enclave):

    docker run -p XXXX:11000/tcp obscuro_enclave