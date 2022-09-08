# Obscuro network deployer

This tool automates the creation of an SGX-enabled Azure VM running a full Obscuro network and its associated L1 
network.

## Deployment

* Install the Azure CLI by following the instructions [here](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
* Set up file-based authentication by following the instructions [here](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-file-based-authentication)
* Set the `AZURE_AUTH_LOCATION` environment variable to where the authorization file is located
* Log into the Azure CLI using `az login`
* Choose values for the `vm_password` and `email_recipient` fields in the `vm-params.json` file
  * `vm_password` will be the Azure VM's SSH password. It must meet the requirements set out [here](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/faq#what-are-the-password-requirements-when-creating-a-vm-)
  * `email_recipient` will be the email used for notifications
* Create the VM using the `main` method in `network_deployer.go`
* SSH into the VM using the `obscuro` user (e.g. `ssh obscuro@XX.XX.XXX.XXX`)
* Start the components (two Geth nodes, two Obscuro hosts, and two Obscuro enclave containers in Docker) by running
  `sh go-obscuro/tools/azuredeployer/networkdeployer/run.sh`
  * After running the script, if you run `ps`, you should see two `geth-v1.10.17` processes and two host processes
  * If you run `docker container ls`, you should see two Docker containers built from the `obscuro_enclave` image

## Usage

* Each component type will produce logs:
  * The host logs are found under `~/host_logs.txt`
  * The enclave logs can be viewed using `docker logs <container-id>`
  * The Geth network logs are found under `~/go-obscuro/integration/.build/geth/<run-id>/node_logs.txt`

* The first Obscuro host can be connected to remotely via WS on port `13001`
  * This repo contains a tool under `tools/obscuroclient/` that can be used to connect remotely to the running 
    Obscuro host and retrieve the current block head height
