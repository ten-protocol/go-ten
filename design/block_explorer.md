# Block explorer design

## Scope

The design for the Obscuro block explorers, tools to allow users to make sense of the contents of the Obscuro chain. 
There will be two tools - a public block explorer that only displays public information, and a private block explorer 
that also displays private information belonging to that user.

## Requirements

[TODO: HARMONISE INFO SHOWN ACROSS PUBLIC AND PRIVATE EXPLORER?]
[TODO: DO WE WANT TO DISPLAY ROLLUPS?]

* Public block explorer
  * Anyone can access the public block explorer via their browser
  * The public block explorer does not decrypt any information using viewing keys
  * The public block explorer displays the following information:
    * Encrypted transactions in mempool [TODO: IS THIS A NUMBER OR A LIST?]
    * Encrypted transactions validated [TODO: IS THIS A NUMBER OR A LIST?]
    * Transaction metadata:
      * Fees 
      * Timestamps
      * Block number
      * Transaction status
    * (Should) Average number of transactions per rollup
    * (Should) The bytecode and optionally sourcecode (via GitHub) of any deployed contracts
  * [TODO: WHAT DECRYPTION CAPABILITIES FOR THE PUBLIC BLOCK EXPLORER? ANY REQUIREMENT TO EMBED THESE CAPABILITIES IN THE TOOL?]
* Private block explorer
  * The private block explorer is run locally
  * The private block explorer is tied to a specific user, and uses their viewing key to decrypt information in a 
    secure way
  * The private block explorer should display the following information (highest priority first):
    * Latest blocks and their aggregator [TODO: SHOULD THE BLOCKS HAVE THEIR TXS DECRYPTED AUTOMATICALLY? ON DEMAND?]
    * Latest transactions [TODO: CLARIFY THAT THIS IS JUST FOR THE CURRENT USER], with information on [TODO: I ASSUME THIS INFO IS DECRYPTED?]:
      * Status
      * From
      * To
      * State [TODO: CLARIFY WHAT THIS IS]
      * Timestamp
    * Number of transactions
    * The bytecode and optionally sourcecode (via GitHub) of any deployed contracts [TODO: CLARIFY WHETHER THIS IS JUST FOR CURRENT USER]
    * Ability to search for transactions by: [TODO: ADD SOME REQUIREMENT HERE AROUND SPEED AND NOT RESCANNING ENTIRE BLOCKCHAIN - I ASSUME THIS IS WHERE A LOT OF COMPLEXITY COMES IN?]
      * Address
      * Transaction hash
      * Block
      * Token [TODO: CLARIFY WHAT THIS MEANS]
    * Network health (e.g. time to complete gossip)
    * Revelation periods and associated transactions and contracts [TODO: CLARIFY WHAT THIS MEANS]

## Design

TBD - Will depend on which alternative we select from the below.

## Known limitations

TBD - Will depend on which alternative we select from the below.

## Alternatives considered

### Build our own private and public block explorers

TODO

  * DOCUMENT WHAT'S EXPENSIVE FROM AN ENG PERSPECTIVE (DIGESTING CHAIN INTO DB, GETTING CONTRACTS?)
  * TALK ABOUT NEED FOR DESIGN WORK, NO CURRENT CAPABILITY IN-HOUSE
  * TALK ABOUT BENEFIT OF SHARED ARCH + LIBS FOR PUBLIC + PRIVATE EXPLORERS
  * TALK ABOUT BENEFIT OF REUSING EXISTING STACK (GOLANG)

### Fork BlockScout for the public block explorer

[BlockScout](https://github.com/blockscout/blockscout) is an open-source block explorer, used by Secret Network among 
others (see [here](https://explorer.secret.dev/)).

In theory, this would give us a block explorer "for free". In practice, we'd need to customise BlockScout to some 
extent because it cannot handle the fact that some information about the Obscuro chain is returned in an encrypted 
form. For example, vanilla BlockScout correctly displays the number of Obscuro blocks, but it considers every block to 
have zero transactions, because it chokes on the encrypted transaction contents being returned.

Forking BlockScout would require us to develop skills we don't have currently (e.g. it is written in Elixir), and would 
require us to maintain the fork. Blockscout is currently c. 270k lines of code, 20% larger than the Obscuro codebase as 
of this writing.

Meanwhile, we'd also have to develop the private block explorer, so we'd be maintaining two block explorers with no 
common architecture.
