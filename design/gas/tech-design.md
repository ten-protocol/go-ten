# Technical Design Of Implementing Gas Mechanics

In this document we cover the technical architecture and requirements of the gas mechanics feature. 

## Requirements

1. High Level requirements for the gas mechanics:
   * Gas pricing should accomodate computational cost of the obscuro platform.
   * Gas pricing should be dynamic based on tx demand.
   * Gas pricing should accomodate L1 costs of cross chain messages.
   * Gas pricing should accomodate L1 costs for tx calldata.
   * Gas should be priced in a custom token rather than ETH.
   * There needs to be a safety mechanism that prevents L2 stalling due to gas bugs.

2. Low level requirements for the gas mechanics:
   * Gas should be fundable from the L1 and automatically provisioned to L2 accounts in order to bootstrap ability to pay for transactions
   * The protocol to provision gas should use the cross chain messaging. In other words it should be trustless and decentralised.
   * The protocol should be open to everyone and shouldn't be restricted to the standard Obscuro bridge for example. **This might not be feasible** 
   * The protocol needs to achieve maximum EVM compatibility.

## Scope

1. Gas related smart contracts
2. Cross chain pricing
3. Protocol wiring to fill gas balances of addresses
4. Rewards distribution?


## Assumptions

An L2 DEX will be available on launch. If there are no partners that can deploy on time, then we should take an existing AMM and deploy it and provide the initial liquidity to a single ETH/OBX pool.

## Options for hooking in OBX to be used as gas

There are 3 phases of determining the gas required on Obscuro:
1. Calculate calldata cost when tx is received. If gas provided can't cover the calldata at minimum then we reject the transaction.
2. Compute the transaction, EVM handles the gas pricing using the configured cost per opcode and deducts accordingly.
3. Ensure cross chain messages are paid for during computation inside the contract whenever `publishMessage` is called.

Normally, addresses have balances associated with them. When a transaction is sent, the gas payment is taken from the balance. All of the EVM code is configured for this. Smart contracts, however, are unable to pay for gas and the EVM provides no ways to mutate the gas, only to read some stuff about it. This means that from the perspective of a contract `msg.value` and `address.balance` can be treated differently than the gas. **Due to this we have two potential options to facilitate gas:**
1. `address.balance` is equal to OBX.
   * `msg.value` therefore is also equal to OBX along with anything gas related.
   * Gas funding means that the enclave will increase the balance of an address by the funded amount when the instruction comes in.
   * Balance is reduced beforehand to pay for calldata, and further reduced by the codebase.
2. `address.balance` is equal to ETH.
    * `msg.value` is also equal to ETH.
    * `OBX` is a deployed ERC20.
    * Gas funding in this instance means minting ERC20 equal to the funded amount.
    * `address.balance` for gas purposes becomes a pointer to a storage slot for `address` under `OBX` ERC20
    * Deducting gas means we reduce ERC20 balance (and increase to whomever gas is paid)
    * Bridged over ETH is added to balance by the platform

Method 1 is much simpler and doesn't require some grand intervention to get working. It is not without its problems. The highlights are:
* Is technically EVM equivalent. Currency behind gas, msg.value, address.balance is the same.
* Is logically EVM incompatible. The backing currency is NOT ETH. If a contract has a requirement that it is, it will require porting it over.
* Might create some `unique` issues regarding voting in a DAO when it becomes a thing. Highly depends on the final DAO so no concrete issues, but probability is very high.

Method 2 highlights:
* Harder to implement. Would require investigation to determine if feasible.
* Is mostly logically equivalent. The main difference would be for rewards paid for gas spent; Instead of sending over value one would have to transfer ERC20.
* `OBX` is a contract and not a pure balance. This means any wired in DAO mechanics will work automatically if gas is moved around using the EVM instead of writing to the storage slots directly. 
* Programmable gas - If EVM is used, then this opens up the possibility to extend what gas means and how it interacts entirely on chain. Think of ERC4337, but on chain. "Address `0xAAAA` is submitting a transaction; Lets deduct balance by calling `OBX.deductGas(0xAAAA, amount)`." - What happens underneath might be a simple balance substraction from `0xAAAA` or it might be substracting from someone volunteering to pay, or exchanging another preferred currency for the OBX amount automatically. 


**Note that the two methods will have differing implementations at the point of adding and deducting gas.** Apart from this, everything else remains the same as it makes no difference where the OBX is stored.

## Components

* Gas station smart contract - The contract on the L1 that facilitates gas funding.
* Gas module - the go module that encapsulates the gas logic.
* Pricing module - the go module that analyzes the L1 prices and provides predictions.

## Safety measures

As gas is required to change the state of the chain, if there was suddenly no gas available there would be no way to mutate the chain, at least from the perspective of accounts and what belongs to them. One way to have sudden 'no gas available' would be due to some catastrophic bug. The upgrade process might be capable of getting us out of such a situation, but we can still add a further layer of safety. 

**"Gas Anarchy Mode"** will be a trigger in the management contract controlled by the Obscuro DAO. When blocks are consumed, if gas anarchy event is encountered coming from the management contract, Obscuro will start processing transactions that have their gasPrice equal to 0 (which translates to no cost at all). This would allow for people to unlock their assets and bridge them back regardless of the gas mechanics working. The following rollups will have to be funded somehow, but it will be largely a manual process and negotiation with the community depending on the issue. 


## Design of the L1 to L2 gas funding of accounts

We will have the `ObscuroGasStation` (please recommend a better name, has to sound like a proper oil monopoly) contract on L1.

It has the following interface: 

```solidity
interface ObscuroGasStation {
    function fundAccount(
        address receiver, //L2 address that will receive the OBX
        uint256 amount //amount of OBX
    ) external;
}
```

`fundAccount` will transfer the OBX `amount` from `msg.sender` to the `receiver`. There are two ways to do this:
  1. Using the `ObscuroBridge`. The gas contract will publish an additional cross chain message that it has funded an account.
  2. Making the gas contract act as a `portal` that can burn and mint `OBX` on both L1 and L2 using cross chain messages.

The first one can still be used by other bridges in order to fund accounts, but it would be somewhat weird interaction.
The second makes the gas station contract also be a bridge. Additionally, if hacked it can mint as much OBX as it wants.
`OBX` will also be a wrapped contract on the L2 if the bridge is used.

In both cases the gas station produces its own cross chain messages. The enclave will listen for them and for every encountered such message will 
call `OnGasFunded` on the receiver from the gas module that encapsulates the access to gas balances.

From there gas will be added using the appropriate way - increment the balance of the account under the storage OR create and run a free synthetic transaction that calls the `CrossChainMessenger`, which then calls the bridge/l2 gas station.

## Design of gas consumption

Gas is normally paid to the `coinbase` address. This is the address that has "mined" the block. In our instance this will be the sequencer who computed the result of the transactions. We cannot however pay all of the gas to the `coinbase` address as some of the gas is meant to be paying for L1 costs. Due to this, some of the gas paid will be redirected to another address of a predeployed contract.

The predeployed L2 contract will be called `ObscuroRevenueAgency`. It's functionality is simple - reward the nominated address of whomever published a rollup. For gas consumption, the relevant functionality would be the `depositGas()` function that will synthetically be called free of charge.

When a transaction arrives it will be fed into the `Gas module` which will then ask the `Pricing module` about the L1 gas price to use for calculations.
This gas price will be an adjusted one with priced in risk, but this will be transparent to the gas module. Using the quote, the `Gas module` will calculate the ETH price wanted for the transaction. It will then get a quote for OBX/ETH and determine if the transaction satisfies minimum price. If it does not, it will get rejected without being posted anywhere (as otherwise we'd have to pay for the calldata).

When a batch is being processed, gas being paid is removed from the balance's location per transaction. Then from the gaspool for the transaction we remove the OBX required for calldata and pay it to the `ObscuroRevenueAgency`. The remainder is used for compute costs.

Computation then proceeds normally, any gas left in the pool is refunded.

## Design of Obscuro Revenue Agency

This contract both collects & rewards fees. It encapsulates the gas economy relevant to the L1.
Depending on economy choices rewards are either in `OBX` or `ETH`

Rewards function by consuming cross chain messages from the management contract. In those messages the information included would be the rollup number, nominated address of receiver and gas consumed. When such a message is received by the contract it will distribute the rewards for the appropriate rollup period. 

```solidity 
struct RollupSubmittedMessage {
    address nominatedReceiver; // L2 address that receives rewards for rollup
    uint256 rollupStart; // the number of the first batch in the rollup
    uint256 rollupEnd; // the number of the last batch in the rollup
    uint256 rollupNumber; // Sequential rollup index
}
```


Collecting fees functionality will vary depending on the economic architecture chosen, but on the high level it takes OBX and puts in in the bucket representing this rollup time period. If we were to use spot L2 swaps, this is the function that would be performing them. This method should return/revert if the OBX is insufficient to cover the expected ETH price.  

The contract will have the following interface:

```solidity
interface ObscuroRevenueAgency {
    function claimReward(message Structs.CrossChainMessage) external; // onlyCrossChainSender(managementContract)
    function payGas(
        uint256 amount, // The amount of gas available for this payment
        uint256 wantsETH, // If spot swap is used, then this is the minimum amount of ETH we want in exchange else this should revert. 
    ) external
}
```


## Design of the Pricing module

The pricing module is responsible for giving a predicted gas price. In order to make the prediction it consumes L1 blocks and looks at how the gas price is trending. If we wanted a simple implementation like what Optimism have done, we could simply return the latest gasPrice + 10%. Alternatively we could employ a statistical model and have more advanced pricing.

The following interface should approximately look something like this:

```go
interface PricingTracker {
    SubmitBlock(types.L1Block) // Update based on L1 block.
    GetEstimation(atBlock types.L1Block) uint64 // Retrieve gas price estimation at block.
}
```


## Design of the Gas module

The gas module is responsible for providing gas estimates, adding OBX to accounts and moving OBX from accounts into gas pools. L1 Blocks are analyzed in order to determine who gets gas.

```go
interface GasHandler {
    SubmitBlock(types.L1Block)
   // OnGasFunded(amount *big.Int)
    GetTotalEstimate(types.Transaction) *big.Int
    GetComputationEstimate(types.Transaction) *big.Int
    GetCalldataEstimate(types.Transaction) *big.Int
    CreateGasPool(types.Transaction, ...) *types.ObscuroGasPool
}
```

