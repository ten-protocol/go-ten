package log

import (
	"os"

	gethlog "github.com/ethereum/go-ethereum/log"
)

var (
	// these are the keys of the log entries
	ErrKey   = "err"
	SubIdKey = "subId"
	CfgKey   = "cfg"
	TxKey    = "tx"
	CmpKey   = "cmp"

	// The high level components which will have their own logging context
	EnclaveCmp   = "enclave"
	HostCmp      = "host"
	HostRPCCmp   = "host_rpc"
	TxInjectCmp  = "tx_inject"
	TestLogCmp   = "test_log"
	P2PCmp       = "p2p"
	RPCClientCmp = "rpc_client"
	DeployerCmp  = "deployer"
	NetwMngCmp   = "network_manager"
	WalletExtCmp = "wallet_extension"

	TestGethNetwComp = "test_geth_network"
	EthereumL1Cmp    = "l1_host"

	NodeId    = "id"
	NetworkId = "id"

	// output type
	SysOut = "sys_out"
)

func New(component string, level int, out string, ctx ...interface{}) gethlog.Logger {
	context := append(ctx, CmpKey, component)
	l := gethlog.New(context...)
	var s gethlog.Handler
	if out == SysOut {
		s = gethlog.StreamHandler(os.Stdout, gethlog.TerminalFormat(false))
	} else {
		s1, err := gethlog.FileHandler(out, gethlog.TerminalFormat(false))
		if err != nil {
			panic(err)
		}
		s = s1
	}
	l.SetHandler(gethlog.LvlFilterHandler(gethlog.Lvl(level), s))
	return l
}
