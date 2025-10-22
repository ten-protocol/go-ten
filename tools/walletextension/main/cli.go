package main

import (
	"flag"
	"fmt"
	"time"

	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const (
	// Flag names, defaults and usages.
	walletExtensionHostName    = "host"
	walletExtensionHostDefault = "127.0.0.1"
	walletExtensionHostUsage   = "The host where the wallet extension should open the port."
	walletExtensionPortName    = "port"
	walletExtensionPortDefault = 3000
	walletExtensionPortUsage   = "The port on which to serve the wallet extension. Default: 3000."

	walletExtensionPortWSName    = "portWS"
	walletExtensionPortWSDefault = 3001
	walletExtensionPortWSUsage   = "The port on which to serve websocket JSON RPC requests. Default: 3001."

	nodeHostName    = "nodeHost"
	nodeHostDefault = "erpc.sepolia-testnet.ten.xyz"
	nodeHostUsage   = "The host on which to connect to the Obscuro node. Default: `erpc.sepolia-testnet.ten.xyz`."

	nodeHTTPPortName    = "nodePortHTTP"
	nodeHTTPPortDefault = 80
	nodeHTTPPortUsage   = "The port on which to connect to the Obscuro node via RPC over HTTP. Default: 80."

	nodeWebsocketPortName    = "nodePortWS"
	nodeWebsocketPortDefault = 81
	nodeWebsocketPortUsage   = "The port on which to connect to the Obscuro node via RPC over websockets. Default: 81."

	logPathName    = "logPath"
	logPathDefault = "sys_out"
	logPathUsage   = "The path to use for the wallet extension's log file"

	databasePathName    = "databasePath"
	databasePathDefault = ".obscuro/gateway_database.db"
	databasePathUsage   = "The path for the wallet extension's database file. Default: .obscuro/gateway_database.db"

	logLevelFlagName    = "logLevel"
	logLevelFlagDefault = "info"
	logLevelFlagUsage   = "Log level for wallet extension (critical, error, warn, info, debug, trace). Default: info"

	dbTypeFlagName    = "dbType"
	dbTypeFlagDefault = "sqlite"
	dbTypeFlagUsage   = "Defined the db type (sqlite or mariaDB)"

	dbConnectionURLFlagName    = "dbConnectionURL"
	dbConnectionURLFlagDefault = ""
	dbConnectionURLFlagUsage   = "If dbType is set to mariaDB, this must be set. ex: obscurouser:password@tcp(127.0.0.1:3306)/ogdb"

	tenChainIDName      = "tenChainID"
	tenChainIDDefault   = 443
	tenChainIDFlagUsage = "ChainID of TEN network that the gateway is communicating with"

	storeIncomingTxs        = "storeIncomingTxs"
	storeIncomingTxsDefault = true
	storeIncomingTxsUsage   = "Flag to enable storing incoming transactions in the database for debugging purposes. Default: true"

	rateLimitUserComputeTimeName    = "rateLimitUserComputeTime"
	rateLimitUserComputeTimeDefault = 10 * time.Second
	rateLimitUserComputeTimeUsage   = "rateLimitUserComputeTime represents how much compute time is user allowed to used in rateLimitWindow time. If rateLimitUserComputeTime is set to 0, rate limiting is turned off. Default: 10s."

	rateLimitWindowName    = "rateLimitWindow"
	rateLimitWindowDefault = 1 * time.Minute
	rateLimitWindowUsage   = "rateLimitWindow represents time window in which we allow one user to use compute time defined with rateLimitUserComputeTimeMs  Default: 1m"

	rateLimitMaxConcurrentRequestsName    = "maxConcurrentRequestsPerUser"
	rateLimitMaxConcurrentRequestsDefault = 3
	rateLimitMaxConcurrentRequestsUsage   = "Number of concurrent requests allowed per user. Default: 3"

	insideEnclaveFlagName    = "insideEnclave"
	insideEnclaveFlagDefault = false
	insideEnclaveFlagUsage   = "Flag to indicate if the program is running inside an enclave. Default: false"

	encryptionKeySourceFlagName    = "encryptionKeySource"
	encryptionKeySourceFlagDefault = ""
	encryptionKeySourceFlagUsage   = "Source of the encryption key for the gateway database. It can be set to empty (read from the sealed key), URL of another gateway with which we can exchange the key or 'new' to generate a new key (but only if sealed key is not present). Default: empty"

	enableTLSFlagName    = "enableTLS"
	enableTLSFlagDefault = false
	enableTLSFlagUsage   = "Flag to enable TLS/HTTPS"

	tlsDomainFlagName    = "tlsDomain"
	tlsDomainFlagDefault = ""
	tlsDomainFlagUsage   = "Domain name for TLS certificate"

	encryptingCertificateEnabledFlagName    = "encryptingCertificateEnabled"
	encryptingCertificateEnabledFlagDefault = false
	encryptingCertificateEnabledFlagUsage   = "Flag to enable encrypting certificate functionality. Default: false"

	disableCachingFlagName    = "disableCaching"
	disableCachingFlagDefault = false
	disableCachingFlagUsage   = "Flag to disable response caching in the gateway. Default: false"

	frontendURLFlagName    = "frontendURL"
	frontendURLFlagDefault = "https://uat-gw-testnet.ten.xyz"
	frontendURLFlagUsage   = "The frontend URL that is allowed to access restricted CORS endpoints. Default: https://uat-gw-testnet.ten.xyz"

	sessionKeyExpirationThresholdFlagName    = "sessionKeyExpirationThreshold"
	sessionKeyExpirationThresholdFlagDefault = 10 * time.Minute
	sessionKeyExpirationThresholdFlagUsage   = "Threshold for session key expiration. Session keys older than this duration will be considered expired. If set to 0, session key expiration is disabled. Default: 24h"

	sessionKeyExpirationIntervalFlagName    = "sessionKeyExpirationInterval"
	sessionKeyExpirationIntervalFlagDefault = 10 * time.Second
	sessionKeyExpirationIntervalFlagUsage   = "How often the session key expiration service runs. Default: 2h"
)

// getLogLevelInt converts string log level to integer value
func getLogLevelInt(level string) int {
	switch level {
	case "critical":
		return 0
	case "error":
		return 1
	case "warn":
		return 2
	case "info":
		return 3
	case "debug":
		return 4
	case "trace":
		return 5
	default:
		return 3 // default to info level
	}
}

func parseCLIArgs() wecommon.Config {
	walletExtensionHost := flag.String(walletExtensionHostName, walletExtensionHostDefault, walletExtensionHostUsage)
	walletExtensionPort := flag.Int(walletExtensionPortName, walletExtensionPortDefault, walletExtensionPortUsage)
	walletExtensionPortWS := flag.Int(walletExtensionPortWSName, walletExtensionPortWSDefault, walletExtensionPortWSUsage)
	nodeHost := flag.String(nodeHostName, nodeHostDefault, nodeHostUsage)
	nodeHTTPPort := flag.Int(nodeHTTPPortName, nodeHTTPPortDefault, nodeHTTPPortUsage)
	nodeWebsocketPort := flag.Int(nodeWebsocketPortName, nodeWebsocketPortDefault, nodeWebsocketPortUsage)
	logPath := flag.String(logPathName, logPathDefault, logPathUsage)
	databasePath := flag.String(databasePathName, databasePathDefault, databasePathUsage)
	logLevel := flag.String(logLevelFlagName, logLevelFlagDefault, logLevelFlagUsage)
	dbType := flag.String(dbTypeFlagName, dbTypeFlagDefault, dbTypeFlagUsage)
	dbConnectionURL := flag.String(dbConnectionURLFlagName, dbConnectionURLFlagDefault, dbConnectionURLFlagUsage)
	tenChainID := flag.Int(tenChainIDName, tenChainIDDefault, tenChainIDFlagUsage)
	storeIncomingTransactions := flag.Bool(storeIncomingTxs, storeIncomingTxsDefault, storeIncomingTxsUsage)
	rateLimitUserComputeTime := flag.Duration(rateLimitUserComputeTimeName, rateLimitUserComputeTimeDefault, rateLimitUserComputeTimeUsage)
	rateLimitWindow := flag.Duration(rateLimitWindowName, rateLimitWindowDefault, rateLimitWindowUsage)
	rateLimitMaxConcurrentRequests := flag.Int(rateLimitMaxConcurrentRequestsName, rateLimitMaxConcurrentRequestsDefault, rateLimitMaxConcurrentRequestsUsage)
	insideEnclaveFlag := flag.Bool(insideEnclaveFlagName, insideEnclaveFlagDefault, insideEnclaveFlagUsage)
	encryptionKeySource := flag.String(encryptionKeySourceFlagName, encryptionKeySourceFlagDefault, encryptionKeySourceFlagUsage)
	enableTLSFlag := flag.Bool(enableTLSFlagName, enableTLSFlagDefault, enableTLSFlagUsage)
	tlsDomainFlag := flag.String(tlsDomainFlagName, tlsDomainFlagDefault, tlsDomainFlagUsage)
	encryptingCertificateEnabled := flag.Bool(encryptingCertificateEnabledFlagName, encryptingCertificateEnabledFlagDefault, encryptingCertificateEnabledFlagUsage)
	disableCaching := flag.Bool(disableCachingFlagName, disableCachingFlagDefault, disableCachingFlagUsage)
	frontendURL := flag.String(frontendURLFlagName, frontendURLFlagDefault, frontendURLFlagUsage)
	sessionKeyExpirationThreshold := flag.Duration(sessionKeyExpirationThresholdFlagName, sessionKeyExpirationThresholdFlagDefault, sessionKeyExpirationThresholdFlagUsage)
	sessionKeyExpirationInterval := flag.Duration(sessionKeyExpirationIntervalFlagName, sessionKeyExpirationIntervalFlagDefault, sessionKeyExpirationIntervalFlagUsage)
	flag.Parse()

	return wecommon.Config{
		WalletExtensionHost:            *walletExtensionHost,
		WalletExtensionPortHTTP:        *walletExtensionPort,
		WalletExtensionPortWS:          *walletExtensionPortWS,
		NodeRPCHTTPAddress:             fmt.Sprintf("%s:%d", *nodeHost, *nodeHTTPPort),
		NodeRPCWebsocketAddress:        fmt.Sprintf("%s:%d", *nodeHost, *nodeWebsocketPort),
		LogPath:                        *logPath,
		DBPathOverride:                 *databasePath,
		LogLevel:                       getLogLevelInt(*logLevel),
		DBType:                         *dbType,
		DBConnectionURL:                *dbConnectionURL,
		TenChainID:                     *tenChainID,
		StoreIncomingTxs:               *storeIncomingTransactions,
		RateLimitUserComputeTime:       *rateLimitUserComputeTime,
		RateLimitWindow:                *rateLimitWindow,
		RateLimitMaxConcurrentRequests: *rateLimitMaxConcurrentRequests,
		InsideEnclave:                  *insideEnclaveFlag,
		EncryptionKeySource:            *encryptionKeySource,
		EnableTLS:                      *enableTLSFlag,
		TLSDomain:                      *tlsDomainFlag,
		EncryptingCertificateEnabled:   *encryptingCertificateEnabled,
		DisableCaching:                 *disableCaching,
		FrontendURL:                    *frontendURL,
		SessionKeyExpirationThreshold:  *sessionKeyExpirationThreshold,
		SessionKeyExpirationInterval:   *sessionKeyExpirationInterval,
	}
}
