package common

import (
	"encoding/json"
	"math/big"
	"time"

	tengenesis "github.com/ten-protocol/go-ten/go/enclave/genesis"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common/log"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
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
	HOC                ERC20 = "HOC"
	POC                ERC20 = "POC"
	HOCAddr                  = "f3a8bd422097bFdd9B3519Eaeb533393a1c561aC"
	pocAddr                  = "9802F661d17c65527D7ABB59DAAD5439cb125a67"
	bridgeAddr               = "deB34A740ECa1eC42C8b8204CBEC0bA34FDD27f3"
	hocOwnerKeyHex           = "6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682" // Used only in sim tests. Fine
	pocOwnerKeyHex           = "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8" // Used only in sim tests. Fine
	TestnetPrefundedPK       = "8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b" // The genesis main account private key.
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
func DefaultEnclaveConfig() *enclaveconfig.EnclaveConfig {
	return &enclaveconfig.EnclaveConfig{
		NodeID:                    "",
		HostAddress:               "127.0.0.1:10000",
		RPCAddress:                "127.0.0.1:11000",
		L1ChainID:                 1337,
		TenChainID:                443,
		WillAttest:                false, // todo (config) - attestation should be on by default before production release
		ManagementContractAddress: gethcommon.BytesToAddress([]byte("")),
		LogLevel:                  int(gethlog.LvlInfo),
		LogPath:                   log.SysOut,
		UseInMemoryDB:             true, // todo (config) - persistence should be on by default before production release
		EdgelessDBHost:            "",
		SqliteDBPath:              "",
		ProfilerEnabled:           false,
		MinGasPrice:               big.NewInt(params.InitialBaseFee),
		TenGenesis:                "",
		DebugNamespaceEnabled:     false,
		MaxBatchSize:              1024 * 55,
		MaxRollupSize:             1024 * 128,
		GasPaymentAddress:         gethcommon.HexToAddress("0xd6C9230053f45F873Cb66D8A02439380a37A4fbF"),
		BaseFee:                   new(big.Int).SetUint64(params.InitialBaseFee),

		// Due to hiding L1 costs in the gas quantity, the gas limit needs to be huge
		// Arbitrum with the same approach has gas limit of 1,125,899,906,842,624,
		// whilst the usage is small. Should be ok since execution is paid for anyway.
		GasLocalExecutionCapFlag:  300_000_000_000,
		GasBatchExecutionLimit:    30_000_000,
		RPCTimeout:                5 * time.Second,
		StoreExecutedTransactions: true,
		DecompressionLimit:        10 * 1024 * 1024,
	}
}

func TestnetGenesisJSON() string {
	accts := make([]tengenesis.Account, 1)
	amount, success := big.NewInt(0).SetString("7500000000000000000000000000000", 10)
	if !success {
		panic("failed to set big.Int from string")
	}
	accts[0] = tengenesis.Account{
		Address: gethcommon.HexToAddress("A58C60cc047592DE97BF1E8d2f225Fc5D959De77"),
		Amount:  amount,
	}
	gen := &tengenesis.Genesis{
		Accounts: accts,
	}

	genesisBytes, err := json.Marshal(gen)
	if err != nil {
		panic(err)
	}
	return string(genesisBytes)
}
