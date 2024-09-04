# Block explorer design

## Scope

The design for the TEN block explorers, tools to allow users to make sense of the contents of the TEN chain. 
There will be two tools - a public block explorer that only displays public information, and a private block explorer 
that also displays private information belonging to that user.
 
## Requirements

* Public block explorer
  * The public block explorer is hosted online and can be accessed by anyone via their browser
  * The public block explorer displays the following information:
    * Network statistics
      * Number of rollups
      * (Should) Number of transactions
      * (Should) Average time per rollup
      * (Should) Number of wallet addresses
    * A feed of the latest rollups and the latest transactions
    * (Should) The bytecode and optionally sourcecode (via GitHub) of any deployed contracts
  * Any rollup can be retrieved and displayed, alongside its associated L1 block and (for TestNet only) its decrypted 
    transaction blob. A rollup can be retrieved in any of the following ways:
    * Clicking on a rollup in the list of latest rollups
    * Clicking on a transaction in the list of latest transaction (the rollup containing the transaction is shown)
    * Searching for a rollup by number
    * Searching for a transaction by hash (the rollup containing the transaction is shown)

* Private block explorer
  * The private block explorer is run locally
  * The private block explorer is tied to a specific user, and uses their viewing key to decrypt information in a 
    secure way
  * The private block explorer has the same capabilities as the public block explorer. In addition, it allows the 
    following:
    * Displaying the decrypted transaction and transaction receipt if the user is allowed to view them
    * Searching for transactions by:
      * Address [TODO: CLARIFY WHAT THIS MEANS]
      * Token symbol [TODO: CLARIFY WHAT THIS MEANS]

## Design

We will build our own private and public block explorers.

Design is TBD.

## Known limitations

TBD

## Alternatives considered

### Build our own private and public block explorers

The downside of this approach is that we'd have to write all the block explorer logic from scratch. While displaying 
individual rollups and transactions is relatively simple, the complexity emerges when processing the data to present a 
useful view to the customer. For example, a user may wish to view all the transfers to their address by an ERC-20 
contract, which requires walking the chain and storing the results locally so that they don't have to be continuously 
recomputed.

However, a decisive upside of building our own block explorers is that Ten's rules about data visibility mean that 
an off-the-shelf block explorer is unlikely to be fit for purpose in various ways, and will require extensive 
customisation. We talk about that in section `Fork an existing block explorer for the public block explorer`, below.

By writing our own block explorers, we can also reuse our existing stack (i.e. Go).

### Fork an existing block explorer

The only suitable, open-source block explorer that we are aware of is 
[BlockScout](https://github.com/blockscout/blockscout).

#### Fork BlockScout

BlockScout is an open-source block explorer, used by Secret Network among others (see 
[here](https://explorer.secret.dev/)).

In theory, this would give us a block explorer "for free". In practice, we'd need to customise BlockScout to a large  
extent, even for the public block explorer, for two reasons:

* It cannot handle the fact that some information about the TEN chain is returned in an encrypted form. For 
  example, if vanilla BlockScout is connected to an TEN host, it correctly displays the number of TEN blocks, 
  but it considers every block to have zero transactions, because it chokes on the encrypted transaction contents being 
  returned
* Every advanced block explorer has some customised handling of standard contracts. For example, for ERC-20, it will 
  process the chain to allow a given address to see its entire holdings of various tokens. In Ten, this processing 
  would have to happen inside the enclave, since the block explorer would not have access to the transaction contents. 
  Since BlockScout is not written with this in mind, it would entail a large amount of custom code. This is especially
  true since the logical place to do this sensitive processing is inside the enclave, but BlockScout is not written in 
  Go, and thus it's logic is not reusable

Even further customisations would be required for the private block explorer, where we would have to introduce handling 
of viewing keys.

Forking BlockScout would require us to develop skills we don't have currently (e.g. it is written in Elixir).

Once BlockScout was forked, we'd have to maintain the fork. Blockscout is currently c. 270k lines of code, 20% larger 
than the TEN codebase as of this writing.
