# TEN enclave service Docker image

The Docker image defined by `enclave.Dockerfile` creates a Docker image for an TEN enclave service running in SGX. 
To build the image, run:

    docker build -t enclave -f dockerfiles/enclave.Dockerfile .

This is a required step before running the Docker integration tests.

To run the image as a container, where `XXX` is the port on which to expose the enclave service's RPC endpoints on the 
local machine, and `YYY` is the address of the node that this enclave service is for as an integer:

    docker run -p XXX:11000/tcp enclave --hostID YYY --address :11000

By default, the image runs the TEN enclave service in SGX simulation mode. To run the enclave service in 
non-simulation mode instead, run:

    docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p XXX:11000/tcp enclave --hostID YYY --address :11000 --willAttest=true

Stop and remove all TEN docker containers:

    docker rm $(docker stop $(docker ps -a -q --filter ancestor=enclave --format="{{.ID}}"))

## Data directory
The ego docker images setup a home directory in `/home/ten/` and within that `go-ten/` contains the code 
while `data/` is used for persistence by the enclave.

That `/home/ten/data` directory is mounted inside the enclave as `/data`, any files written to it by the enclave process
should be sealed using the private enclave key, so it will be unreadable to anyone that accesses the container.