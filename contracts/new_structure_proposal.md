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
│   │       ├── CrossChain  # depends on "Cross_chain_messaging". It is the only input to the cross-chain message bus.
│   │       ├── NetworkEnclaveRegistry
│   │       └── RollupContract # todo - rename to DataAvailabilityRegistry
│   │
│   ├── Cross_chain_messaging/ # the message bus containing both L1 and L2 contracts.
│   │   ├── common/
│   │   │   ├── IMessageBus
│   │   │   ├── ICrossChainMessenger
│   │   │   ├── CrossChainMessenger
│   │   ├── lib/
│   │   │   ├── CrossChainEnabledObscuro # rename to CrossChainEnabledTEN
│   │   ├── L1/
│   │   │   ├── IMerkleTreeMessageBus
│   │   │   ├── MerkleTreeMessageBus
│   │   └── L2/ 
│   │       └── MessageBus
│   │
│   ├── TEN_system/  # "system" contracts deployed automatically on the TEN network and managed by the platform. Also configuration convenstions supported by the TEN platform
│   │   ├── config/
│   │   │   └── IContractTransparencyConfig # if implemented by a SC it will configure custom data "visibility rules"
│   │   ├── lib/
│   │   │   ├── IFees  # todo needs a proper explanation of how it works.
│   │   │   ├── IOnBlockEndCallback # custom decentralised platform logic based on the transactions.
│   │   │   ├── IPublicCallback     # SC can register decoupled actions to be executed 
│   │   │   ├── Logger
│   │   │   └── Transaction
│   │   ├── impl/
│   │   │   ├── Fees 
│   │   │   ├── SystemDeployer
│   │   │   ├── TransactionPostProcessor
│   │   │   └── EthereumBridge
│   │   └── utils/
│   │       └── ZenBase - example of `IOnBlockEndCallback` usage to reward tx activity by automatically issuing "ZEN" tokens. 
│   │
│   └── reference_bridge/ # reference implementation of a simple mint/burn bridge using the `Cross_chain_messaging`.
│       ├── common/
│       ├── L1/
│       │ - same as now
│       └── L2/
│           ├── lib/
│           │   ├── IBridge
│           │   └── ITokenFactory
│           └── impl/
│               └── EthereumBridge

```
