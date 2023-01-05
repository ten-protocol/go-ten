## Solidity smart contracts

In this subdirectory you can find the solidity smart contracts for the platform.
Under the `/management` subdirectory you can find the root contract of Obscuro - ManagementContract.sol
It dictates the possible state of the Layer 2 and drives the process. 

Under `/messaging` you can find the cross chain messaging contracts. Inside of it under the `/messenger` subdirectory sits an example implementation of cross chain message relayer that utilizes the message bus.

Under `/bridge` you can find the cross chain enabled erc20 standard bridge. 

The contracts interfaces are defined from the perspective of the API consumer contracts and what they should be aware of.

### Compiling

You can use `npx hardhat compile` in order to produce the artifacts for this smart contracts. The results should be generated into
`contracts/artifacts`

### Dependencies

`import "@openzeppelin/...` is enabled using a remapping to the module downloaded by npm under `../node_modules`. This also means that the version of it is maintained in `../package.json`. Notice that the version should be fixed on a concrete release like this `4.5.0` rather than using matchers `^4.5.0` in order to avoid getting unwanted updates.
