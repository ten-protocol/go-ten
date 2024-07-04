package eth2network

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type PosEth2Network interface {
	Start() error
	Stop() error
}

type PosImpl struct {
	scriptPath string
}

func NewPosEth2Network(scriptPath string) PosEth2Network {
	return &PosImpl{scriptPath: scriptPath}
}
func (n *PosImpl) Start() error {
	return runShellScript(n.scriptPath)
}

func (n *PosImpl) Stop() error {
	return stopProcesses()
}

func runShellScript(scriptPath string) error {
	cmd := exec.Command("/bin/bash", scriptPath)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start script: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("script execution failed: %w\nOutput: %s", err, out.String())
	}

	fmt.Printf("Script output: %s\n", out.String())
	return nil
}

func stopProcesses() error {
	// Find processes using the specified ports
	cmd := exec.Command("/bin/bash", "-c", "lsof -i :12000 -i :30303 -t")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to find processes: %w\nOutput: %s", err, out.String())
	}

	pids := strings.Fields(out.String())
	if len(pids) == 0 {
		return fmt.Errorf("no processes found on UDP:12000 or TCP:30303")
	}

	// Kill each process
	for _, pid := range pids {
		killCmd := exec.Command("kill", "-9", pid)
		err := killCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to kill process %s: %w", pid, err)
		}
		fmt.Printf("Killed process %s\n", pid)
	}

	return nil
}

func main() {
	scriptPath := filepath.Join(".", "pos_eth2_network.sh")

	network := NewPosEth2Network(scriptPath)
	err := network.Start()
	if err != nil {
		fmt.Printf("Failed to start the network: %v\n", err)
		return
	}

	// Example of waiting for the network to be ready
	time.Sleep(30 * time.Second)

	// Connect to the Geth node to check if it's running
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		fmt.Printf("Failed to connect to the Geth node: %v\n", err)
		return
	}
	defer client.Close()

	// Example: Get the latest block number
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		fmt.Printf("Failed to get block number: %v\n", err)
		return
	}
	fmt.Printf("Latest block number: %d\n", blockNumber)

	// Stop the network
	err = network.Stop()
	if err != nil {
		fmt.Printf("Failed to stop the network: %v\n", err)
		return
	}

	fmt.Println("Network stopped successfully")
}
