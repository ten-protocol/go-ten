# Outstanding Work for Mainnet

## I. Operational Node 
	- configuration that makes sense
	- robust startup scripts 
 	- Enclave HA
 	- HA Host DB
 	- Mempool
 	- More robust P2P (?)
 	- Harden Ethereum Wallet. More robust L1 tx signing (with manual confirmation?)
 	- HA L1 provider 
  - Obscuro RPC - cleanup
  - secret sharing. Review and harden
  - address compression edge case 

## II. Crypto
	- Implement revelation logic
	- implement key derivation for rollup encryption

## III. Upgrade
	- Implement operational upgrade process (described in design doc). Includes L1 component, RPC endpoints, etc
	- Implement key splitting and sharing to N parties (with stake?)

## V. Security
  - L1 validation
  - Add challenges
  - review everything

## VII. Cross chain
  - Sort out finality for xchain messages, separate from the DA rollup.

## VIII. Management contract
  - 

## IX.  Gateway
  - SGX
  - certificate generation
  - tooling to determine cert used for connection
  - smart contract logic + staking + slashing

