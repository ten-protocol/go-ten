### The Service package

This package contains the different service types, currently Sequencer and Validator. It should only contain the driving logic built on top of the different components from the components package.
Those services expect that the state would already be preinitialized when they are created so the responsibility for this for now sits within the enclave. Example of this is resynchronising the stateDB
upon a restart.