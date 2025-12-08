package helpful

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
	"github.com/ten-protocol/go-ten/integration/simulation/devnetwork"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// TestBackupSharedSecret tests the backup shared secret functionality
// It generates a new key pair, sets the backup encryption key, calls BackupSharedSecret,
// decrypts the result to verify it's valid, and then starts a new node with the backed up secret
// to verify that the secret can be successfully provisioned to a new node
func TestBackupSharedSecret(t *testing.T) {
	// networktest.TestOnlyRunsInIDE(t)

	// Generate a new ECDSA key pair for the backup encryption
	backupPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate backup key: %v", err)
	}

	// Get the compressed public key
	backupPublicKeyBytes := crypto.CompressPubkey(&backupPrivateKey.PublicKey)
	backupPublicKeyHex := hex.EncodeToString(backupPublicKeyBytes)

	t.Logf("Generated backup key pair")
	t.Logf("Public key (hex): %s", backupPublicKeyHex)

	networktest.Run(
		"backup-shared-secret",
		t,
		env.LocalDevNetwork(devnetwork.WithBackupEncryptionKey(backupPublicKeyHex)),
		actions.Series(
			// Wait for the network to be initialized and for the shared secret to be generated
			actions.CreateAndFundTestUsers(1), // This ensures the network is up and running

			// Custom action to test the backup functionality
			actions.RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
				return testBackupSharedSecretRPC(ctx, network, backupPrivateKey, t)
			}),
		),
	)
}

func testBackupSharedSecretRPC(ctx context.Context, network networktest.NetworkConnector, backupPrivateKey *ecdsa.PrivateKey, t *testing.T) (context.Context, error) {
	// Connect to the sequencer RPC
	rpcAddress := network.SequencerRPCAddress()
	t.Logf("Connecting to RPC at: %s", rpcAddress)

	client, err := gethrpc.Dial(rpcAddress)
	if err != nil {
		return ctx, fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	// Call ten_backupSharedSecret
	var encryptedSecretHex string
	err = client.CallContext(ctx, &encryptedSecretHex, rpc.BackupSharedSecret)
	if err != nil {
		return ctx, fmt.Errorf("failed to call ten_backupSharedSecret: %w", err)
	}
	t.Logf("Retrieved encrypted shared secret (length: %d bytes)", len(encryptedSecretHex)/2)

	// Decode the hex-encoded encrypted secret
	encryptedSecret, err := hex.DecodeString(encryptedSecretHex[2:]) // Remove "0x" prefix if present
	if encryptedSecretHex[:2] != "0x" {
		encryptedSecret, err = hex.DecodeString(encryptedSecretHex)
	}
	if err != nil {
		return ctx, fmt.Errorf("failed to decode encrypted secret: %w", err)
	}

	// Decrypt the shared secret using our private key
	eciesPrivateKey := ecies.ImportECDSA(backupPrivateKey)
	decryptedSecret, err := eciesPrivateKey.Decrypt(encryptedSecret, nil, nil)
	if err != nil {
		return ctx, fmt.Errorf("failed to decrypt shared secret: %w", err)
	}

	t.Logf("Successfully decrypted shared secret (length: %d bytes)", len(decryptedSecret))

	// Verify the decrypted secret has the expected length (32 bytes for the shared secret)
	expectedSecretLength := 32
	if len(decryptedSecret) != expectedSecretLength {
		return ctx, fmt.Errorf("decrypted secret has unexpected length: expected %d, got %d", expectedSecretLength, len(decryptedSecret))
	}

	logger.Info("Backup shared secret test passed",
		"encryptedLength", len(encryptedSecret),
		"decryptedLength", len(decryptedSecret))

	t.Logf("✓ Backup shared secret test passed!")
	t.Logf("  - Encrypted secret length: %d bytes", len(encryptedSecret))
	t.Logf("  - Decrypted secret length: %d bytes", len(decryptedSecret))

	// Convert the decrypted secret to hex for provisioning
	decryptedSecretHex := hex.EncodeToString(decryptedSecret)
	t.Logf("Decrypted secret (hex): %s", decryptedSecretHex)

	// Now test starting a new node with the backed up secret
	devNetwork, ok := network.(*devnetwork.InMemDevNetwork)
	if !ok {
		return ctx, fmt.Errorf("expected InMemDevNetwork, got %T", network)
	}

	t.Logf("Starting new validator node with provisioned secret...")
	err = startNodeWithProvisionedSecret(ctx, devNetwork, decryptedSecretHex, t)
	if err != nil {
		return ctx, fmt.Errorf("failed to start node with provisioned secret: %w", err)
	}

	return ctx, nil
}

func startNodeWithProvisionedSecret(ctx context.Context, network *devnetwork.InMemDevNetwork, secretHex string, t *testing.T) error {
	// Get the current network configuration
	cfg := network.TenConfig()
	l1Data := network.L1SetupData()
	nodeWallets := network.Wallets()

	// Calculate the next node index (sequencer=0, validators=1..N)
	nextNodeIdx := 1 + network.NumValidators()
	t.Logf("Creating new validator node with index: %d", nextNodeIdx)

	// Create a new TenConfig with the provisioned secret
	newNodeConfig := &devnetwork.TenConfig{
		PortStart:           cfg.PortStart,
		InitNumValidators:   cfg.InitNumValidators,
		BatchInterval:       cfg.BatchInterval,
		RollupInterval:      cfg.RollupInterval,
		CrossChainInterval:  cfg.CrossChainInterval,
		NumNodes:            cfg.NumNodes,
		TenGatewayEnabled:   cfg.TenGatewayEnabled,
		NumSeqEnclaves:      cfg.NumSeqEnclaves,
		L1BeaconPort:        cfg.L1BeaconPort,
		L1BlockTime:         cfg.L1BlockTime,
		DeployerPK:          cfg.DeployerPK,
		BackupEncryptionKey: cfg.BackupEncryptionKey,
		SharedSecret:        secretHex, // Provision the backed up secret
	}

	// Create a new node operator with the provisioned secret
	l1Client := network.L1Network().GetClient(nextNodeIdx % network.L1Network().NumNodes())
	nodeWallet := nodeWallets.NodeWallets[nextNodeIdx%len(nodeWallets.NodeWallets)]

	newNode := devnetwork.NewInMemNodeOperator(
		nextNodeIdx,
		newNodeConfig,
		common.Validator,
		l1Data,
		l1Client,
		nodeWallet,
		network.Logger(),
	)

	// Start the new node
	t.Logf("Starting the new node...")
	err := newNode.Start()
	if err != nil {
		return fmt.Errorf("failed to start new node: %w", err)
	}

	// Wait for the node to be ready
	t.Logf("Waiting for the new node to be ready...")
	time.Sleep(5 * time.Second)

	// Verify the node is running by connecting to its RPC
	rpcAddress := newNode.HostRPCHTTPAddress()
	t.Logf("Connecting to new node RPC at: %s", rpcAddress)

	client, err := gethrpc.Dial(rpcAddress)
	if err != nil {
		// Try to stop the node before returning error
		_ = newNode.Stop()
		return fmt.Errorf("failed to connect to new node RPC: %w", err)
	}
	defer client.Close()

	// Test that the node is responsive by calling a simple RPC method
	var blockNumber string
	err = client.CallContext(ctx, &blockNumber, "eth_blockNumber")
	if err != nil {
		_ = newNode.Stop()
		return fmt.Errorf("failed to call eth_blockNumber on new node: %w", err)
	}

	t.Logf("✓ New node started successfully with provisioned secret!")
	t.Logf("  - Node RPC address: %s", rpcAddress)
	t.Logf("  - Current block number: %s", blockNumber)

	// Stop the new node to clean up
	t.Logf("Stopping the new node...")
	err = newNode.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop new node: %w", err)
	}

	t.Logf("✓ Test completed successfully!")

	return nil
}
