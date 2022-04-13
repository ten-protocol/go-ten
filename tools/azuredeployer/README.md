# Obscuro enclave service Azure deployer

Obscuro uses [EGo](https://www.ego.dev/) to run the Obscuro enclave service inside SGX. EGo supports running the 
enclave in simulation mode using the `OE_SIMULATION=1` environment variable.

This tool automates the creation of an SGX-enabled Azure VM that is set up to run the enclave service outside of 
simulation mode.

## Usage

* Install the Azure CLI by following the instructions [here](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
* Set up file-based authentication by following the instructions [here](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-file-based-authentication)
* Choose values for the `vm_password` and `email_recipient` fields in the `vm-params.json` file
  * `vm_password` will be the Azure VM's SSH password. It must meet the requirements set out [here](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/faq#what-are-the-password-requirements-when-creating-a-vm-)
  * `email_recipient` will be the email used for notifications
* Create the VM using the `main` method under `main/azure_deployer.go`
* SSH into the VM using the `obscuro` user (e.g. `ssh obscuro@XX.XX.XXX.XXX`)
* Start the enclave service outside of simulation mode using the following command, replacing `$1` with an integer 
  representing the 20 bytes of the node's address:

      sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p 11000:11000/tcp obscuro_enclave --nodeID $1 --address :11000

The enclave service is now running and exposed on port 11000.
