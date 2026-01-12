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

	// HTTP endpoint rate limiting (for /join - can be used by other endpoints in the future too)
	// Burst is automatically set to 1.5x the rate for both limiters
	HTTPRateLimitGlobalRate float64 // Global requests per second (0 = disabled)
	HTTPRateLimitPerIPRate  float64 // Per-IP requests per second (0 = disabled)

	InsideEnclave                 bool // Indicates if the program is running inside an enclave
	EncryptionKeySource           string
	EnableTLS                     bool
	TLSDomain                     string
	EncryptingCertificateEnabled  bool
	DisableCaching                bool
	FrontendURL                   string        // Frontend URL allowed for restrictive CORS
	SessionKeyExpirationThreshold time.Duration // Threshold for session key expiration
	SessionKeyExpirationInterval  time.Duration // How often the session key expiration service runs
}
