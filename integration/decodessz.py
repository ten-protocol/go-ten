import os
from eth2spec.phase0.spec import *


def process_genesis_state():
    # Read and deserialize the genesis state from the file
    with open('genesis.ssz', 'rb') as f:
        genesis_state_res = BeaconState.deserialize(f, os.stat('genesis.ssz').st_size)

    return genesis_state_res


genesis_state = process_genesis_state()
genesis_time = genesis_state.genesis_time
print(genesis_state)
