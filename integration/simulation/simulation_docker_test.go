package simulation

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"strconv"
	"testing"
	"time"
)

// TODO - Use individual Docker containers for the Obscuro nodes and Ethereum nodes.

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes live in the same process, the enclaves run in individual Docker containers, and the Ethereum nodes are mocked out.
func TestDockerNodesMonteCarloSimulation(t *testing.T) {
	params := SimParams{
		NumberOfNodes:         3,
		NumberOfWallets:       5,
		AvgBlockDurationUSecs: uint64(250_000),
		SimulationTimeSecs:    15,
	}
	params.AvgNetworkLatency = params.AvgBlockDurationUSecs / 15
	params.AvgGossipPeriod = params.AvgBlockDurationUSecs / 3
	params.SimulationTimeUSecs = params.SimulationTimeSecs * 1000 * 1000
	efficiencies := EfficiencyThresholds{0.2, 0.3, 0.4}

	// todo - joel - when using api below, other enclaves are not initialised, only the first one is. is due to a mismatch in the address

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	var enclavePorts []string
	for i := 0; i < params.NumberOfNodes; i++ {
		// We assign an enclave port to each enclave service on the network.
		enclavePorts = append(enclavePorts, fmt.Sprintf("%d", enclaveStartPort+i))
	}

	var containerIDs []string
	for i, port := range enclavePorts {
		nodeID := strconv.FormatInt(int64(i+1), 10)
		containerConfig := &container.Config{Image: "obscuro_enclave", Cmd: []string{"--nodeID", nodeID}}
		hostConfig := &container.HostConfig{
			PortBindings: nat.PortMap{"11000/tcp": []nat.PortBinding{{"localhost", port}}},
		}

		resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
		if err != nil {
			panic(err)
		}
		containerIDs = append(containerIDs, resp.ID)
	}

	defer terminateDockerContainers(err, cli, ctx, containerIDs)

	for _, id := range containerIDs {
		if err = cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
	}

	testSimulation(t, CreateBasicNetworkOfDockerNodes, params, efficiencies)
}

func terminateDockerContainers(err error, cli *client.Client, ctx context.Context, containerIDs []string) {
	// todo - joel - not being called - why?
	time.Sleep(3 * time.Second)
	println("jjj calling defer block")
	for _, id := range containerIDs {
		println("jjj handling a container")
		timeout := 1 * time.Second
		err = cli.ContainerStop(ctx, id, &timeout)
		err = cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
	}
}
