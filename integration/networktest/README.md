# Network test
Network tests provide standard structures for running functional, load and scenario tests against a network.

They are designed to be network agnostic, able to run against both local dev networks and remote testnets.

## Actions
Tests are structured using "Actions" which can be run in series of parallel in a nested tree down to very small steps.

Action's `Run()` method takes a context and a network connector, it returns a context and potential error.

Actions can use the context to pass values to tests that are running later in the flow, so properties that an action require 
from the context are part of that action's API. We use the ActionKey type for the context keys.

## Running tests
All tests can be found under the `/tests` directory. They are grouped into a few packages to try to keep similarly used 
tests together.

### `/helpful`
These are ad-hoc helpers, for example:

- running a local network indefinitely
- smoke test that checks a network isn't completely broken

### `/load`
These tests are designed to apply simulated user activity for a period of time to stress the network.

### `/nodescenario`
These tests require a network that provides access to its nodes through the `NodeOperator` interface.

They can be used to test scenarios such as rejoining the network after a restart, sequencer failover, nodes losing connectivity.

### `/ci`
These are the only tests that will run during the CI builds by default, they should be quick and not fragile.

## UserWallet
In `/userwallet` is a high-level client that bundles a simulated user's private key, an RPC client and manages the nonce and viewing key.

It aims to make testing easier by mimicking the functionality of software and hardware wallets in the real world (a high-level 
interface for interacting with the crypto chain for a user's account(s))