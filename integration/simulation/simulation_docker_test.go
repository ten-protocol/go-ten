package simulation

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// TODO - Use individual Docker containers for the Obscuro nodes and Ethereum nodes.

var (
	enclaveDockerImg  = "obscuro_enclave"
	nodeIDFlag        = "--nodeID"
	addressFlag       = "--address"
	address           = ":11000"
	enclaveDockerPort = "11000/tcp"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes live in the same process, the enclaves run in individual Docker containers, and the Ethereum nodes are mocked out.
func TestDockerNodesMonteCarloSimulation(t *testing.T) {
	logFile := setupTestLog()
	defer logFile.Close()

	params := SimParams{
		NumberOfNodes:         10,
		NumberOfWallets:       5,
		AvgBlockDurationUSecs: uint64(250_000),
		SimulationTimeSecs:    15,
	}
	params.AvgNetworkLatency = params.AvgBlockDurationUSecs / 15
	params.AvgGossipPeriod = params.AvgBlockDurationUSecs / 3
	params.SimulationTimeUSecs = params.SimulationTimeSecs * 1000 * 1000

	efficiencies := EfficiencyThresholds{0.2, 0.3, 0.4}

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
		println(fmt.Sprintf("This test requires the `%s` Docker image to be built using the `dockerfiles/enclave` Dockerfile. Terminating.", enclaveDockerImg))
		return
	}

	// We create the Docker containers and set up a hook to terminate them at the end of the test.
	containerIDs := createDockerContainers(ctx, cli, params.NumberOfNodes)
	defer terminateDockerContainers(ctx, cli, containerIDs)

	// We start the Docker containers.
	for _, id := range containerIDs {
		if err = cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
	}

	testSimulation(t, CreateBasicNetworkOfDockerNodes, params, efficiencies)
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
		enclavePorts = append(enclavePorts, fmt.Sprintf("%d", enclaveStartPort+i))
	}

	containerIDs := make([]string, len(enclavePorts))
	for idx, port := range enclavePorts {
		nodeID := strconv.FormatInt(int64(idx+1), 10)
		containerConfig := &container.Config{Image: enclaveDockerImg, Cmd: []string{nodeIDFlag, nodeID, addressFlag, address}}
		hostConfig := &container.HostConfig{
			PortBindings: nat.PortMap{nat.Port(enclaveDockerPort): []nat.PortBinding{{HostIP: localhost, HostPort: port}}},
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
		timeout := -time.Nanosecond // A negative timeout means forceful termination.
		_ = cli.ContainerStop(ctx, id, &timeout)
		_ = cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
	}

	if err := cli.Close(); err != nil {
		panic(err)
	}
}
