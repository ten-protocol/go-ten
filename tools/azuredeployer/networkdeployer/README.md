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

Immediately upon creation, the VM will be running two Obscuro hosts, two Obscuro enclaves, and two Geth nodes on a new 
Ethereum network.
