
### Version Tracking:
* Each contract's version is now tracked in the NetworkConfig contract
* The ContractVersion struct stores the name, version string, and implementation address
* Clear history of upgrades and current versions

### Upgrade Events:
* ContractUpgraded event is emitted for individual upgrades
* BatchUpgradeCompleted event is emitted for batch upgrades
* Events include both old and new versions for tracking changes

### Batch Upgrade Support:
* The recordBatchUpgrade function allows upgrading multiple contracts in one transaction
* Generates a unique hash for each batch upgrade
* Maintains atomicity of upgrades

```
// 1. Deploy new implementation contracts
CrossChain newCrossChainImpl = new CrossChain();
DataAvailabilityRegistry newDARImpl = new DataAvailabilityRegistry();
... 

// 2. Upgrade the proxies
proxyAdmin.upgradeAndCall(
    crossChainProxy,
    address(newCrossChainImpl),
    abi.encodeWithSelector(newCrossChainImpl.initialize.selector, ...)
);

proxyAdmin.upgradeAndCall(
    daRegistryProxy,
    address(newDARImpl),
    abi.encodeWithSelector(newDARImpl.initialize.selector, ...)
);
..

// 3. Record the upgrades in NetworkConfig
string[] memory names = new string[](2);
string[] memory versions = new string[](2);
address[] memory implementations = new address[](2);

names[0] = "CrossChain";
versions[0] = "2.0.0";
implementations[0] = address(newCrossChainImpl);

names[1] = "DataAvailabilityRegistry";
versions[1] = "2.0.0";
implementations[1] = address(newDARImpl);

// This one we could use with Stefans batched safe tx model 
networkConfig.recordBatchUpgrade(names, versions, implementations);
```
# Bridge Upgrade

Given the old bridge isn't recognised as a system contract, we can just deploy the new version and when we upgrade the nodes it will get passed in. We can have some sort of ForkManager which will determine whether to use old bridge logic ie message bus or the new logic ()

```
// fork_manager.go
package crosschain

import (
    "github.com/ethereum/go-ethereum/common"
)

// ForkConfig represents a hard fork configuration
type ForkConfig struct {
    Name        string
    BlockNumber uint64
    IsActive    bool
    // Add any other fork-specific configurations here
    UseNewBridge bool
}

// ForkManager handles fork transitions
type ForkManager struct {
    forks map[string]*ForkConfig
    logger *log.Logger
}

func NewForkManager(logger *log.Logger) *ForkManager {
    return &ForkManager{
        forks: make(map[string]*ForkConfig),
        logger: logger,
    }
}

// RegisterFork registers a new fork configuration
func (fm *ForkManager) RegisterFork(name string, blockNumber uint64, useNewBridge bool) {
    fm.forks[name] = &ForkConfig{
        Name:        name,
        BlockNumber: blockNumber,
        IsActive:    false,
        UseNewBridge: useNewBridge,
    }
}

// ActivateFork activates a fork at the specified block
func (fm *ForkManager) ActivateFork(name string, currentBlock uint64) error {
    fork, exists := fm.forks[name]
    if !exists {
        return fmt.Errorf("fork %s not registered", name)
    }
    if currentBlock < fork.BlockNumber {
        return fmt.Errorf("fork block not reached")
    }
    fork.IsActive = true
    return nil
}

// ShouldUseNewBridge checks if we should use the new bridge for a given block
func (fm *ForkManager) ShouldUseNewBridge(blockNumber uint64) bool {
    for _, fork := range fm.forks {
        if fork.IsActive && blockNumber >= fork.BlockNumber && fork.UseNewBridge {
            return true
        }
    }
    return false
}
```

```
// When processing a new block
func (e *enclaveAdminService) SubmitL1Block(ctx context.Context, blockData *common.ProcessedL1Data) (*common.BlockSubmissionResponse, common.SystemError) {
    // Check if we need to activate any forks
    for name, fork := range n.forkManager.forks {
        if !fork.IsActive && block.NumberU64() >= fork.BlockNumber {
            
            //TODO Check for presence of upgraded event
            err := n.forkManager.ActivateFork(name, block.NumberU64())
            if err != nil {
                n.logger.Error("Failed to activate fork", "name", name, "error", err)
            }
        }
    }

    //TODO how does the L2 fork matter with L1 fork?

    result, rollupMetadata, err := e.ingestL1Block(ctx, blockData)
	if err != nil {
		// only critical errors ie duplicate block or signed rollup error are returned so we can continue processing if non-critical
		return nil, e.rejectBlockErr(ctx, fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

    if result.IsFork() {
		e.logger.Info(fmt.Sprintf("Detected fork at block %s with height %d", blockHeader.Hash(), blockHeader.Number))
	}
    // block processing
}
```

# Claude Suggested Improvements

### Add Upgrade Validation:
Add a function to validate that all required contracts are at compatible versions
This can help prevent upgrades that would break system compatibility

### Add Upgrade Scheduling:
Allow scheduling upgrades for a future block
This gives time for nodes to prepare for the upgrade
Can include a grace period for nodes to update

### Add Upgrade Rollback Support:
Store previous implementations
Allow rolling back to previous versions if needed
Include a time window for rollbacks

### Add Upgrade Documentation:
Include a field for upgrade notes/description
Store changelog information
Link to upgrade documentation
