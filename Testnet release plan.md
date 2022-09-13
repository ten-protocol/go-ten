# Development milestones

This is the large task breakdown I propose we add to the next development milestones.
There will be bug fixes and UX improvements as we discover them, which will be treated as high priority.

Note that there is no protocol work listed here. These are the tasks that are orthogonal to the protocol.


## V2
1. Events subscription
2. Connect to a well-known Ethereum testnet
   a) support restarting of nodes
3. Improvements to Management contract
   a) Support withdrawals
4. Mempool improvements
   a) Remove failed transactions


## V3
1. Support a limited number of user operated nodes
   a) Catching up robustness
   b) Script robustness
   c) Instructions
   d) Try outside of Azure?
2. Design and implement reveleation period
3. Finalise cryptography
   a) key derivation (both the obscuro network key, and the symmetric rollup encryption keys) 
   b) Calculate hashes on startup
   c) Check that the rollup producer is attested
   d) Add entropy to txs
   e) Encrypt some errors returned from enclave
4. Design and implement the final mempool
5. Add some metrics
    - gauge


## V4
1. Design and implement dynamic Bridge
2. Design and implement tx fees + aggregator rewards


## V5
1. Optimise sharing of the secret
2. Support a larger number of aggregators
   a) gossip
   b) dynamic discovery

3. Start implementing failure scenarios 	
   a) Malicious actors

4. Design and implement new Wallet extension experience

5. Optimise and implement the l1 interaction
    - how to ignore failed transactions
    - ensure the l1 blocks are valid	