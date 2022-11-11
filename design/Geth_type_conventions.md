# Code Conventions Proposal - Geth Type Renames

This proposal seeks to establish a standard on how the enclave codebase should use the standard `geth` types like `types.Transaction`, `types.Block` and the likes.
Currently the code base sparsely uses renamed types like the ones defined in `common/types.go`. The suggestion of this proposal is that such type renames should be the standard
and not the exception. 

## Example and thought process behind the suggestion

Lets take a look at an example function from the code:
```go
// SubmitBlock is used to update the enclave with an additional L1 block.
func (rc *RollupChain) SubmitBlock(block types.Block, isLatest bool) (*common.BlockSubmissionResponse, error) {
```

Notice how the comment explains what type of block should be passed to it. However the same information can be conveyed without having a comment at all:

```go
func (rc *RollupChain) SubmitBlock(block common.L1.Block, isLatest bool) (*common.BlockSubmissionResponse, error) {
```

This function signature makes it immediately obvious that this function should only be called with L1 blocks. While we might now that there aren't any L2 blocks and thus
every `types.Block` is an L1 block, this information is implied by having knowledge of the code base and how Obscuro works. This is not ideal for an open source code base,
which also needs to be highly auditable by security researchers. 

## Reasons

1. For newcomers to the codebase
    * Renamed types like `common.L1.Log`, `common.L1.Receipt`, `common.L2.Receipt` reduce the contextual information required before grasping how concrete functions operate and on what.
    * Makes it easier to safely contribute. Code reviews of PRs matching those conventions should also be easier. 
    * Code becomes more auditable
2. General
    * This is an open source project and this change will make it more inviting for people to look into it. 
    * Generally reduced confusion and easier collaboration. It's better for the code the express the intent of how it should be used, instead of explaining the intent on discord.
    * It will reduce the chance of accidental mistakes. Using the geth types interchangably and operating on them increases the chances that someone will slip in some type coming from the wrong layer into a processor expecting the opposite. This might cause very serious exploits. 
    * Overall future proofing. As the team grows it would be far easier to collaborate with a self explanatory code base.
    * The code base already contains comments that explain where the type is coming from, so making it standard to do so seems reasonable
    * It is a non intrusive change, the geth code will accept our type renames just fine either way.

## The suggest convention

Inside of the enclave code base:

 * **Do** use common type renames that make it explicit where should this type be coming from. 
 * **Do not** use any of the geth types in functions that mutate the enclave/rollup or whatever state.
 * **Exception** using the geth types in common utility functions that perform some abstract operation like `func ShortHash(hash common.Hash) uint64` is **okay**. 

There might be geth types which are completely fine and irrelevant between chains, but the idea is trying to capture types like `Transaction(s)`, `Receipt(s)`, `Log(s)` and the rest "layer-contextual" types.
This might also apply to Addresses, `common.L1.Address`, `common.L2.Address` do imply a specific intent in some instances, but of course apply common sense if it is required to describe such an intent. 