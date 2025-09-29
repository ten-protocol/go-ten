package common

import "time"

// Config contains the configuration required by the WalletExtension.
type Config struct {
	WalletExtensionHost     string
	WalletExtensionPortHTTP int
	WalletExtensionPortWS   int

	NodeRPCHTTPAddress      string
	NodeRPCWebsocketAddress string

	LogPath        string
	LogLevel       int    // Log level for the application (0=critical, 1=error, 2=warn, 3=info, 4=debug)
	DBPathOverride string // Overrides the database file location. Used in tests.

	DBType          string
	DBConnectionURL string

	TenChainID       int
	StoreIncomingTxs bool

	RateLimitUserComputeTime       time.Duration
	RateLimitWindow                time.Duration
	RateLimitMaxConcurrentRequests int

	InsideEnclave       bool // Indicates if the program is running inside an enclave
	EncryptionKeySource string

	// Azure Managed HSM / Key Vault integration (optional)
	AzureEnableHSM               bool   // Enable Azure HSM secret fallback/write
	AzureHSMName                 string // Managed HSM name (e.g. EnclaveSigningV2)
	AzureKeyName                 string // Secret/Key name to store/read encryption key
	AzureResourceGroup           string // Resource group for convenience when using CLI auth
	AzureSubscriptionID          string // Subscription ID for CLI auth
	AzureTenantID                string // Tenant ID for CLI auth
	AzureReadSecret              bool   // If true, try reading the secret when other methods fail
	AzureWriteSecret             bool   // If true, write the generated key to Azure HSM
	EnableTLS                    bool
	TLSDomain                    string
	EncryptingCertificateEnabled bool
	DisableCaching               bool
	FrontendURL                  string // Frontend URL allowed for restrictive CORS
}
