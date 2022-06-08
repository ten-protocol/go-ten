# Obscuro enclave service Azure deployer

This tool automates the creation of an SGX-enabled Azure VM that is set up to run the Obscuro enclave service outside 
of simulation mode.

## Usage

* Install the Azure CLI by following the instructions [here](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
* Set up file-based authentication by following the instructions [here](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-file-based-authentication)
* Set the `AZURE_AUTH_LOCATION` environment variable to where the authorization file is located
* Log into the Azure CLI using `az login`
* Choose values for the `vm_password` and `email_recipient` fields in the `vm-params.json` file
  * `vm_password` will be the Azure VM's SSH password. It must meet the requirements set out [here](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/faq#what-are-the-password-requirements-when-creating-a-vm-)
  * `email_recipient` will be the email used for notifications
* Create the VM using the `main` method in `enclave_deployer.go`
* SSH into the VM using the `obscuro` user (e.g. `ssh obscuro@XX.XX.XXX.XXX`), and the ip that was output by the enclave_deployer
* Repeat the previous steps until enough machines were created
* In the `TestOnAzureEnclaveNodesMonteCarloSimulation`, set the ips of the machines created
* Running it will output a docker command that has to run on each server 

## Testing changes

Changes can be tested by checking out a branch on the VM or scp-ing the changed files over to the box. You will need to 
re-run the docker build from inside the obscuro-playground dir on the VM:

    sudo docker build -t obscuro_enclave -f dockerfiles/enclave_local_build.Dockerfile .

## Multiple nodes

You can run the azure_deployer tool multiple times to run more than one node.

Note the IP logged out at the end of each run to SSH to the box (and to configure them into simulation_multi_azure_enclaves_test.go).

Each of these resource groups will need to be deleted separately through the Azure dashboard webapp when you are finished with them.

If you want to properly test the enclave attestation you should pass the `-willAttest` flag to the docker run commands on the VMs (add it to the commands output by the test).
This will only work if there are no local test enclaves basically only if len(vmIPs) == numberOfNodes).