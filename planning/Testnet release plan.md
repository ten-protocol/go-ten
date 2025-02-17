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
   a) key derivation (both the TEN network key, and the symmetric rollup encryption keys) 
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
   - create a mechanism based on incentives where anyone should be able to respond to a secret request , and be rewarded for that 
2. Support a larger number of aggregators
   a) gossip
   b) dynamic discovery

3. Start implementing failure scenarios 	
   a) Malicious actors

4. Design and implement new Wallet extension experience
   - with our current approach of having a local encryption/decryption extension, we won't be able to support mobile wallets
   - if we go with a javascript wrapper that developers have to implement, we won't be able to use MetaMask (anywhere)
   - so an important outcome here is not necessarily that we 'design and implement mobile experience', but rather we revisit that whole piece around UX before going into production

5. Optimise and implement the l1 interaction
    - how to ignore failed transactions
    - ensure the l1 blocks are valid	