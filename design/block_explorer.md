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
    * Encrypted transactions in mempool
    * Encrypted transactions validated
    * Fees, timestamps and other metadata [TODO: CLARIFY WHAT OTHER METADATA]
    * (Should) Average number of transactions per rollup
    * (Should) Any contracts that have been deployed [TODO: CLARIFY WHAT SHOULD BE SHOWN]
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
    * Any contracts that have been deployed [TODO: CLARIFY WHAT SHOULD BE SHOWN, AND WHETHER THIS IS JUST FOR CURRENT USER]
    * Ability to search for transactions by: [TODO: ADD SOME REQUIREMENT HERE AROUND SPEED AND NOT RESCANNING ENTIRE BLOCKCHAIN - I ASSUME THIS IS WHERE A LOT OF COMPLEXITY COMES IN?]
      * Address
      * Transaction hash
      * Block
      * Token [TODO: CLARIFY WHAT THIS MEANS]
    * Network health (e.g. time to complete gossip)
    * Revelation periods and associated transactions and contracts [TODO: CLARIFY WHAT THIS MEANS]

## Design

TBD

## Known limitations

TBD

## Alternatives considered

TBD
