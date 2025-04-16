package l1contractdeployer

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/ten-protocol/go-ten/go/node"
)

// StoreNetworkCfgInKeyVault stores the network configuration in the Azure Key Vault.
// It requires credentials to be set in the environment, see: https://docs.microsoft.com/en-us/azure/key-vault/general/authentication?tabs=azure-cli#set-environment-variables
// Note: these details are not secrets, but it is convenient to store them in KV alongside network secrets for infra systems access
func StoreNetworkCfgInKeyVault(ctx context.Context, vaultURL string, env string, networkConfig *node.NetworkConfig) error {
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

	// envPrefix allows us to use the same KV instance for multiple environments, if env="uat" then prefix="UAT-
	envPrefix := ""
	if env != "" {
		envPrefix = strings.ToUpper(env) + "-"
	}

	// Store each contract address as a secret (matching env var config names from TenConfig)
	secrets := map[string]string{
		envPrefix + "NETWORK_L1_CONTRACTS_NETWORKCONFIG":   networkConfig.NetworkConfigAddress,
		envPrefix + "NETWORK_L1_CONTRACTS_CROSSCHAIN":      networkConfig.CrossChainAddress,
		envPrefix + "NETWORK_L1_CONTRACTS_ROLLUP":          networkConfig.DataAvailabilityRegistryAddress,
		envPrefix + "NETWORK_L1_CONTRACTS_ENCLAVEREGISTRY": networkConfig.EnclaveRegistryAddress,
		envPrefix + "NETWORK_L1_CONTRACTS_MESSAGEBUS":      networkConfig.MessageBusAddress,
		envPrefix + "NETWORK_L1_STARTHASH":                 networkConfig.L1StartHash,
	}

	for name, value := range secrets {
		// key must not contain underscore, replace with hyphen
		name = strings.ReplaceAll(name, "_", "-")
		_, err := client.SetSecret(ctx, name, azsecrets.SetSecretParameters{
			Value: &value,
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to set secret %s: %w", name, err)
		}
	}

	return nil
}
