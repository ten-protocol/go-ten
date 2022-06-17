This package dubplicates the geth "rawdb" package.
It contains logic to wrap access to the key value store.
The only changes are around the used prefixes, and the removal of the "ancients" which we don't use for now. 

Note 1: We had to duplicate the geth code, since we're storing both rollup and block information, and so we need different convetions.
Note 2: This needs to be reviewed, and maybe we can find a way to use the geth methods directly. (maybe by having two instances of the kv store, one for blocks and one for rollups)