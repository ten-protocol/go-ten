package nodescenario

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
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
	networktest.TestOnlyRunsInIDE(t)

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

			// extract the backup
			actions.RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
				rpcAddress := network.GetValidatorNode(1)
				client, err := gethrpc.Dial(rpcAddress.HostRPCHTTPAddress())
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

				return context.WithValue(ctx, "SharedSecret", hexutils.BytesToHex(decryptedSecret)), nil
			}),

			actions.SleepAction(15*time.Second),

			// stop all nodes
			actions.StopSequencerHost(),
			actions.StopSequencerEnclave(0),
			actions.StopValidatorEnclave(0),
			actions.StopValidatorEnclave(1),
			actions.StopValidatorEnclave(2),
			actions.StopValidatorHost(0),
			actions.StopValidatorHost(1),
			actions.StopValidatorHost(2),

			// wait for nodes to stop
			actions.SleepAction(4*time.Second),

			// start a brand new node with the shared secret configured
			actions.StartNewValidatorNode(actions.ApplySharedSecretFromContext()),
			actions.WaitForValidatorHealthCheck(devnetwork.DefaultTenConfig().InitNumValidators, 10*time.Second),
			actions.SleepAction(20*time.Second),
			actions.VerifyOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) error {
				// Get the first test user's address and check their balance
				user, err := actions.FetchTestUser(ctx, 0)
				if err != nil {
					return fmt.Errorf("failed to fetch test user: %w", err)
				}
				newValidator := network.GetValidatorNode(devnetwork.DefaultTenConfig().InitNumValidators)
				client, err := obsclient.DialWithAuth(newValidator.HostRPCHTTPAddress(), user.Wallet(), testlog.Logger())
				if err != nil {
					return fmt.Errorf("failed to connect to RPC: %w", err)
				}
				defer client.Close()

				// verify the node has a non-zero block height
				height, err := client.BatchNumber()
				if err != nil {
					return fmt.Errorf("failed to get batch number: %w", err)
				}
				t.Logf("New validator block height: %d", height)
				if height == 0 {
					return fmt.Errorf("expected non-zero block height on new validator")
				}

				// verify the chain state is accessible (funded user wallet has non-zero balance)
				balance, err := client.BalanceAt(ctx, nil)
				if err != nil {
					return fmt.Errorf("failed to get balance: %w", err)
				}
				t.Logf("Test user balance on new validator: %s", balance.String())
				if balance.Cmp(big.NewInt(0)) == 0 {
					return fmt.Errorf("expected non-zero balance for test user on new validator")
				}

				t.Logf("Backup shared secret successfully provisioned to new node and verified")
				return nil
			}),
		),
	)
}
