```
contracts/
├── src/
│   ├── L1_management_contract/  # set of contracts deployed to Ethereum that manage the state of the TEN network
│   │   ├── lib/
│   │   │   ├── ICrossChain
│   │   │   ├── INetworkEnclaveRegistry
│   │   │   ├── IRollupContract
│   │   │   └── NetworkConfig
│   │   └── impl/
│   │       ├── CrossChain
│   │       ├── NetworkEnclaveRegistry
│   │       └── RollupContract
│   │
│   ├── Cross_chain_messaging/ # the message bus
│   │   ├── common/
│   │   │   ├── IMessageBus
│   │   │   ├── ICrossChainMessenger
│   │   │   ├── CrossChainMessenger
│   │   ├── lib/
│   │   │   ├── CrossChainEnabledObscuro
│   │   ├── L1/
│   │   │   ├── IMerkleTreeMessageBus
│   │   │   ├── MerkleTreeMessageBus
│   │   └── L2/ 
│   │       └── MessageBus
│   │
│   ├── TEN_system/  # contracts deployed automatically on the TEN network and interfaces supported by the TEN platform
│   │   ├── lib/
│   │   │   ├── IContractTransparencyConfig
│   │   │   ├── IFees
│   │   │   ├── IOnBlockEndCallback
│   │   │   ├── IPublicCallback
│   │   │   ├── Logger
│   │   │   └── Transaction
│   │   ├── impl/
│   │   │   ├── Fees
│   │   │   ├── SystemDeployer
│   │   │   ├── TransactionPostProcessor
│   │   │   └── EthereumBridge
│   │   └── utils/
│   │       └── ZenBase
│   │
│   └── reference_bridge/
│       ├── common/
│       ├── L1/
│       └── L2/
│           ├── lib/
│           │   ├── IBridge
│           │   └── ITokenFactory
│           └── impl/
│               └── EthereumBridge

```