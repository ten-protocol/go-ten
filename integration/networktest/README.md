# Network test
Network tests provide standard structures for running functional, load and scenario tests against a network.

They are designed to be network agnostic, able to run against both local dev networks and remote testnets.



## Running tests
All tests can be found under the `/tests` directory. They are grouped into a few packages to try to keep similarly used 
tests together.

### `/ci`
These are the only ones that will run during the CI builds by default, they should be quick and not fragile.

### `/helpful`
These are ad-hoc helpers, for example:

- running a local network indefinitely
- smoke test that checks a network isn't completely broken
- function to send funds to a given account on the configured network

### `/load`
These tests are designed to stress the network in various ways and collect metrics to analyse afterwards.

### `/nodescenario`
These tests require a network that provides access to its nodes through the `NodeOperator` interface.

They can be used to test scenarios such as rejoining the network after a restart, sequencer failover, nodes losing connectivity.

## Simulating traffic
A lot of the tests here involve simulating traffic from user accounts and then verifying what happened.

The `/traffic` package provides utilities for orchestrating simulated users and collecting data about their actions to verify.

## UserWallet
In `/userwallet` is a high-level client that bundles a simulated user's private key, an RPC client and manages the nonce and viewing key.

It aims to make testing easier by mimicking the functionality of software and hardware wallets in the real world (a high-level 
interface for interacting with the crypto chain for a user's account(s))