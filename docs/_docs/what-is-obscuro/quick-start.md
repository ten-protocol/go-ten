---
---
# Developer quick start

The only difference between an Ten and an Ethereum (or Arbitrum) dApp is that on Ten you can hide the internal 
state of the contract. 

The most obvious example is that an ERC20 token deployed on Ten will not respond to balance requests unless you are 
the account owner.

In Ten, the internal node database is encrypted, and the contract execution is also encrypted inside the TEE.
The calls to [getStorageAt](https://docs.alchemy.com/reference/eth-getstorageat) are disabled, so all data access 
requests will be performed through view functions which are under the control of the smart contract developer.

Nobody (which includes node operators and the sequencer) can access the internal state of a contract.

**The only thing you have to do when porting a dApp to Ten is to add a check in your view functions comparing 
the `tx.origing` and `msg.sender` against the accounts allowed to access that data.**

The snippet below illustrates this for an [ERC20 token](https://github.com/ten-protocol/sample-applications/blob/main/number-guessing-game/contracts/ERC20.sol#L25). 

```solidity
function balanceOf(address tokenOwner) public view override returns (uint256) {
    require(tx.origin == tokenOwner || msg.sender == tokenOwner, "Only the token owner can see the balance.");
    return balances[tokenOwner];
}
```

_Note that this works because in Ten all calls to view functions are authenticated._ 