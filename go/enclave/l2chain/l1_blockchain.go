package l2chain

import (
	"errors"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ten-protocol/go-ten/go/common/log"

	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	gethDir             = "geth"
	chainDataDir        = "chaindata"
	chainDataAncientDir = "chaindata/ancient"
	trieCacheDir        = "triecache"
	ethashDir           = "ethash"
	// todo (#1471) - use a constant that makes sense outside of the simulation.
	dataDirRoot = "../.build/simulations/gethDataDir"
)

// todo (#1471) - add the constants used in this file to the config framework

// NewL1Blockchain creates a Geth BlockChain object. `genesisJSON` is the Genesis block config in JSON format.
// A Geth node can be made to output this using the `dumpgenesis` startup command.
func NewL1Blockchain(genesisJSON []byte, logger gethlog.Logger) *core.BlockChain {
	dataDir := createDataDir(logger)

	db := createDB(dataDir, logger)
	cacheConfig := createCacheConfig()
	trieDB := createTrie(db, cacheConfig)
	genesis := createGenesis(genesisJSON, logger)
	chainConfig := createChainConfig(db, trieDB, genesis, logger)
	engine, err := createEngine(chainConfig, db)
	if err != nil {
		logger.Crit(err.Error())
	}
	vmConfig := createVMConfig()
	shouldPreserve := createShouldPreserve()
	txLookupLimit := ethconfig.Defaults.TxLookupLimit // Default.

	blockchain, err := core.NewBlockChain(db, cacheConfig, genesis, nil, engine, vmConfig, shouldPreserve, &txLookupLimit)
	if err != nil {
		logger.Crit("l1 blockchain could not be created. ", log.ErrKey, err)
	}
	return blockchain
}

func createTrie(db ethdb.Database, cacheConfig *core.CacheConfig) *trie.Database {
	// Open trie database with provided config
	return trie.NewDatabase(db, &trie.Config{
		Preimages: cacheConfig.Preimages,
	})
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
	root := path.Join(dataDir, gethDir, chainDataDir) // Defaults to `geth/chaindata` in the node's data directory.
	cache := 2048                                     // Default.
	handles := 2048                                   // Default.
	namespace := ""                                   // Defaults to `eth/db/chaindata`.
	readonly := false                                 // Default.

	db, err := rawdb.NewLevelDBDatabase(root, cache, handles, namespace, readonly)
	if err != nil {
		logger.Crit("l1 blockchain database could not be created. ", log.ErrKey, err)
	}
	return db
}

func createCacheConfig() *core.CacheConfig {
	return &core.CacheConfig{
		TrieCleanLimit: 4096 * 15 / 100, // Default. 15% of 4096MB allowance for internal caching on mainnet.
		// TrieCleanJournal:    path.Join(dataDir, gethDir, trieCacheDir),  // Defaults to `geth/triecache` in the node's data directory.
		// TrieCleanRejournal:  ethconfig.Defaults.TrieCleanCacheRejournal, // Default.
		TrieCleanNoPrefetch: false,                          // Default.
		TrieDirtyLimit:      4096 * 25 / 100,                // Default. 25% of 4096MB allowance for internal caching on mainnet.
		TrieDirtyDisabled:   false,                          // Default.
		TrieTimeLimit:       ethconfig.Defaults.TrieTimeout, // Default.
		SnapshotLimit:       4096 * 10 / 100,                // Default. 10% of 4096MB allowance for internal caching on mainnet.
		Preimages:           false,                          // Default.
	}
}

func createChainConfig(db ethdb.Database, triedb *trie.Database, genesis *core.Genesis, logger gethlog.Logger) *params.ChainConfig {
	chainConfig, _, err := core.SetupGenesisBlockWithOverride(
		db,
		triedb,
		genesis,
		nil, // Default.
	)
	if err != nil {
		logger.Crit("l1 blockchain genesis block could not be created. ", log.ErrKey, err)
	}
	return chainConfig
}

// Recreates `eth/ethconfig/config.go/CreateConsensusEngine()`.
func createEngine(chainConfig *params.ChainConfig, db ethdb.Database) (consensus.Engine, error) {
	// If proof-of-authority is requested, set it up
	if chainConfig.Clique != nil {
		return beacon.New(clique.New(chainConfig.Clique, db)), nil
	}
	// If defaulting to proof-of-work, enforce an already merged network since
	// we cannot run PoW algorithms and more, so we cannot even follow a chain
	// not coordinated by a beacon node.
	if !chainConfig.TerminalTotalDifficultyPassed {
		return nil, errors.New("ethash is only supported as a historical component of already merged networks")
	}
	return beacon.New(ethash.NewFaker()), nil
}

func createVMConfig() vm.Config {
	return vm.Config{
		EnablePreimageRecording: false, // Default.
	}
}

func createGenesis(genesisJSON []byte, logger gethlog.Logger) *core.Genesis {
	genesis := &core.Genesis{}
	err := genesis.UnmarshalJSON(genesisJSON)
	if err != nil {
		logger.Crit("l1 blockchain genesis JSON could not be parsed. ", log.ErrKey, err)
	}
	return genesis
}

// We indicate that no blocks are authored by local accounts, and thus all blocks are discarded during reorgs.
func createShouldPreserve() func(header *types.Header) bool {
	return func(header *types.Header) bool {
		return false
	}
}

type BlockIngestionType struct {
	// IsLatest is true if this block was the canonical head of the L1 chain at the time it was submitted to enclave
	// (if false then we are behind and catching up, expect to be fed another block immediately afterwards)
	IsLatest bool

	// Fork is true if the ingested block is on a different branch to previously known head
	// (resulting in rewinding of one or more blocks that we had previously considered canonical)
	Fork bool

	// PreGenesis is true if there is no stored L1 head block.
	// (L1 head is only stored when there is an L2 state to associate it with. Soon we will start consuming from the
	// genesis block and then, we should only see one block ingested in a 'PreGenesis' state)
	PreGenesis bool
}
