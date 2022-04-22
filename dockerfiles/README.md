# Obscuro enclave service Docker image

The Docker image defined by `enclave.Dockerfile` creates a Docker image for an Obscuro enclave service running in SGX. 
To build the image, run:

    docker build -t obscuro_enclave -f dockerfiles/enclave.Dockerfile .

This is a required step before running the Docker integration tests.

To run the image as a container, where `XXX` is the port on which to expose the enclave service's RPC endpoints on the 
local machine, and `YYY` is the address of the node that this enclave service is for as an integer:

    docker run -p XXX:11000/tcp obscuro_enclave --nodeID YYY --address :11000

By default, the image runs the Obscuro enclave service in SGX simulation mode. To run the enclave service in 
non-simulation mode instead, run:

    docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p XXX:11000/tcp obscuro_enclave --nodeID YYY --address :11000


Stop and remove all obscuro docker containers:

    docker rm $(docker stop $(docker ps -a -q --filter ancestor=obscuro_enclave --format="{{.ID}}"))