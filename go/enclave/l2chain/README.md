This is where the bulk of the TEN specific logic is.
The entry point is the `L2Chain`.

Ethereum Blocks and Rollups produced by peers are fed into this datastructure, and it decides which is the canonical chain, 
it produces rollups, and is able provide information when requested via RPC. 