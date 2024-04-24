package log

import (
	"os"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// These are the keys of the log entries
const (
	ErrKey           = "err"
	CtrErrKey        = "ctr_err"
	SubIDKey         = "subscription_id"
	CfgKey           = "cfg"
	TxKey            = "tx"
	DurationKey      = "duration"
	DurationMilliKey = "durationMilli"
	BundleHashKey    = "bundle"
	BatchHashKey     = "batch"
	BatchHeightKey   = "batch_height"
	BatchSeqNoKey    = "batch_seq_num"
	RollupHashKey    = "rollup"
	CmpKey           = "component"
	NodeIDKey        = "node_id"
	EnclaveIDKey     = "enclave_id"
	NetworkIDKey     = "network_id"
	BlockHeightKey   = "block_height"
	BlockHashKey     = "block_hash"
	PackageKey       = "package"
)

// Logging is grouped by the component where it was initialised
const (
	EnclaveCmp      = "enclave"
	HostCmp         = "host"
	HostRPCCmp      = "host_rpc"
	TxInjectCmp     = "tx_inject"
	TestLogCmp      = "test_log"
	P2PCmp          = "p2p"
	RPCClientCmp    = "rpc_client"
	DeployerCmp     = "deployer"
	NetwMngCmp      = "network_manager"
	WalletExtCmp    = "wallet_extension"
	TestGethNetwCmp = "test_geth_network"
	EthereumL1Cmp   = "l1_host"
	TenscanCmp      = "tenscan"
	CrossChainCmp   = "cross_chain"
)

// Used when the logger has to write to Sys.out
const (
	SysOut = "sys_out"
)

// New - helper function used to create a top level logger for a component.
func New(component string, level int, out string, ctx ...interface{}) gethlog.Logger {
	context := append(ctx, CmpKey, component)
	l := gethlog.New(context...)
	var s gethlog.Handler
	if out == SysOut {
		s = gethlog.StreamHandler(os.Stdout, TenLogFormat())
	} else {
		s1, err := gethlog.FileHandler(out, TenLogFormat())
		if err != nil {
			panic(err)
		}
		s = s1
	}
	l.SetHandler(gethlog.LvlFilterHandler(gethlog.Lvl(level), s))
	return l
}
