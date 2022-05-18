# Obscuro network deployer

This tool automates the creation of an SGX-enabled Azure VM running a full Obscuro network and its associated L1 
network.

## Usage

* Install the Azure CLI by following the instructions [here](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
* Set up file-based authentication by following the instructions [here](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-file-based-authentication)
* Set the `AZURE_AUTH_LOCATION` environment variable to where the authorization file is located
* Log into the Azure CLI using `az login`
* Choose values for the `vm_password` and `email_recipient` fields in the `vm-params.json` file
  * `vm_password` will be the Azure VM's SSH password. It must meet the requirements set out [here](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/faq#what-are-the-password-requirements-when-creating-a-vm-)
  * `email_recipient` will be the email used for notifications
* Create the VM using the `main` method in `network_deployer.go`
* SSH into the VM using the `obscuro` user (e.g. `ssh obscuro@XX.XX.XXX.XXX`)
* Start the components (two Obscuro hosts, two Obscuro enclaves, and two Geth nodes) by running
  `sh obscuro-playground/tools/azuredeployer/networkdeployer/run.sh`
  * After running the script, if you run `ps`, you should see two `geth-v1.10.17` processes, two `enclave` processes, 
    and two `host` processes
* Each component type will produce logs:
  * The host logs are found under `~/host_logs.txt`
  * The enclave logs are found under `~/enclave_logs.txt`
  * The Geth network logs are found under `~/obscuro-playground/integration/.build/geth/<run-id>/node_logs.txt`

// todo - joel - say how this network is exposed
