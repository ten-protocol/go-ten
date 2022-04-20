package enclave

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

// NewL1Blockchain creates a Geth BlockChain object. `genesisJSON` is the Genesis block config in JSON format. A Geth
// node can be made to output this using the `dumpgenesis` startup command.
func NewL1Blockchain(genesisJSON []byte) *core.BlockChain {
	dataDir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		panic(err)
	}

	db := createDB(dataDir)
	cacheConfig := createCacheConfig(dataDir)
	chainConfig := createChainConfig(db, genesisJSON)
	engine := createEngine(dataDir)
	vmConfig := createVMConfig()
	shouldPreserve := createShouldPreserve()
	txLookupLimit := ethconfig.Defaults.TxLookupLimit // Default.

	blockchain, err := core.NewBlockChain(db, cacheConfig, chainConfig, engine, vmConfig, shouldPreserve, &txLookupLimit)
	if err != nil {
		panic(fmt.Errorf("l1 blockchain could not be created: %w", err))
	}
	return blockchain
}

func createDB(dataDir string) ethdb.Database {
	root := path.Join(dataDir, "geth/chaindata")            // Defaults to `geth/chaindata` in the node's data directory.
	cache := 2048                                           // Default.
	handles := 2048                                         // Default.
	freezer := path.Join(dataDir, "geth/chaindata/ancient") // Defaults to `geth/chaindata/ancient` in the node's data directory.
	namespace := ""                                         // Defaults to `eth/db/chaindata`.
	readonly := false                                       // Default.

	db, err := rawdb.NewLevelDBDatabaseWithFreezer(root, cache, handles, freezer, namespace, readonly)
	if err != nil {
		panic(fmt.Errorf("l1 blockchain database could not be created: %w", err))
	}
	return db
}

func createCacheConfig(dataDir string) *core.CacheConfig {
	return &core.CacheConfig{
		TrieCleanLimit:      4096 * 15 / 100,                            // Default. 15% of 4096MB allowance for internal caching on mainnet.
		TrieCleanJournal:    path.Join(dataDir, "geth/triecache"),       // Defaults to `geth/triecache` in the node's data directory.
		TrieCleanRejournal:  ethconfig.Defaults.TrieCleanCacheRejournal, // Default.
		TrieCleanNoPrefetch: false,                                      // Default.
		TrieDirtyLimit:      4096 * 25 / 100,                            // Default. 25% of 4096MB allowance for internal caching on mainnet.
		TrieDirtyDisabled:   false,                                      // Default.
		TrieTimeLimit:       ethconfig.Defaults.TrieTimeout,             // Default.
		SnapshotLimit:       4096 * 10 / 100,                            // Default. 10% of 4096MB allowance for internal caching on mainnet.
		Preimages:           false,                                      // Default.
	}
}

func createChainConfig(db ethdb.Database, genesisJSON []byte) *params.ChainConfig {
	genesis := &core.Genesis{}
	err := genesis.UnmarshalJSON(genesisJSON)
	if err != nil {
		panic(fmt.Errorf("l1 blockchain genesis JSON could not be parsed: %w", err))
	}

	chainConfig, _, err := core.SetupGenesisBlockWithOverride(
		db,
		genesis,
		nil, // Default.
		nil, // Default.
	)
	if err != nil {
		panic(fmt.Errorf("l1 blockchain genesis block could not be created: %w", err))
	}
	return chainConfig
}

// Recreates the golden path through `eth/ethconfig/config.go/CreateConsensusEngine()`.
func createEngine(dataDir string) consensus.Engine {
	engine := ethash.New(ethash.Config{
		PowMode:          ethash.ModeNormal,                          // Default.
		CacheDir:         path.Join(dataDir, "geth/ethash"),          // Defaults to `geth/ethash` in the node's data directory.
		CachesInMem:      ethconfig.Defaults.Ethash.CachesInMem,      // Default.
		CachesOnDisk:     ethconfig.Defaults.Ethash.CachesOnDisk,     // Default.
		CachesLockMmap:   ethconfig.Defaults.Ethash.CachesLockMmap,   // Default.
		DatasetDir:       "",                                         // Defaults to `~/Library/Ethash` in the node's data directory.
		DatasetsInMem:    ethconfig.Defaults.Ethash.DatasetsInMem,    // Default.
		DatasetsOnDisk:   ethconfig.Defaults.Ethash.DatasetsOnDisk,   // Default.
		DatasetsLockMmap: ethconfig.Defaults.Ethash.DatasetsLockMmap, // Default.
		NotifyFull:       false,                                      // Default.
	}, nil, false) // Defaults.
	interface{}(engine).(*ethash.Ethash).SetThreads(-1)
	return beacon.New(engine)
}

func createVMConfig() vm.Config {
	return vm.Config{
		EnablePreimageRecording: false, // Default.
	}
}

// We indicate that no blocks are authored by local accounts, and thus all blocks are discarded during reorgs.
func createShouldPreserve() func(header *types.Header) bool {
	return func(header *types.Header) bool {
		return false
	}
}
