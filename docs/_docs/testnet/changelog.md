---
---
# Obscuro Testnet Change Log

## November 2022-11-22 (v0.7)
  * A variety of stability related issues are fixed within this release. 
  * Inclusion of a health endpoint for system status monitoring. 
  * It is now possible to run an Obscuroscan against a locally deployed testnet. For more information see 
    [building and running a local testnet](https://github.com/obscuronet/go-obscuro/blob/main/README.md#building-and-running-a-local-testnet) 
    in the project readme.
  * A list of the relevant PRs addressed in this release is as below;
    * `12a04c40` Checks whether the head rollup is nil (#859)
    * `619d39b4` Clarify that blocks are L1 blocks (#858)
    * `01884de0` Removes endpoint to get L1 height from Obs node (#856)
    * `9b975f3d` eth_getBlockByNumber and eth_getBlockByHash responds based on the batches, not the rollups (#855)
    * `87588e54` Stores batch at the correct point (#854)
    * `f4d37f6e` Remove geth EVM trace logger (#853)
    * `fcc02555` Distribute and store batches (#850)
    * `243f7ef7` Replace panics with logger.Crit in the enclave (#844)
    * `5f97c1a4` Returns errors from DB methods, instead of `found` bools and critical log events (#842)
    * `3dd03cdc` Uses a write lock instead of a read lock (#847)
    * `c039df7d` Locks the subscription list in LogEventManager for threadsafety (#846)
    * `b1bfed47` Gets number of subs in threadsafe way (#845)
    * `76fde61d` Revert "Fix EVM error casting to use pointer variable" (#841)
    * `d225a75c` Adds methods to host DB for batches (#837)
    * `f3d60127` Fix EVM error casting to use pointer variable (#840)
    * `8e21374b` Fix issues with submit block errors (#838)
    * `ddecd719` Fixes concurrency bug in subscription manager (#839)
    * `12f34d46` Create blockprovider and use it for awaitSecret (#813)
    * `9e524e7f` Removes HeaderWithHashes type (#836)
    * `756b7c16` Removes ExtRollup/ExtRollupWithHash split (#835)
    * `17940c7b` Fixing node start out of sync (#832)
    * `81b8d9c8` Testnet DNS now point to node1 (#827)
    * `de9dbc6f` Cleans up the GetLatestTransactions API method (#833)
    * `6932e020` Fixes grabbing a rollup via ObscuroScan (#829)
    * `c9e978f0` Adds booleans to DB methods to indicate whether was found. (#831)
    * `849ea7aa` Fetch latest Rollup Head now returns error (#826)
    * `bc652690` Adding health check endpoint (#825)
    * `8f049ff9` Handle all errors from ethcall and estimate gas (#823)
    * `9107d571` Fixes eth call error propagation (#822)
    * `976d872c` Remove unused test APIs. Rename RPC method constants for clarity (#821)
    * `3a6f197f` Stop in-mem nodes properly. Prune unused in-mem RPC methods (#820)
    * `84e7c615` Provides logger for Obscuroscan (#819)
    * `088d8f50` Dynamic estimate gas (#815)
    * `4478ffbd` Fix the bridge address to pass the checksums (#812)
    * `ef0e04d9` Downgrade the spammy log message (#810)
    * `6cb0d85a` Have sims test the eth_blockNumber endpoint (#809)
    * `d83c201e` Confusing description of `DB.writeRollupNumber`. Minor clean-up (#791)
    * `5faab414` Fix to use the dev build of the contract deployer (#807)


## November 2022-11-08 (v0.6)
  * The Number Guessing Game has been removed from static and auto deployment scripts, and is now hosted 
    [in a sample applications repository](https://github.com/obscuronet/sample-applications). Given the move for 
    Testnet to be long-running (or at least restartable without contract disappearance), the Guessing Game must be 
    persisted across software updates, and redeployed manually if needed in the same way other applications are.
  * The list of sensitive RPC API methods, where the request and response is encrypted in transit, now covers 
    `eth_call`, `eth_estimateGas`, `eth_getBalance`, `eth_getLogs`, `eth_getTransactionByHash`, `eth_getTransactionCount`, 
    `eth_getTransactionReceipt` and `eth_sendRawTransaction`. See the Obscuro
    [documentation](https://docs.obscu.ro/api/sensitive-apis/) for more details. 
  * Calls to wait for a transaction receipt are now blocking, whereas previously they would return an error meaning the
    client side code needed to perform a specific wait and poll loop. The example on how to [programmatically deploy
    a contract](https://docs.obscu.ro/testnet/deploying-a-smart-contract-programmatically/) has been updated accordingly.
  * The ability to start a faucet server against a local testnet deployment is now supported via a docker 
    container. For more information see the Obscuro 
    [readme](https://github.com/obscuronet/go-obscuro#building-and-running-a-local-faucet).
  * Updates to the [Events](https://github.com/obscuronet/go-obscuro/blob/main/design/Events_design.md) design 
    inclusion of the [Fast Finality](https://github.com/obscuronet/go-obscuro/blob/main/design/fast_finality.md) design.
  * The [Obscuro docs site](https://docs.obscu.ro/) is now searchable. 
  * Testnet is now officially termed `Evan's Cat`.

* ObscuroScan:
  * ObscuroScan supports a single API at [/rollup/](http://testnet.obscuroscan.io/rollup/) which allows web clients to 
    access a JSON representation of rollups and encrypted transactions. Further details 
    [here](https://docs.obscu.ro/testnet/obscuroscan)

## October 2022-10-21 (v0.5)
* Event Subscriptions:
  * Event subscriptions for logs are now supported via the eth_subscribe and eth_getLogs approaches. This has been 
    tested using both the ethers and web3js libraries. Note that eth_newFilter is not currently supported. For more 
    information see [the events design](https://github.com/obscuronet/go-obscuro/blob/main/design/Events_design.md).

## September 2022-09-22 (v0.4)
* Wallet extension:
  * The wallet extension now supports separate ports for HTTP and WebSocket connections. Use the `--port` and `--portWS` 
    command line options respectively for each. For more information see the
    [Wallet extension](https://docs.obscu.ro/wallet-extension/wallet-extension) documentation. 
* Event subscription:
  * An early preview of event subscriptions is available in this release, though note that this is still undergoing 
    testing and feature enhancements and therefore is liable to issues and instability. For more information on the 
    functionality available reach out to the development team on the discord 
    [active testnet developers](https://discord.com/channels/916052669955727371/1004752710077259838) channel. 
* Transaction receipts:
  * Only return receipts for transactions which were included in a canonical rollup.

## September 2022-09-07 (v0.3)
* Tokens / ERC20 contracts
  * The ERC20 'HOC' and 'POC' tokens are now funded with 18 decimal places of precision. Previously funding of 50 
    tokens was erroneously made as 50 10^-18. This means tokens imported into Metamask will display correctly. Note that
    the number guessing game pay to play still costs 1 10^-18 HOC tokens.
* Wallet extension:
  * Viewing keys are now persisted across wallet extension restarts
  * Enhanced logging for registering of viewing keys

## August 2022-08-22 (v0.2)
* Account balances:
  * Added correct calculation of account balances (previously, all accounts were allocated infinite funds).
* Tokens / ERC20 contracts
  * The two pre-installed ERC20 contracts deployed are now named 'HOC' and 'POC', replacing the previous tokens of 'JAM' 
    and 'ETH'. Contract addresses remain the same as before respectively. The tokens have restricted `balanceOf` and 
    `allowance` calls such that only the owner of the account can view details which should be private to them. See 
    `go-obscuro\integration\erc20contract\ObsERC20.sol` for more information. 
  * Testnet now supports a faucet to distribute native OBX on request. Previously pre-funding of accounts meant that 
    no native tokens were required to execute transactions on Obscuro - this is now not the case and native tokens 
    must be requested. Allocation of native OBX, along with HOC and POC tokens is currently not supported automatically 
    and a request to Obscuro Labs should be made on the Faucet Discord channel.
* Gas prices:
  * The node operator can configure the minimum gas price their aggregator will accept on startup.
* Wallet extension 
  * The wallet extension now supports multiple viewing keys through a single running instance. For more information see 
    the [wallet extension design doc](https://github.com/obscuronet/go-obscuro/blob/main/design/wallet_extension.md), 
    specifically `Handling eth_call requests` for how the required viewing key is selected when multiple keys are 
    registered.

## August 2022
* Testnet launch:
  * Testnet preview launched to limited number of application developers.
  * ObscuroScan block explorer for Testnet launched.
  * Number Guessing Game smart contract deployed to Testnet.
* Obscuro Docsite launched.
* Account balances:
  * Added correct calculation of account balances (previously, all accounts were allocated infinite funds).
  * Introduced network faucet account.
  * Obscuro enclaves services can configure the minimum gas price they'll accept
* ``block.difficulty`` will return a true random number generated inside the secure enclave.
