# Obscuro Gas Mechanics economics

## Background

Layer 2's main source of income traditionally is gas profits. Every L2 transaction needs to be submitted in a rollup
on the L1. In the EVM we have storage, memory and calldata. In the same order they are most expensive to least expensive. The optimal way to facilitate data availability while reducing costs has been to encode the transactions as calldata which is never accessed. This means that every transaction has a well known gas cost for each byte it consists of priced in terms of calldata. 

The problem however is that while we know how much gas each transaction will cost beforehand, we do not know what the price of gas would be when the rollup is submitted. This results in all the layer 2's having extended gas mechanics and calculation formulas that factor in price movements in order to minimize rollups that end up at a loss. 

With Obscuro we have another layer on top of the traditional gas price problem - We also want gas to be paid in OBX on the layer 2 rather than ETH. Other L2's all function by using ETH for gas. As traditionally we can't pay with OBX on the L1 for gas we will have to convert the OBX to ETH and then pay. This means that not only we are at risk of gas prices spiking, but the trade price of OBX/ETH going against us. ETH can jump across the board or OBX can fall across the board, or even liquidity pools we depend on might not be as liquid as we need them at the time of the transaction and cause slippage. 

The problem is further compounded by MEV. Using L1 DEXes with the protocol exposes us to profiteers. While slippage can be controlled, we'd still be losing gas for the transaction which will make us suffer further. Worst case scenario we'd have to bail on a DEX in the same transaction that is publishing a rollup. 

And finally, using OBX for gas means that while we are EVM compatible, addresses paying with native currency and the balance of each account would carry different meaning as opposed to the L1. `address.balance` will be the OBX balance on obscuro instead of ETH as it is on mainnet. This isn't an economic problem and won't be discussed further in this paper, but is nonetheless important for the conclusion - Using ETH to pay for gas on L2 is much easier and carries less risk. Solving all the problems is very hard and competitive layer 2's have bailed trying to for the time being. 

## Isolated Addons

There are ways to reduce the risks and improve the profits of the protocol. I'm calling them "Isolated", because in isolation they solve/improve something, but as anything economic there is still no way to guarantee combining them will yield success.

### Alternatives to L1 DEX

When discussing L1 DEXes and how they expose us to MEV an obvious question comes to mind - `"Can't we use a DEX on the L2 and swap OBX for ETH there?"` In theory, yes! Technical details aside, having an automatic swap of the pending OBX to ETH and paying it to the address that published the rollup would remove the MEV exposure. This approach introduces other issues: 

- The liquidity on layer 2's tends to be far worse than L1.
- Price on the L2 is decoupled from the L1 price. 
- Relying on the DEX means that one must be deployed. To deploy one, gas needs to be paid. But paying gas requires a DEX. Then liquidity needs to go in, same issue. Circular dependency problem. Leaving details out for later, solving this would essentially boil down to a launch ceremony where gas gets enabled later on. 
- Publisher needs to have funds on L1 to pay in full, before being able to receive payment from the L2 fees.

However, there is one major benefit when using L2 DEX for gas mechanics - The paid OBX for the gas fee can be converted to ETH on the spot implicitly for every transaction. This is without any impacts on L1 gas costs or L2 gas costs. The implicit computation can be configured to either be fully free or have a fixed OBX gas cost for every transaction regardless of its size. Note that we only need to convert enough to afford the calldata cost on L1. The rest can stay in OBX.

This spot conversion allows us to lock in the ETH required to pay for the transaction on the L1 on the spot, which greatly reduces the risk exposure to pair movements on the L1. The more time passed, the greater the chance of volatility between two assets. When using a layer 1 DEX, the timeframe between getting the fee and converting it is significantly bigger so the risk is way bigger. Arguably this benefit compensates for any additional overheads using L2 DEX brings.

The main challenge with spot conversions is the same challenge every layer 2 has with gas - getting reliable current L1 gas price and adding a safety buffer to it in case of volatility. 


#### The problem of L2 illiquidity in AMM pools

Most projects that extend into L2 tend to share the same property of having significantly less TVL. This is true for both Arbitrum and Optimism and is not just about small projects. Even the biggest projects tend to suffer from lack of liquidity outside the "core pools". Core for example being USDC/ETH. 


In our instance, given that this would be liquidity for OBX/ETH on the platform where OBX is used it might not be such an issue. One could say this is actually a "core pool". Furthermore good fee incentives on the pool along with having guaranteed platform induced demand might make providing liquidity on the L2 even more lucrative than on the L1. The guaranteed demand is due to the pool being core part of the gas protocol and thus everyone being forced to use it.

#### The problem of OBX/ETH pairs being decoupled

Arbitrage will ensure the pairs do not slip for extended periods of time, but mandatory fees will dictate the minimum possible slippage. As the protocol applies pressure on the L2 and makes ETH more expensive in terms of OBX, its reasonable to expect ETH will tend to be more expensive on the L2 most of the time. If the total arbitrage fees sum up to 1% than the price diff between L1 and L2 pairs will only get fixed once it goes above 1% which virtually guarantees near constant cost overhead.

Minimizing arbitrage slippage would also depend on automated players getting in. Relying on the odd trade to balance out will most certainly not work. Being EVM compatible should mean that given a well known DEX interface, arbitrage software that is already used in other places should work with us. There might be minor work required to accomodate the bridge. We should still assume the worst and prepare to subsidize gas costs for months, as price slip might mean fees are extra expensive with us. Spot conversions protect us from losses, but users will bear the cost which ultimately is bad for us.

Another angle to the arbitrage problem is withdrawal lock periods. Long periods to get out the OBX might discourage arbitrage. Liquidity bridge protocols that provide very quick withdrawals might help, but they still would need to deploy on us. 

Ultimately arbitrage might be something we have to do on our own initially, even at some very primitive level. It might be easier to pick up some open source basic (but working) bot and wire in the logic for our bridges.


#### The problem of publisher needing L1 funds beforehand

This is not really a serious problem. Even in the most basic scenario a good amount of funds can be kept on L1 allowing for payments to come in later. Note that this problem is there regardless if the swap is done before rollup submission or after. The rollup being successfully submitted is a prerequisite to take the payment out of the bridge. 

This is also a problem when swapping on L1 DEX too as the contract can't pay the transaction's gas. 

If we wanted to solve this and enable scenarios where address doesn't keep a lot of funds on hand we can use ERC4337.
It allows us to circumvent having to pay for gas ourselves; Someone else will boot the cost of the transaction, but of course this means that the contract would have to pay his address the gas cost along with profit which ends up being additional overhead.


### Prebuying module

While using an L2 DEX will be great in the long run, the "bootstrap" period would be less than ideal due to liquidity and lack of automated arbitrage.
We can add an additional part to the protocol that allows prebuying ETH for OBX. This ETH will exclusively be moved around using the cross chain messaging.
The mechanism will be rather simple - when the implicit L2 swap is being performed, the price will be compared against the fixed offering of the prebuy buffer.
If the buffer is cheaper, OBX will be paid into its L2 contract. Before rollup, the l2 prebuy contract will be called in order to create a cross chain message.
This cross chain message would transfer the OBX to the L1 and put it in the prebuy contract, which once again using the full amount of OBX it had initially will purchase ethereum at L1 spot price. 

If the buffer is empty and cannot cover the OBX swap, then the L2 dex is used as a fallback.

The L2 buffer is filled by the L1 prebuy contract everytime a rollup is submitted and OBX is exchanged. Upon swapping for ETH it is immediately bridged back to the L2. This means that the block that marks the rollup as successful also provides fresh ETH.

The L2 prebuy contract is not accessible by anything else but the enclave. It cannot be used to get better prices than what is available by external parties.


#### Balancing the buffer

If the OBX price drops, then even if the same amount of OBX is perpetually maintained by the prebuy contract, the resulting ETH will be getting less and less. In this scenario we can manually add more OBX, perhaps taken from the sequencer's profits in order to extend the available ETH. In reality the temporary nature of this mechanism means that once the L2 OBX/ETH pool becomes solid there is little need to fill up this prebuy mechanism. 

### Hedge modules

As mentioned in the whitepaper, the currency exposure can be hedged with futures or options IF such a market existed.
With spot conversions there is little sense to hedge. The only real thing to hedge against is the gas price volatility.
If we do not use spot conversions however, hedging becomes important.

One of the most modern option approaches is [Opyn's squeeth](https://www.opyn.co/?ct=NL). The market there is USDC/ETH and it is not exactly a traditonal option contract. However it does provide the same capabilities in a perpetual way. Holding squeeth is essentially a perpetual contract that overtime loses value, the premium paid. Volatility is reflected in the premium and squeeth minting prices.

We can use the current deployed squeeth to hedge against ethereum price going up suddenly in between rollups. This can be done by creating a buffer that is filled with Squeeth and aimed at maintaining some cap of ETH value represented in Squeeth. When ETH shoots off while OBX stays at same dollar value, what the swap lacks to cover the gas cost can be withdrawn from the buffer. When the OBX swap turns a profit the buffer is refilled. 

This however does not help in the case where ETH remains same price, but OBX is dropping. 

An alternative approach would be to deploy our own OBX Squeeth that allows us to cover both OBX dropping and ETH shooting up. This would have a questionable market given the size of ETH/USDC. 

### Selling calldata options

If we ignore OBX and think in terms of L1 ETH pricing of L2 transactions, we have something non obvious - We are selling an option contract for X amount of ETH to deliver of Y amount of bytes. The price is current spot gas price + gas price premium in ETH. This contract is for a fixed amount of time that will always be executed. (I guess this is more like a future, but I have no clue how volatility is priced there) 

So users pay ETH for a contract that you will pay the gas on the L1 for their transaction. The implication of this is that if we were to go on the extreme offensive, we could in theory use established option pricing to provide extremely fair prices. Instead of having some % markup. And this would make us ... option sellers. [Pricing model?](https://www.investopedia.com/terms/b/blackscholes.asp)

Furthermore if we remember OBX, we can stack the optionality: 

1. Get transaction option contract price in ETH.
2. Get OBX price for ETH option contract. (or future, idk)

Now this isn't entirely risk free, but in theory it should allow us to operate without instant swaps or any other party tricks. Given a long enough amount of time and not using too much capital, we should easily be profitable. As pretty much everyone who sells theta is. **That said, I don't think this is the way to go for the time being.** When Obscuro and OBX mature and become somewhat stable we could transition to this after proper back testing.

