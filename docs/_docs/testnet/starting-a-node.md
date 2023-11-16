---
---
# Starting a node
How to start a node in the Evan's Cat testnet.

## Requirements
- SGX enabled VM
- Docker

## Steps
#### - Create an SGX enabled VM
Recommended Standard DC4s v2 (4 vcpus, 16 GiB memory) in Azure.

#### - Install Docker & Go 

```
sudo apt-get update \
    && curl -fsSL https://get.docker.com -o get-docker.sh && sh ./get-docker.sh \
    && sudo snap refresh && sudo snap install --channel=1.20 go --classic 
```

#### - Download Ten repo


Make sure to use the latest `<version>` at https://github.com/ten-protocol/go-ten/tags.

```
  git clone --depth 1 -b <version> https://github.com/ten-protocol/go-ten.git /home/ten/go-ten
```

#### - Start Ten Node

To start the ten node some information is required to populate the starting script.

- (host_public_p2p_addr) The external facing address of the network. Where outside peers will connect to. Must be open to outside connections.`curl https://ipinfo.io/ip` provides the external IP.
- (private_key) Private Key to issue transactions into the Layer 1
- (host_id) Public Key derived from the Private Key

```
go run /home/ten/go-ten/go/node/cmd \
     -is_genesis="false" \
     -node_type=validator \
     -is_sgx_enabled="true" \
     -host_id="0xD5C925bb6147aF6b6bB6086dC6f7B12faa1ab0ff" \
     -l1_host="testnet-eth2network.uksouth.azurecontainer.io" \
     -management_contract_addr="0x7d13152F196bDEebBD6CC53CD43e0CdAf97CbdE6" \
     -message_bus_contract_addr="0x426E82B481E2d0Bd6A1664Cccb24FFc76C0AD2f9" \
     -l1_start="0x190f89a5f68a880f1cd2a67e0ed17980c7f012503279a764acc78a538d7e188f" \
     -private_key="f19601351ab594b04f21bc1d577e03cc62290a5efea8198af8bdfb19dad035b3" \
     -sequencer_id="0xc272459070A881BfA28aB3E810f9b19E4F468531" \
     -host_public_p2p_addr="$(curl https://ipinfo.io/ip):10000" \
     -host_p2p_port=10000 \
     -enclave_docker_image="testnetobscuronet.azurecr.io/obscuronet/enclave:latest" \
     -host_docker_image="testnetobscuronet.azurecr.io/obscuronet/host:latest" \
     start
```

## - (Alternatively) Steps required to run a node on Alibaba SGX
Setup an Alibaba node to provide SGX to docker .

Instance Type g7t - Linux Distro : Alibaba Linux 3

### Install docker
```
sudo yum install -y yum-utils \
    && sudo yum-config-manager \
    --add-repo https://download.docker.com/linux/centos/docker-ce.repo \
    && sudo yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin \
    && sudo systemctl start docker
```

### Install sgx
```
yum install -y yum-utils && \
yum-config-manager --add-repo \
https://enclave-cn-hongkong.oss-cn-hongkong-internal.aliyuncs.com/repo/alinux/enclave-expr.repo \
&& yum install -y libsgx-ae-le libsgx-ae-pce libsgx-ae-qe3 libsgx-ae-qve \
libsgx-aesm-ecdsa-plugin libsgx-aesm-launch-plugin libsgx-aesm-pce-plugin \
libsgx-aesm-quote-ex-plugin libsgx-dcap-default-qpl libsgx-dcap-ql \
libsgx-dcap-quote-verify libsgx-enclave-common libsgx-launch libsgx-pce-logic \
libsgx-qe3-logic libsgx-quote-ex libsgx-ra-network libsgx-ra-uefi \
libsgx-uae-service libsgx-urts sgx-ra-service sgx-aesm-service \
&& yum install -y sgxsdk
```