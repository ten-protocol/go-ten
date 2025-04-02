package l1contractdeployer

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/ten-protocol/go-ten/go/node"
)

// StoreNetworkCfgInKeyVault stores the network configuration in the Azure Key Vault.
// It requires credentials to be set in the environment, see: https://docs.microsoft.com/en-us/azure/key-vault/general/authentication?tabs=azure-cli#set-environment-variables
// Note: these details are not secrets, but it is convenient to store them in KV alongside network secrets for infra systems access
func StoreNetworkCfgInKeyVault(ctx context.Context, vaultURL string, networkConfig *node.NetworkConfig) error {
	// Create a credential using the default Azure credential chain
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to create credential: %w", err)
	}

	// Create a client
	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Store each contract address as a secret (matching env var config names from TenConfig)
	secrets := map[string]string{
		"NETWORK_L1_CONTRACTS_NETWORKCONFIG":   networkConfig.NetworkConfigAddress,
		"NETWORK_L1_CONTRACTS_CROSSCHAIN":      networkConfig.CrossChainAddress,
		"NETWORK_L1_CONTRACTS_ROLLUP":          networkConfig.DataAvailabilityRegistryAddress,
		"NETWORK_L1_CONTRACTS_ENCLAVEREGISTRY": networkConfig.EnclaveRegistryAddress,
		"NETWORK_L1_CONTRACTS_MESSAGEBUS":      networkConfig.MessageBusAddress,
		"NETWORK_L1_STARTHASH":                 networkConfig.L1StartHash,
	}

	for name, value := range secrets {
		_, err := client.SetSecret(ctx, name, azsecrets.SetSecretParameters{
			Value: &value,
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to set secret %s: %w", name, err)
		}
	}

	return nil
}
