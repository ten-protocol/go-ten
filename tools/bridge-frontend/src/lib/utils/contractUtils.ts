import { showToast } from "@/src/components/ui/use-toast";
import { ToastType } from "@/src/types";
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { ethers } from "ethers";

const handleError = (error: any, customMessage: string) => {
  console.error(customMessage, error);
  if (error?.message?.includes("unknown account")) {
    throw new Error(
      "Ensure your wallet is unlocked and an account is connected"
    );
  }
  return error;
};

const constructMerkleTree = (leafEntries: any[], msgHash: string) => {
  showToast(ToastType.INFO, "Constructing Merkle tree");
  const tree = StandardMerkleTree.of(leafEntries, ["string", "bytes32"]);
  const proof = tree.getProof(["v", msgHash]);

  return { tree, proof };
};

const estimateGas = async (
  receiver: string,
  value: string,
  bridgeContract: ethers.Contract
) => {
  try {
    if (!value || isNaN(Number(value))) {
      throw new Error("Invalid value provided for gas estimation.");
    }

    const parsedValue = ethers.utils.parseEther(value);
    showToast(ToastType.INFO, "Estimating gas for the transaction");
    return await bridgeContract?.estimateGas.sendNative(receiver, {
      value: parsedValue,
    });
  } catch (error) {
    return handleError(error, "Error estimating gas");
  }
};

const estimateAndPopulateTx = async (
  receiver: string,
  value: string,
  gasPrice: any,
  bridgeContract: ethers.Contract
) => {
  try {
    showToast(ToastType.INFO, "Estimating gas for the transaction");
    const estimatedGas = await estimateGas(receiver, value, bridgeContract);

    showToast(ToastType.INFO, "Populating transaction");
    const populatTxResp = await bridgeContract?.populateTransaction.sendNative(
      receiver,
      {
        value: ethers.utils.parseEther(value),
        gasPrice,
        gasLimit: estimatedGas,
      }
    );
    console.log("🚀 ~ useContractService ~ populatTxResp:", populatTxResp);
    return populatTxResp;
  } catch (error) {
    return handleError(error, "Error populating transaction");
  }
};

const extractAndProcessValueTransfer = async (
  txReceipt: any,
  messageBusContract: ethers.Contract,
  provider: ethers.providers.JsonRpcProvider
) => {
  try {
    showToast(
      ToastType.INFO,
      "Extracting and processing value transfer event data"
    );

    const valueTransferEvent = txReceipt.logs.find(
      (log: any) =>
        log.topics[0] ===
        ethers.utils.id("ValueTransfer(address,address,uint256,uint64)")
    );
    console.log("🚀 ~ useContract ~ valueTransferEvent:", valueTransferEvent);

    if (!valueTransferEvent) {
      throw new Error("ValueTransfer event not found in the logs");
    }

    showToast(ToastType.INFO, "ValueTransfer event found; parsing event data");

    const valueTransferEventData =
      messageBusContract?.interface.parseLog(valueTransferEvent);

    if (!valueTransferEventData) {
      throw new Error("ValueTransfer event data not found");
    }

    showToast(
      ToastType.INFO,
      "ValueTransfer event data found; fetching block data"
    );

    const block = await provider.send("eth_getBlockByHash", [
      ethers.utils.hexValue(txReceipt.blockHash),
      true,
    ]);

    if (!block) {
      throw new Error("Block not found");
    }

    showToast(ToastType.INFO, "Block data found; processing value transfer");

    return { valueTransferEventData, block };
  } catch (error) {
    return handleError(error, "Error processing value transfer");
  }
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
  showToast(ToastType.INFO, "Estimating gas for the transaction");
  while (!gasLimit) {
    try {
      gasLimit = await managementContract.estimateGas.ExtractNativeValue(
        msg,
        proof,
        root,
        {}
      );
      console.log("🚀 ~ gasLimit:", gasLimit);
    } catch (error: any) {
      console.error(`Estimate gas threw error: ${error.reason}`);
    }
    if (Date.now() - startTime >= timeout) {
      console.log("Timed out waiting for gas estimate, using default");
      showToast(
        ToastType.INFO,
        "Timed out waiting for gas estimate; using default"
      );
      return ethers.BigNumber.from(2000000);
    }
    await new Promise((resolve) => setTimeout(resolve, interval));
  }
  return gasLimit;
};

export {
  handleError,
  constructMerkleTree,
  estimateGasWithTimeout,
  extractAndProcessValueTransfer,
  estimateAndPopulateTx,
  estimateGas,
};
