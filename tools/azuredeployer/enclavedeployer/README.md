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
* SSH into the VM using the `obscuro` user (e.g. `ssh obscuro@XX.XX.XXX.XXX`)
* Start the enclave service outside of simulation mode using the following command, replacing `$1` with an integer 
  representing the 20 bytes of the node's address:

      sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p 11000:11000/tcp obscuro_enclave --nodeID $1 --address :11000

The enclave service is now running and exposed on port 11000.

## Testing

The Obscuro enclave service can be tested by running the `TestOnAzureEnclaveNodesMonteCarloSimulation` test.

For this test, the value of `$1` in the enclave service start command must be set to `0`.

## Testing changes

Changes can be tested by checking out a branch on the VM or scp-ing the changed files over to the box. You will need to 
re-run the docker build from inside the obscuro-playground dir on the VM:

    sudo docker build -t obscuro_enclave -f dockerfiles/enclave_local_build.Dockerfile .

## Multiple nodes

You can run the azure_deployer tool multiple times to run more than one node by modifying the resource group name param
at the top of the script, e.g. you could add a number

    resourceGroupName     = "ObscuroEnclaveService2"

Note the IP logged out at the end of each run to SSH to the box (and to configure them into simulation_multi_azure_enclaves_test.go).
Don't forget to increment the nodeID parameter when you call the `docker run` command above.

Each of these resource groups will need to be deleted separately through the Azure dashboard webapp when you are finished with them.