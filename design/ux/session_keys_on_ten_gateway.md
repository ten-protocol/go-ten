# Session keys on Ten

EIP-4337 (Account Abstraction) has popularised the "Session keys" concept to minimize clicks.

We want to offer a similar user-friendly experience to Ten dApp devs.

To test the concept, we will convert the Battleships game to a no-click UX.


## The Ten Gateway as a Session Key manager

*Reminder: A session key (SK) is a key that is not managed in a wallet and can be used to sign behind-the-scenes operations without user action.*

Without smart contract wallets, a SK must have a balance to pay for gas.

The browser can manage SKs on behalf of the user, but it can lose the gas if the browser crashes.


The TG already manages VKs on behalf of users, so adding SKs is relatively straightforward.

The advantage of the TG managing SKs is that it can return the funds to the EOA anytime.


Another advantage is that they can be reused between sessions to avoid unnecessary transactions.


### The flow

The logic of the game UI (javascript) will call an endpoint on the TG that will return a value transfer transaction that the user must sign in their wallet, together with the address of the SK.

The value transfer is required to prepay for the moves. 

*Note: If the user authenticates with a service that is happy to prepay (a paymaster), then the signing step is not required.*

The game will create unsigned move transactions and submit them to the gateway. 
The TG will sign these transactions with the SK and submit them to the network.


## Data ownership

The SKs are a proxy for the EOA, so we want the EOA to read all data belonging to the SKs.
The SKs are tx signers so the platform will treat them as EOAs.

We need to introduce platform level support for SKs.

### The main place will the `externally_owned_account` table:

```
create table if not exists obsdb.externally_owned_account
(
    id      INTEGER AUTO_INCREMENT,
    address binary(20) NOT NULL,
    owner   INTEGER, // this is a new field that will be populated for SKs and will point to the main account
    primary key (id),
    INDEX USING HASH (address)
);
```

Todo : 
 - Query for events:
 - Query for receipts: 
 - T

### RPC endpoint on Node - Register session key

- Receives a session key signed by a VK signed by the EOA.
- Stores it to the `externally_owned_account` table and points to the true EOA
