# TEN Faucet

This tools contains a Faucet to allow allocation of ETH tokens within an TEN network. For more information 
on TEN see the [TEN repo](https://github.com/ten-protocol/go-ten) and [documentation](https://docs.ten.xyz/).

## Repository Structure
The top level structure of the tool is as below;

```
├── Dockerfile                 # Docker file to build container
├── README.md                  # This readme file
├── cmd                        # Source code for the CLI application
├── container_build.sh         # Build a local container
├── container_run.sh           # Run a local container
├── faucet                     # Source code for faucet implementation
├── go.mod                     # Golang dependency management 
└── go.sum                     # Goland dependency checksums
```

## Running a local container
To run a local container and run the Faucet use the below;

```bash
$ ./container_run.sh 
```

By default, when running locally the Faucet will connect to a local testnet started as described in the go-TEN 
project repo [readme](https://github.com/ten-protocol/go-ten#building-and-running-a-local-testnet). It will connect 
to the TEN node running within the local testnet on host `validator-host` and port `13010`. The Faucet opens 
on port `80` within the container, but maps port `8080` on the host machine to this.


## Allocating ETH to an EOA on a local testnet
Allocating ETH to an externally owned account is done through a POST command to the `/fund/eth` endpoint, where the 
data in the POST command specifies the address e.g. for the account `0x0d2166b7b3A1522186E809e83d925d7b0B6db084`

```bash
curl --location --request POST 'http://127.0.0.1:8080/fund/eth' \
--header 'Content-Type: application/json' \
--data-raw '{ "address":"0x0d2166b7b3A1522186E809e83d925d7b0B6db084" }'
```

