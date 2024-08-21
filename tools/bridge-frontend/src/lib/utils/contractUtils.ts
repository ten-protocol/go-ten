import { StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { ethers } from "ethers";

const handleError = (error: any, customMessage: string) => {
  console.error(customMessage, error);
  throw new Error(customMessage);
};

const constructMerkleTree = (leafEntries: any[], msgHash: string) => {
  const tree = StandardMerkleTree.of(leafEntries, ["string", "bytes32"]);
  const proof = tree.getProof(["v", msgHash]);
  return { tree, proof };
};

const estimateGasWithTimeout = async (
  managementContract: ethers.Contract,
  msg: any,
  proof: any,
  root: any,
  timeout = 30000,
  interval = 5000
) => {
  let gasLimit: ethers.BigNumber | null = null;
  const startTime = Date.now();
  while (!gasLimit) {
    try {
      gasLimit = await managementContract.estimateGas.ExtractNativeValue(
        msg,
        proof,
        root,
        {}
      );
    } catch (error: any) {
      console.log(`Estimate gas threw error: ${error.reason}`);
    }
    if (Date.now() - startTime >= timeout) {
      console.log("Timed out waiting for gas estimate, using default");
      return ethers.BigNumber.from(2000000);
    }
    await new Promise((resolve) => setTimeout(resolve, interval));
  }
  return gasLimit;
};

export { handleError, constructMerkleTree, estimateGasWithTimeout };
