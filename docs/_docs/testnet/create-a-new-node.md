---
---
# Create a New Node
Follow the steps below to create a new Obscuro node using Docker and join the Obscuro Testnet.

Obscuro makes use of Intel's Software Guard Extension (SGX) capability to achieve full computational privacy (see the [Obscuro Whitepaper](https://whitepaper.obscu.ro) for full details). As a result there are some hardware compatibility requirements and the Intel SGX driver and SGX Platform Software (PSW) has to be installed.

An attestation process initiated when the node wants to join the Obscuro testnet verifies the node is running a genuine SGX enclave and that it is patched and not vulnerable to any known exploits. The steps below include a means of checking the outcome of the attestation process to avoid a failure at the point of joining the Testnet.

## Prerequisites
* Ubuntu 18.04
* URL to the most recent SGX driver .bin file, available from https://download.01.org/intel-sgx/sgx-linux/. For example https://download.01.org/intel-sgx/sgx-linux/2.16/distro/ubuntu18.04-server/sgx_linux_x64_driver_2.11.054c9c4c.bin

Either,
* Computer hardware running an Intel Xeon CPU with Intel SGX capability (use [Intel's Production Specification advanced search page](https://ark.intel.com/content/www/us/en/ark/search/featurefilter.html?productType=873) to confirm compatible processors)
* Enable Intel Software Guard Extension in the BIOS menu
* Disable Secure Boot in the BIOS menu

Or,
* Azure Confidential Computing virtual machine. By default these use Intel SGX capable processors.

# Joining Testnet as an Aggregator Node
## Install SGX
1. Update the system and install required components:

     ```
     sudo apt update
     sudo apt upgrade
     sudo apt-get install make gcc wget
     ```

1. Download SGX driver:

     `wget "<the URL recorded as part of prerequisites above>"`

1. Set protections to allow for the .bin file execution:

     `chmod +x sgx_linux_x64_driver_<driver version taken from driver download URL>.bin`

1. Install the driver:

     `sudo ./sgx_linux_x64_driver_<driver version taken from driver download URL>.bin`

1. Create or update remount-dev-exec.service to remove /dev as exec and at system startup:

     `sudo nano /etc/systemd/system/remount-dev-exec.service`

     Paste in the following and press Ctrl+O to write the contents and Ctrl+X to exit nano
     ```
          [Unit]
          Description=Remount /dev as exec to allow AESM service to boot and load enclaves into SGX

          [Service]
          Type=oneshot
          ExecStart=/bin/mount -o remount,exec /dev
          RemainAfterExit=true

          [Install]
          WantedBy=multi-user.target
     ```

     ```
     sudo systemctl enable remount-dev-exec
     sudo systemctl start remount-dev-exec
     ```

1. Configure the system to run an Intel SGX application:
     ``` 
     echo 'deb [arch=amd64] https://download.01.org/intelsgx/sgx_repo/ubuntu bionic main' | sudo tee /etc/apt/sources.list.d/intel-sgx.list
     wget -qO - https://download.01.org/intelsgx/sgx_repo/ubuntu/intel-sgx-deb.key | sudo apt-key add
     sudo apt update
     sudo apt-get install libsgx-epid libsgx-quote-ex libsgx-dcap-ql libsgx-enclave-common libsgx-urts sgx-aesm-service libsgx-uae-service autoconf libtool libprotobuf-dev
     ```
## Verify Your SGX Installation
_PLACEHOLDER_

## Install Docker
Before you install Docker Engine for the first time you need to set up the Docker repository.

### Set up the repository
1. Update the `apt` package index and install packages to allow `apt` to use a repository over HTTPS:
     ```
     sudo apt-get update
     sudo apt-get install ca-certificates curl gnupg lsb-release
     ```
1. Add Docker’s official GPG key:

     ```
     sudo mkdir -p /etc/apt/keyrings
     curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
     ```

1. Use the following command to set up the repository:
     ```
     echo \
     "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
     $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
     ```

### Install Docker Engine
1. Update the `apt` package index, and install the latest version of Docker Engine, containered, and Docker Compose, or go to the next step to install a specific version:
     ```
     sudo apt-get update
     sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
     ```
1. To install a specific version of Docker Engine, list the available versions in the repo, then select and install:

   a. List the versions available in your repo:

     `apt-cache madison docker-ce`

   b. Install a specific version using the version string from the second column, for example, `5:20.10.16~3-0~ubuntu-jammy`.

     `sudo apt-get install docker-ce=<VERSION_STRING> docker-ce-cli=<VERSION_STRING> containerd.io docker-compose-plugin`

1. Verify that Docker Engine is installed correctly by running the hello-world image.

     `sudo docker run hello-world`

This command downloads a test image and runs it in a container. When the container runs, it prints a message and exits.

Docker Engine is installed and running. The `docker` group is created but no users are added to it. You need to use `sudo` to run Docker commands. 

## Start Edgeless Database
1. Run EdgelessDB on an SGX-capable system:

     `docker run --name my-edb -p3306:3306 -p8080:8080 --privileged -v /dev/sgx:/dev/sgx -t ghcr.io/edgelesssys/edgelessdb-sgx-4gb`

## Start Enclave
1. Create a Dockerfile so the Docker image can be created:

     `sudo nano enclave.Dockerfile`

     Paste in the following and press Ctrl+O to write the contents and Ctrl+X to exit nano

     ```
     FROM ghcr.io/edgelesssys/ego-dev:latest

     RUN git clone https://github.com/obscuronet/go-obscuro
     RUN cd go-obscuro/go/enclave/main && ego-go build && ego sign main

     ENV OE_SIMULATION=1
     ENTRYPOINT ["ego", "run", "go-obscuro/go/enclave/main/main"]
     EXPOSE 11000
     ```

1. Create a Docker image for the enclave service to run in SGX:

     `sudo docker build -t obscuro_enclave -f enclave.Dockerfile .`

2. Run the Docker image as a container where `XXX` is the port on which to expose the enclave service's RPC endpoints on the local machine, and `YYY` is the public IP address of your node:

     `sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p XXX:11000/tcp obscuro_enclave --hostID YYY --address :11000 --willAttest=true`

## Start Host
_PLACEHOLDER_

## Join Your Node to Obscuro Testnet
_PLACEHOLDER_