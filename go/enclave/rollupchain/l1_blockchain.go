package rollupchain

import (
	"os"
	"path"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/consensus/clique"

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

const (
	gethDir             = "geth"
	chainDataDir        = "chaindata"
	chainDataAncientDir = "chaindata/ancient"
	trieCacheDir        = "triecache"
	ethashDir           = "ethash"
	// TODO - Use a constant that makes sense outside of the simulation.
	dataDirRoot = "../.build/simulations/gethDataDir"
)

// TODO - Add the constants used in this file to the config framework.

// NewL1Blockchain creates a Geth BlockChain object. `genesisJSON` is the Genesis block config in JSON format.
// A Geth node can be made to output this using the `dumpgenesis` startup command.
func NewL1Blockchain(genesisJSON []byte, logger gethlog.Logger) *core.BlockChain {
	dataDir := createDataDir(logger)

	db := createDB(dataDir, logger)
	cacheConfig := createCacheConfig(dataDir)
	chainConfig := createChainConfig(db, genesisJSON, logger)
	engine := createEngine(dataDir, chainConfig, db)
	vmConfig := createVMConfig()
	shouldPreserve := createShouldPreserve()
	txLookupLimit := ethconfig.Defaults.TxLookupLimit // Default.

	blockchain, err := core.NewBlockChain(db, cacheConfig, chainConfig, engine, vmConfig, shouldPreserve, &txLookupLimit)
	if err != nil {
		logger.Crit("l1 blockchain could not be created. ", log.ErrKey, err)
	}
	return blockchain
}

func createDataDir(logger gethlog.Logger) string {
	err := os.MkdirAll(dataDirRoot, 0o700)
	if err != nil {
		logger.Crit("l1 blockchain data directory could not be created. ", log.ErrKey, err)
	}
	dataDir, err := os.MkdirTemp(dataDirRoot, "")
	if err != nil {
		logger.Crit("l1 blockchain data directory could not be created. ", log.ErrKey, err)
	}

	return dataDir
}

func createDB(dataDir string, logger gethlog.Logger) ethdb.Database {
	root := path.Join(dataDir, gethDir, chainDataDir)           // Defaults to `geth/chaindata` in the node's data directory.
	cache := 2048                                               // Default.
	handles := 2048                                             // Default.
	freezer := path.Join(dataDir, gethDir, chainDataAncientDir) // Defaults to `geth/chaindata/ancient` in the node's data directory.
	namespace := ""                                             // Defaults to `eth/db/chaindata`.
	readonly := false                                           // Default.

	db, err := rawdb.NewLevelDBDatabaseWithFreezer(root, cache, handles, freezer, namespace, readonly)
	if err != nil {
		logger.Crit("l1 blockchain database could not be created. ", log.ErrKey, err)
	}
	return db
}

func createCacheConfig(dataDir string) *core.CacheConfig {
	return &core.CacheConfig{
		TrieCleanLimit:      4096 * 15 / 100,                            // Default. 15% of 4096MB allowance for internal caching on mainnet.
		TrieCleanJournal:    path.Join(dataDir, gethDir, trieCacheDir),  // Defaults to `geth/triecache` in the node's data directory.
		TrieCleanRejournal:  ethconfig.Defaults.TrieCleanCacheRejournal, // Default.
		TrieCleanNoPrefetch: false,                                      // Default.
		TrieDirtyLimit:      4096 * 25 / 100,                            // Default. 25% of 4096MB allowance for internal caching on mainnet.
		TrieDirtyDisabled:   false,                                      // Default.
		TrieTimeLimit:       ethconfig.Defaults.TrieTimeout,             // Default.
		SnapshotLimit:       4096 * 10 / 100,                            // Default. 10% of 4096MB allowance for internal caching on mainnet.
		Preimages:           false,                                      // Default.
	}
}

func createChainConfig(db ethdb.Database, genesisJSON []byte, logger gethlog.Logger) *params.ChainConfig {
	genesis := &core.Genesis{}
	err := genesis.UnmarshalJSON(genesisJSON)
	if err != nil {
		logger.Crit("l1 blockchain genesis JSON could not be parsed. ", log.ErrKey, err)
	}

	chainConfig, _, err := core.SetupGenesisBlockWithOverride(
		db,
		genesis,
		nil, // Default.
		nil, // Default.
	)
	if err != nil {
		logger.Crit("l1 blockchain genesis block could not be created. ", log.ErrKey, err)
	}
	return chainConfig
}

// Recreates `eth/ethconfig/config.go/CreateConsensusEngine()`.
func createEngine(dataDir string, chainConfig *params.ChainConfig, db ethdb.Database) consensus.Engine {
	var engine consensus.Engine
	if chainConfig.Clique != nil {
		engine = clique.New(chainConfig.Clique, db)
	} else {
		engine = ethash.New(ethash.Config{
			PowMode:          ethash.ModeNormal,                          // Default.
			CacheDir:         path.Join(dataDir, gethDir, ethashDir),     // Defaults to `geth/ethash` in the node's data directory.
			CachesInMem:      ethconfig.Defaults.Ethash.CachesInMem,      // Default.
			CachesOnDisk:     ethconfig.Defaults.Ethash.CachesOnDisk,     // Default.
			CachesLockMmap:   ethconfig.Defaults.Ethash.CachesLockMmap,   // Default.
			DatasetDir:       "",                                         // Defaults to `~/Library/Ethash` in the node's data directory.
			DatasetsInMem:    ethconfig.Defaults.Ethash.DatasetsInMem,    // Default.
			DatasetsOnDisk:   ethconfig.Defaults.Ethash.DatasetsOnDisk,   // Default.
			DatasetsLockMmap: ethconfig.Defaults.Ethash.DatasetsLockMmap, // Default.
			NotifyFull:       false,                                      // Default.
		}, nil, false) // Defaults.
		interface{}(engine).(*ethash.Ethash).SetThreads(-1) // Disables CPU mining.
	}
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

type blockIngestionType struct {
	// latest is true if this block was the canonical head of the L1 chain at the time it was submitted to enclave
	// (if false then we are behind and catching up, expect to be fed another block immediately afterwards)
	latest bool

	// fork is true if the ingested block is on a different branch to previously known head
	// (resulting in rewinding of one or more blocks that we had previously considered canonical)
	fork bool

	// preGenesis is true if there is no stored L1 head block.
	// (L1 head is only stored when there is an L2 state to associate it with. Soon we will start consuming from the
	// genesis block and then, we should only see one block ingested in a 'preGenesis' state)
	preGenesis bool
}
