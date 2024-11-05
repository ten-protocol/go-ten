package common

import "time"

// Config contains the configuration required by the WalletExtension.
type Config struct {
	WalletExtensionHost            string
	WalletExtensionPortHTTP        int
	WalletExtensionPortWS          int
	NodeRPCHTTPAddress             string
	NodeRPCWebsocketAddress        string
	LogPath                        string
	DBPathOverride                 string // Overrides the database file location. Used in tests.
	VerboseFlag                    bool
	DBType                         string
	DBConnectionURL                string
	TenChainID                     int
	StoreIncomingTxs               bool
	RateLimitUserComputeTime       time.Duration
	RateLimitWindow                time.Duration
	RateLimitMaxConcurrentRequests int
	Debug                          bool
}
