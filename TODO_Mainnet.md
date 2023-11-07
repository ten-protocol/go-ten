# Outstanding Work for Mainnet

## I. Operational Node
  - Configuration that makes sense 
  - Robust startup scripts 
  - Enclave HA 
  - HA Host DB 
  - Mempool 
  - More robust P2P (?)
  - Harden Ethereum Wallet. More robust L1 tx signing (with manual confirmation?)
  - HA L1 provider 
  - Obscuro RPC - cleanup 
  - Secret sharing. Review and harden 
  - Address compression edge case 

## II. Crypto
  - Implement revelation logic 
  - Implement key derivation for rollup encryption

## III. Upgrade
  - Implement operational upgrade process (described in design doc). Includes L1 component, RPC endpoints, etc
  - Implement key splitting and sharing to N parties (with stake?)

## IV. Security
  - L1 validation
  - Add challenges
  - Review everything
  - Prepare docs for external security audit
  - Arrange the audit 
  - Address the outcome of the audit

## V. Cross chain
  - Sort out finality for xchain messages, separate from the DA rollup.

## VI. Management contract
  - Upgradable contracts 
  - Contract signature handling (multisig to deploy new versions ?)

## VII. Gateway
  - SGX 
  - Certificate generation 
  - Tooling to determine cert used for connection 
  - Smart contract logic + staking + slashing

## VIII. UI libraries, clients and tooling, and monitoring
  - Obscuroscan 
  - Obscuro widget 
  - JS libraries 
  - Hardhat plugin

