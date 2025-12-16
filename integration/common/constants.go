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
		NodeID:                          "",
		HostAddress:                     "127.0.0.1:10000",
		RPCAddress:                      "127.0.0.1:11000",
		L1ChainID:                       1337,
		TenChainID:                      443,
		WillAttest:                      false, // todo (config) - attestation should be on by default before production release
		NetworkConfigAddress:            gethcommon.BytesToAddress([]byte("")),
		DataAvailabilityRegistryAddress: gethcommon.BytesToAddress([]byte("")),
		EnclaveRegistryAddress:          gethcommon.BytesToAddress([]byte("")),
		LogLevel:                        int(gethlog.LvlInfo),
		LogPath:                         log.SysOut,
		UseInMemoryDB:                   true, // todo (config) - persistence should be on by default before production release
		EdgelessDBHost:                  "",
		SqliteDBPath:                    "",
		ProfilerEnabled:                 false,
		MinGasPrice:                     big.NewInt(params.InitialBaseFee),
		TenGenesis:                      "",
		DebugNamespaceEnabled:           false,
		MaxBatchSize:                    1024 * 55,
		MaxRollupSize:                   1024 * 128,
		GasPaymentAddress:               gethcommon.HexToAddress("0xd6C9230053f45F873Cb66D8A02439380a37A4fbF"),
		MinBaseFee:                      new(big.Int).SetUint64(params.InitialBaseFee),

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

	// WETH9 contract pre-deployed at the same address used in production config
	contracts := []tengenesis.Contract{
		{
			Address:  gethcommon.HexToAddress("0x1000000000000000000000000000000000000042"),
			Bytecode: "0x6080604052600436106100bc5760003560e01c8063313ce56711610074578063a9059cbb1161004e578063a9059cbb146102cb578063d0e30db0146100bc578063dd62ed3e14610311576100bc565b8063313ce5671461024b57806370a082311461027657806395d89b41146102b6576100bc565b806318160ddd116100a557806318160ddd146101aa57806323b872dd146101d15780632e1a7d4d14610221576100bc565b806306fdde03146100c6578063095ea7b314610150575b6100c4610359565b005b3480156100d257600080fd5b506100db6103a8565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101155781810151838201526020016100fd565b50505050905090810190601f1680156101425780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561015c57600080fd5b506101966004803603604081101561017357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610454565b604080519115158252519081900360200190f35b3480156101b657600080fd5b506101bf6104c7565b60408051918252519081900360200190f35b3480156101dd57600080fd5b50610196600480360360608110156101f457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135811691602081013590911690604001356104cb565b34801561022d57600080fd5b506100c46004803603602081101561024457600080fd5b503561066b565b34801561025757600080fd5b50610260610700565b6040805160ff9092168252519081900360200190f35b34801561028257600080fd5b506101bf6004803603602081101561029957600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610709565b3480156102c257600080fd5b506100db61071b565b3480156102d757600080fd5b50610196600480360360408110156102ee57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610793565b34801561031d57600080fd5b506101bf6004803603604081101561033457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160200135166107a7565b33600081815260036020908152604091829020805434908101909155825190815291517fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9281900390910190a2565b6000805460408051602060026001851615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f8101849004840282018401909252818152929183018282801561044c5780601f106104215761010080835404028352916020019161044c565b820191906000526020600020905b81548152906001019060200180831161042f57829003601f168201915b505050505081565b33600081815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b4790565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600360205260408120548211156104fd57600080fd5b73ffffffffffffffffffffffffffffffffffffffff84163314801590610573575073ffffffffffffffffffffffffffffffffffffffff841660009081526004602090815260408083203384529091529020547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff14155b156105ed5773ffffffffffffffffffffffffffffffffffffffff841660009081526004602090815260408083203384529091529020548211156105b557600080fd5b73ffffffffffffffffffffffffffffffffffffffff841660009081526004602090815260408083203384529091529020805483900390555b73ffffffffffffffffffffffffffffffffffffffff808516600081815260036020908152604080832080548890039055938716808352918490208054870190558351868152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35060019392505050565b3360009081526003602052604090205481111561068757600080fd5b33600081815260036020526040808220805485900390555183156108fc0291849190818181858888f193505050501580156106c6573d6000803e3d6000fd5b5060408051828152905133917f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65919081900360200190a250565b60025460ff1681565b60036020526000908152604090205481565b60018054604080516020600284861615610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190941693909304601f8101849004840282018401909252818152929183018282801561044c5780601f106104215761010080835404028352916020019161044c565b60006107a03384846104cb565b9392505050565b60046020908152600092835260408084209091529082529020548156fea265627a7a7231582091c18790e0cca5011d2518024840ee00fecc67e11f56fd746f2cf84d5b583e0064736f6c63430005110032",
		},
	}

	gen := &tengenesis.Genesis{
		Accounts:  accts,
		Contracts: contracts,
	}

	genesisBytes, err := json.Marshal(gen)
	if err != nil {
		panic(err)
	}
	return string(genesisBytes)
}
