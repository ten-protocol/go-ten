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

#### - Install Docker & Go 1.17

```
sudo apt-get update \
    && curl -fsSL https://get.docker.com -o get-docker.sh && sh ./get-docker.sh \
    && sudo snap refresh && sudo snap install --channel=1.17 go --classic 
```

#### - Download Obscuro repo


Make sure to use the latest `<version>` at https://github.com/obscuronet/go-obscuro/tags.

```
  git clone --depth 1 -b <version> https://github.com/obscuronet/go-obscuro.git /home/obscuro/go-obscuro
```

#### - Start Obscuro Node

To start the obscuro node some information is required to populate the starting script.

- (host_public_p2p_addr) The external facing address of the network. Where outside peers will connect to. Must be open to outside connections.`curl https://ipinfo.io/ip` provides the external IP.
- (private_key) Private Key to issue transactions into the Layer 1
- (host_id) Public Key derived from the Private Key

```
go run /home/obscuro/go-obscuro/go/node/cmd \
  -is_genesis=false \
  -node_type="validator" \
  -is_sgx_enabled="true" \
  -host_id=PublicKeyAddress \
  -l1_host="testnet-gethnetwork.uksouth.azurecontainer.io" \
  -management_contract_addr=0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF \
  -message_bus_contract_addr=0xFD03804faCA2538F4633B3EBdfEfc38adafa259B \
  -private_key="PrivateKeyString" \
  -private_key="HOST:10000" \
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