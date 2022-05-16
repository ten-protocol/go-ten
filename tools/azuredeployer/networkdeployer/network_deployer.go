package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/Azure/go-autorest/autorest"
	"golang.org/x/crypto/ssh"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

const (
	resourceGroupName     = "ObscuroNetwork"
	deploymentName        = "ObscuroNetwork"
	deploymentIPName      = "obscuro-network-ip"
	resourceGroupLocation = "uksouth"
	vmUsername            = "obscuro"

	templateFile   = "tools/azuredeployer/networkdeployer/vm-template.json"
	parametersFile = "tools/azuredeployer/networkdeployer/vm-params.json"
	vmPasswordKey  = "vm_password"

	azureAuthLocation = "AZURE_AUTH_LOCATION"
	subscriptionIDKey = "subscriptionId"

	sshPort     = "22"
	sshTimeout  = 5 * time.Second
	setupScript = `
		sudo apt-get update
		sudo apt-get install -y docker.io
 		sudo apt-get install -y make
		sudo apt-get install -y build-essential
		wget -c https://go.dev/dl/go1.18.2.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local

		if ! [ -d "obscuro-playground" ]; then git clone https://github.com/obscuronet/obscuro-playground; else :; fi
		sudo systemctl enable --now docker
		sudo docker build -t obscuro_enclave - < obscuro-playground/dockerfiles/enclave.Dockerfile

		cd obscuro-playground
		export PATH=$PATH:/usr/local/go/bin
		go mod tidy
		go build -o ./go/obscuronode/host/main/host ./go/obscuronode/host/main/main.go
		go build -o ./go/obscuronode/enclave/main/enclave ./go/obscuronode/enclave/main/main.go
		integration/gethnetwork/build_geth_binary.sh --version=v1.10.17
	`

	// todo - joel - provide script to start all the components, incl. Geth (start with one of each component)
)

// todo - joel - refactor below to share logic with enclave deployer where possible

func main() {
	ctx := context.Background()

	authorizer := getAuthorizer()
	authInfo := readJSON(os.Getenv(azureAuthLocation))
	groupsClient, deploymentsClient, addressClient := createClients((*authInfo)[subscriptionIDKey].(string), authorizer)

	createResourceGroup(ctx, groupsClient)
	createDeployment(ctx, deploymentsClient)
	vmIP := getIPAddress(ctx, addressClient)
	runSetupScript(vmIP)
}

// Authenticate with the Azure services using file-based authentication.
func getAuthorizer() autorest.Authorizer {
	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("failed to retrieve OAuth config: %v", err)
	}
	return authorizer
}

// Create the required clients for interacting with the Azure services.
func createClients(subscriptionID string, authorizer autorest.Authorizer) (resources.GroupsClient, resources.DeploymentsClient, network.PublicIPAddressesClient) {
	groupsClient := resources.NewGroupsClient(subscriptionID)
	groupsClient.Authorizer = authorizer

	deploymentsClient := resources.NewDeploymentsClient(subscriptionID)
	deploymentsClient.Authorizer = authorizer

	addressClient := network.NewPublicIPAddressesClient(subscriptionID)
	addressClient.Authorizer = authorizer

	return groupsClient, deploymentsClient, addressClient
}

// Create a resource group for the deployment.
func createResourceGroup(ctx context.Context, client resources.GroupsClient) {
	log.Printf("Creating resource group %s", resourceGroupName)

	group, err := client.CreateOrUpdate(
		ctx, resourceGroupName, resources.Group{Location: to.StringPtr(resourceGroupLocation)},
	)
	if err != nil {
		log.Fatalf("failed to create resource group: %v", err)
	}
	log.Printf("Created resource group %s", *group.Name)
}

// Create the deployment.
func createDeployment(ctx context.Context, client resources.DeploymentsClient) {
	log.Printf("Creating deployment %s", deploymentName)

	template := readJSON(templateFile)
	params := readJSON(parametersFile)

	deploymentFuture, err := client.CreateOrUpdate(
		ctx, resourceGroupName, deploymentName, resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: params,
				Mode:       resources.DeploymentModeIncremental,
			},
		},
	)
	if err != nil {
		log.Fatalf("failed to initiate deployment: %v", err)
	}

	err = deploymentFuture.FutureAPI.WaitForCompletionRef(ctx, client.BaseClient.Client)
	if err != nil {
		log.Fatalf("could not wait for deployment: %v", err)
	}
	result, err := deploymentFuture.Result(client)
	switch {
	case err != nil:
		log.Fatalf("failed to create deployment: %v", err)
	case result.Name == nil:
		log.Printf("Created deployment %s, but the provisioning state was not communicated back", deploymentName)
	default:
		log.Printf("Created deployment %s: %s", deploymentName, result.Properties.ProvisioningState)
	}
}

// Get the IP address of the deployment.
func getIPAddress(ctx context.Context, client network.PublicIPAddressesClient) string {
	ipAddress, err := client.Get(ctx, resourceGroupName, deploymentIPName, "")
	if err != nil {
		log.Fatalf("could not retrieve deployment's IP address")
	}

	return *ipAddress.PublicIPAddressPropertiesFormat.IPAddress
}

// Run the script to prepare the virtual machine for running the Obscuro network.
func runSetupScript(ipAddress string) {
	params := readJSON(parametersFile)
	vmPassword := (*params)[vmPasswordKey].(map[string]interface{})["value"].(string)

	config := ssh.ClientConfig{
		User:            vmUsername,
		Auth:            []ssh.AuthMethod{ssh.Password(vmPassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
		Timeout:         sshTimeout,
	}

	log.Printf("SSH'ing into VM to complete set-up: %s:%s", ipAddress, sshPort)
	var client *ssh.Client
	var err error
	for {
		client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", ipAddress, sshPort), &config)
		if err == nil {
			break
		}
		time.Sleep(sshTimeout)
		log.Printf("Waiting for VM to be ready...")
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("failed to create SSH session with VM: %v", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if err := session.Run(setupScript); err != nil {
		log.Fatalf("failed to run VM setup script: %v", err)
	}

	log.Printf("VM set-up complete. To connect, use configured password and run: ssh obscuro@%s", ipAddress)
}

// Read a JSON file into a map.
func readJSON(path string) *map[string]interface{} {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file at %s: %v", path, err)
	}
	contents := make(map[string]interface{})
	err = json.Unmarshal(data, &contents)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON from file at %s: %v", path, err)
	}
	return &contents
}
