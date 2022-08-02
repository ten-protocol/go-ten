package network

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/enclave/enclaverunner"
)

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
func createDockerContainers(ctx context.Context, client *client.Client, numOfNodes int, startPort int, mngmtCtrAddr string, erc20Addrs []string) map[string]string {
	containerIDs := map[string]string{}
	for i := 0; i < numOfNodes; i++ {
		nodeID := common.BigToAddress(big.NewInt(int64(i))).Hex()
		containerConfig := &container.Config{
			Image: enclaveDockerImg,
			Cmd: []string{
				"--" + enclaverunner.HostIDName, nodeID,
				"--" + enclaverunner.HostAddressName, fmt.Sprintf("%s:%d", Localhost, startPort+DefaultHostP2pOffset+i),
				"--" + enclaverunner.AddressName, enclaveAddress,
				"--" + enclaverunner.ManagementContractAddressName, mngmtCtrAddr,
				"--" + enclaverunner.Erc20ContractAddrsName, erc20Addrs[0] + "," + erc20Addrs[1],
				"--" + enclaverunner.ViewingKeysEnabledName + "=true",
			},
		}
		r := container.Resources{
			Memory:     2 * 1024 * 1024 * 1024, // 2GB
			MemorySwap: -1,
		}
		port := fmt.Sprintf("%d", startPort+DefaultEnclaveOffset+i)
		hostConfig := &container.HostConfig{
			PortBindings: nat.PortMap{nat.Port(enclaveDockerPort): []nat.PortBinding{{HostIP: Localhost, HostPort: port}}},
			Resources:    r,
		}

		resp, err := client.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
		if err != nil {
			panic(err)
		}
		containerIDs[resp.ID] = port
	}

	return containerIDs
}

// Stops and removes the test Docker containers.
func terminateDockerContainers(ctx context.Context, cli *client.Client, containerIDs map[string]string, containerStreams map[string]*types.HijackedResponse) {
	for id := range containerIDs {
		if containerStreams[id] != nil {
			containerStreams[id].Close()
		}
		err := cli.ContainerStop(ctx, id, nil)
		if err != nil {
			fmt.Printf("Could not stop the container %v : %s\n", id, err)
			continue
		}

		err = cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   false,
			Force:         true,
		})
		if err != nil {
			fmt.Printf("Could not remove the container %v : %s\n", id, err)
			continue
		}

		fmt.Printf("Stopped and removed container %v\n", id)
	}

	if err := cli.Close(); err != nil {
		fmt.Printf("Could not close cli: %s\n", err)
	}
}
