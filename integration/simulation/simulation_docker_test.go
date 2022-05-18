package simulation

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"

	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// TODO - Use individual Docker containers for the Obscuro nodes and Ethereum nodes.

const (
	enclaveDockerImg  = "obscuro_enclave"
	nodeIDFlag        = "--nodeID"
	addressFlag       = "--address"
	enclaveAddress    = ":11000"
	enclaveDockerPort = "11000/tcp"
	dockerTestEnv     = "DOCKER_TEST_ENABLED"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes live in the same process, the enclaves run in individual Docker containers, and the Ethereum nodes are mocked out.
// $> docker rm $(docker stop $(docker ps -a -q --filter ancestor=obscuro_enclave --format="{{.ID}}") will stop and remove all images
func TestDockerNodesMonteCarloSimulation(t *testing.T) {
	if os.Getenv(dockerTestEnv) == "" {
		t.Skipf("set the variable to run this test: `%s=true`", dockerTestEnv)
	}

	setupTestLog("docker")

	simParams := params.SimParams{
		NumberOfNodes:             10,
		NumberOfObscuroWallets:    5,
		AvgBlockDuration:          400 * time.Millisecond,
		SimulationTime:            25 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.3,
		L2ToL1EfficiencyThreshold: 0.5,
		MgmtContractLib:           ethereum_mock.NewMgmtContractLibMock(),
		ERC20ContractLib:          ethereum_mock.NewStableTokenContractLibMock(),
	}
	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 20
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 2

	for i := 0; i < simParams.NumberOfNodes+1; i++ {
		simParams.EthWallets = append(simParams.EthWallets, datagenerator.RandomWallet())
	}

	// We create a Docker client.
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	// We check the required Docker images are available.
	if !dockerImagesAvailable(ctx, cli) {
		// We don't cause the test to fail here, because we want users to be able to run all the tests in the repo
		// without having to build the Docker images.
		println(fmt.Sprintf("This test requires the `%s` Docker image to be built using `dockerfiles/enclave.Dockerfile`. Terminating.", enclaveDockerImg))
		return
	}

	// We create the Docker containers and set up a hook to terminate them at the end of the test.
	containerIDs := createDockerContainers(ctx, cli, simParams.NumberOfNodes)
	defer terminateDockerContainers(ctx, cli, containerIDs)

	// We start the Docker containers.
	for _, id := range containerIDs {
		if err = cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
	}

	testSimulation(t, network.NewBasicNetworkOfNodesWithDockerEnclave(), &simParams)
}

// Checks the required Docker images exist.
func dockerImagesAvailable(ctx context.Context, cli *client.Client) bool {
	images, _ := cli.ImageList(ctx, types.ImageListOptions{})
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if strings.Contains(tag, enclaveDockerImg) {
				return true
			}
		}
	}
	return false
}

// Creates the test Docker containers.
func createDockerContainers(ctx context.Context, client *client.Client, numOfNodes int) []string {
	var enclavePorts []string
	for i := 0; i < numOfNodes; i++ {
		// We assign an enclave port to each enclave service on the network.
		enclavePorts = append(enclavePorts, fmt.Sprintf("%d", network.EnclaveStartPort+i))
	}

	containerIDs := make([]string, len(enclavePorts))
	for idx, port := range enclavePorts {
		nodeID := strconv.FormatInt(int64(idx+1), 10)
		containerConfig := &container.Config{Image: enclaveDockerImg, Cmd: []string{nodeIDFlag, nodeID, addressFlag, enclaveAddress}}
		hostConfig := &container.HostConfig{
			PortBindings: nat.PortMap{nat.Port(enclaveDockerPort): []nat.PortBinding{{HostIP: network.Localhost, HostPort: port}}},
		}

		resp, err := client.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
		if err != nil {
			panic(err)
		}
		containerIDs[idx] = resp.ID
	}

	return containerIDs
}

// Stops and removes the test Docker containers.
func terminateDockerContainers(ctx context.Context, cli *client.Client, containerIDs []string) {
	for _, id := range containerIDs {
		// timeout := -time.Nanosecond // A negative timeout means forceful termination.
		err1 := cli.ContainerStop(ctx, id, nil)
		if err1 != nil {
			fmt.Printf("Could not stop the container %v : %s\n", id, err1)
			continue
		}

		err2 := cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   false,
			Force:         true,
		})
		if err2 != nil {
			fmt.Printf("Could not remove the container %v : %s\n", id, err2)
			continue
		}

		fmt.Printf("Stopped and removed container %v\n", id)
	}

	if err := cli.Close(); err != nil {
		panic(err)
	}
}
