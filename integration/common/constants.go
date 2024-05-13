package common

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/wallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// The Contract addresses are the result of the deploying a smart contract from hardcoded owners.
// The "owners" are keys which are the de-facto "admins" of those erc20s and are able to transfer or mint tokens.
// The contracts and addresses cannot be random for now, because there is hardcoded logic in the core
// to generate synthetic "transfer" transactions for each erc20 deposit on ethereum
// and these transactions need to be signed. Which means the platform needs to "own" ERC20s.

// ERC20 - the supported ERC20 tokens. A list of made-up tokens used for testing.
type ERC20 string

const (
	HOC            ERC20 = "HOC"
	POC            ERC20 = "POC"
	HOCAddr              = "f3a8bd422097bFdd9B3519Eaeb533393a1c561aC"
	pocAddr              = "9802F661d17c65527D7ABB59DAAD5439cb125a67"
	bridgeAddr           = "deB34A740ECa1eC42C8b8204CBEC0bA34FDD27f3"
	hocOwnerKeyHex       = "6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
	pocOwnerKeyHex       = "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
)

var HOCOwner, _ = crypto.HexToECDSA(hocOwnerKeyHex)

// HOCContract - address of the deployed "hocus" erc20 on the L2
var HOCContract = gethcommon.BytesToAddress(gethcommon.Hex2Bytes(HOCAddr))

var POCOwner, _ = crypto.HexToECDSA(pocOwnerKeyHex)

// POCContract - address of the deployed "pocus" erc20 on the L2
var POCContract = gethcommon.BytesToAddress(gethcommon.Hex2Bytes(pocAddr))

// BridgeAddress - address of the virtual bridge
var BridgeAddress = gethcommon.BytesToAddress(gethcommon.Hex2Bytes(bridgeAddr))

// ERC20Mapping - maps an L1 Erc20 to an L2 Erc20 address
type ERC20Mapping struct {
	Name ERC20

	// L1Owner   wallet.Wallet
	L1Address *gethcommon.Address

	Owner     wallet.Wallet // for now the wrapped L2 version is owned by a wallet, but this will change
	L2Address *gethcommon.Address
}

// DefaultEnclaveConfig returns an EnclaveConfig with default values.
func DefaultEnclaveConfig() *config.EnclaveConfig {
	return &config.EnclaveConfig{
		HostID:                    gethcommon.BytesToAddress([]byte("")),
		HostAddress:               "127.0.0.1:10000",
		Address:                   "127.0.0.1:11000",
		NodeType:                  common.Sequencer,
		L1ChainID:                 1337,
		ObscuroChainID:            443,
		WillAttest:                false, // todo (config) - attestation should be on by default before production release
		ValidateL1Blocks:          false,
		GenesisJSON:               nil,
		ManagementContractAddress: gethcommon.BytesToAddress([]byte("")),
		LogLevel:                  int(gethlog.LvlInfo),
		LogPath:                   log.SysOut,
		UseInMemoryDB:             true, // todo (config) - persistence should be on by default before production release
		EdgelessDBHost:            "",
		SqliteDBPath:              "",
		ProfilerEnabled:           false,
		MinGasPrice:               big.NewInt(params.InitialBaseFee),
		ObscuroGenesis:            "",
		DebugNamespaceEnabled:     false,
		MaxBatchSize:              1024 * 55,
		MaxRollupSize:             1024 * 64,
		GasPaymentAddress:         gethcommon.HexToAddress("0xd6C9230053f45F873Cb66D8A02439380a37A4fbF"),
		BaseFee:                   new(big.Int).SetUint64(params.InitialBaseFee),

		// Due to hiding L1 costs in the gas quantity, the gas limit needs to be huge
		// Arbitrum with the same approach has gas limit of 1,125,899,906,842,624,
		// whilst the usage is small. Should be ok since execution is paid for anyway.
		GasLocalExecutionCapFlag: 300_000_000_000,
		GasBatchExecutionLimit:   300_000_000_000,
		RPCTimeout:               5 * time.Second,
	}
}
