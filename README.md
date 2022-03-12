# Obscuro Simulation
This project is a breadth-first pass through the functionality of a full obscuro system.

## Usage
1. Run main.go, and configure parameters inside.

2. Run simulation_test.go 

## High level overview

There are four main actors:
- L1 nodes (Miners). All of which are configured in a L1Network
- L2 nodes (Aggregators), configured in a L2Network. Each L2 node has to be connected to an L1 node via RPC.
- Users (Wallets)
- The Network

The networks handle the communication between nodes, and, to simulate reality, it introduces a random latency. 


### L1

Similar to real life, the miners produce blocks based on an algorithm that simulates proof of work. 
These blocks include all transactions from the mempool.

There are two types of relevant L1 transactions:

- Deposits: Users who want to move money to the L2
- Rollups: Aggregators who won a POBI round, publish a rollup as an L1 tx

Every time a miner decides that there is a new block on the canonical chain it will notify the connected aggregator via RPC.

### L2
On the L2 things are a bit more complicated.

There are users submitting transactions. 
For now these can only be Transfer or Withdrawal transactions.

Aggregators have to process these transactions, and include them in the rollups they create according to the POBI protocol.
 



