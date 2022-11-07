# Obscuro Testnet Change Log

## Novemeber 2022-11-08 (v0.6)
  * The Number Guessing Game has been removed from static and auto deployment scripts, and is now hosted 
    [in a sample applications repository](https://github.com/obscuronet/sample-applications). Given the move for 
    Testnet to be long-running (or at least restartable without contract disappearance), the Guessing Game must be 
    persisted across software updates, and redeployed manually if needed in the same way other applications are.
  * Calls to wait for a transaction receipt are now blocking, whereas previous they would return an error meaning the 
    client side code needed to perform a specific wait and poll loop. The example on how to programmatically deploy
    a contract has been updated accordingly. 

## October 2022-10-21 (v0.5)
* Event Subscriptions:
  * Event subscriptions for logs are now supported via the eth_subscribe and eth_getLogs approaches. This has been 
    tested using both the ethers and web3js libraries. Note that eth_newFilter is not currently supported. For more 
    information see [the events design](https://github.com/obscuronet/go-obscuro/blob/main/design/Events_design.md).

## September 2022-09-22 (v0.4)
* Wallet extension:
  * The wallet extension now supports separate ports for HTTP and WebSocket connections. Use the `--port` and `--portWS` 
    command line options respectively for each. For more information see the
    [Wallet extension](https://docs.obscu.ro/wallet-extension/wallet-extension.html) documentation. 
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
