package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
)

func main() {
	var (
		nodeURL      = flag.String("node-url", "http://localhost:13000", "Node URL")
		newCount     = flag.Int("historical-count", 0, "New historical transaction count (required)")
		deploymentID = flag.String("deployment-id", "", "Deployment ID for tracking (optional)")
	)
	flag.Parse()

	if *newCount == 0 {
		log.Fatal("historical-count is required and must be greater than 0")
	}

	// Create RPC client
	client, err := rpc.Dial(*nodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to node: %v", err)
	}
	defer client.Close()

	// Create ObsClient
	obsClient := obsclient.NewObsClient(client)

	// Get current historical count
	currentHistoricalCount, err := obsClient.GetHistoricalTransactionCount()
	if err != nil {
		log.Printf("Warning: Failed to get current historical count: %v", err)
	} else {
		fmt.Printf("Current historical transaction count: %d\n", currentHistoricalCount)
	}

	// Update the historical count directly in the database
	fmt.Printf("Updating historical transaction count to: %d\n", *newCount)

	// For now, we'll use a direct SQL update via RPC
	// In a real implementation, you might want to add a specific RPC method for this
	var result interface{}
	err = client.Call(&result, "scan_updateHistoricalTransactionCount", *newCount)
	if err != nil {
		log.Fatalf("Failed to update historical transaction count: %v", err)
	}

	fmt.Println("Historical transaction count updated successfully!")

	// Verify the update
	updatedCount, err := obsClient.GetHistoricalTransactionCount()
	if err != nil {
		log.Printf("Warning: Failed to verify updated count: %v", err)
	} else {
		fmt.Printf("Updated historical transaction count: %d\n", updatedCount)
	}

	// Record deployment if ID provided
	if *deploymentID != "" {
		fmt.Printf("Recording deployment: %s\n", *deploymentID)
		err = obsClient.RecordDeployment(*deploymentID)
		if err != nil {
			log.Printf("Warning: Failed to record deployment: %v", err)
		} else {
			fmt.Println("Deployment recorded successfully!")
		}
	}
}
