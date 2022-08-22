# Obscuro Testnet Change Log

## August 2022-08-22
* Tokens / ERC20 contracts
  * The two pre-installed ERC20 contracts deployed are now named 'HOC' and 'POC', replacing the previous tokens of 'JAM' 
    and 'ETH'. Contract addresses remain the same as before respectively. The tokens have restricted `balanceOf` and 
    `allowance` calls such that only the owner of the account can view details which should be private to them. See 
    `go-obscuro\integration\erc20contract\ObsERC20.sol` for more information. 
  * Testnet now supports a faucet to distribute native OBX on request. Previously pre-funding of accounts meant that 
    no native tokens were required to execute transactions on Obscuro - this is now not the case and native tokens 
    must be requested. Allocation of native OBX, along with HOC and POC tokens is currently not supported automatically 
    and a request to Obscuro Labs should be made on the Faucet Discord channel.  
* Wallet Extension 
  * The wallet extension now supports multiple viewing keys through a single running instance. For more information see
    `go-obscuro/design/wallet_extension.md`, specifically `Handling eth_call requests` for the approach to selecting 
    the required viewing key when multiple keys are registered.

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
