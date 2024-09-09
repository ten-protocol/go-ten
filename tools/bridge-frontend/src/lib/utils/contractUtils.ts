import { showToast } from "@/src/components/ui/use-toast";
import { ToastType, TransactionStep } from "@/src/types";
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { ethers } from "ethers";
import { handleStorage } from "../utils";
import {
  addPendingBridgeTransaction,
  removePendingBridgeTransaction,
  updatePendingBridgeTransaction,
} from "./txnUtils";
import { handleError } from "./walletUtils";

const sendTransactionStep = async (
  bridgeContract: ethers.Contract,
  signer: ethers.Signer,
  receiver: string,
  value: string
) => {
  const gasPrice = await signer.provider?.getGasPrice();

  const tx = await estimateAndPopulateTx(
    receiver,
    value,
    gasPrice as ethers.BigNumber,
    bridgeContract
  );

  const txResponse = await signer.sendTransaction(
    tx as ethers.providers.TransactionRequest
  );

  addPendingBridgeTransaction({
    txHash: txResponse.hash,
    resumeStep: TransactionStep.TransactionSubmission,
    receiver,
    value,
    timestamp: Date.now(),
  });

  showToast(ToastType.INFO, "Transaction sent; waiting for confirmation");
  return txResponse;
};

const confirmTransactionStep = async (
  txResponse: ethers.providers.TransactionResponse
): Promise<ethers.providers.TransactionReceipt> => {
  const txReceipt = await getTransactionReceipt(txResponse);
  return txReceipt as ethers.providers.TransactionReceipt;
};

const extractEventDataStep = async (
  messageBusContract: ethers.Contract,
  provider: ethers.providers.JsonRpcProvider,
  txReceipt: ethers.providers.TransactionReceipt
) => {
  const { valueTransferEventData, block } =
    (await extractAndProcessValueTransfer(
      txReceipt,
      messageBusContract,
      provider
    )) || {
      valueTransferEventData: null,
      block: null,
    };

  if (!valueTransferEventData || !block) {
    throw new Error("Failed to extract value transfer event data");
  }

  // update pending txn after event extraction
  updatePendingBridgeTransaction(txReceipt.transactionHash, {
    resumeStep: TransactionStep.MerkleTreeConstruction,
    valueTransferEventData,
    block,
  });

  return { valueTransferEventData, block };
};

const constructMerkleTreeStep = async (
  txHash: string,
  valueTransferEventData: ethers.utils.LogDescription,
  block: any
) => {
  const { tree, proof } = constructMerkleTree(
    JSON.parse(atob(block?.crossChainTree)),
    ethers.utils.keccak256(
      new ethers.utils.AbiCoder().encode(
        ["address", "address", "uint256", "uint64"],
        [
          valueTransferEventData?.args.sender,
          valueTransferEventData?.args.receiver,
          valueTransferEventData?.args.amount.toString(),
          valueTransferEventData?.args.sequence.toString(),
        ]
      )
    )
  );

  // update pending txn after tree construction
  updatePendingBridgeTransaction(txHash, {
    resumeStep: TransactionStep.GasEstimation,
    tree,
    proof,
  });

  return { tree, proof };
};

// estimate gas (with re-estimation if needed)
const estimateGasStep = async (
  managementContract: ethers.Contract,
  txHash: string,
  valueTransferEventData: ethers.utils.LogDescription,
  proof: any,
  root: string
) => {
  if (!valueTransferEventData) {
    throw new Error("Value transfer event data not found");
  }

  const currentTimestamp = Date.now();

  // Retrieve and parse last gas estimate time from storage
  const lastGasEstimateTimeString = handleStorage.get("lastGasEstimateTime");
  const lastGasEstimateTime = lastGasEstimateTimeString
    ? parseInt(lastGasEstimateTimeString, 10)
    : null;

  console.log("🚀 ~ lastGasEstimateTime:", lastGasEstimateTime);

  // If the gas estimate was done more than a minute ago, re-estimate it
  if (lastGasEstimateTime && currentTimestamp - lastGasEstimateTime > 60000) {
    console.log("Re-estimating gas after timeout");
  }

  const gasLimit = await estimateGasWithTimeout(
    managementContract,
    valueTransferEventData?.args,
    proof,
    root
  );

  // Update pending transaction after gas estimation
  updatePendingBridgeTransaction(txHash, {
    resumeStep: TransactionStep.GasEstimation,
    gasLimit,
  });

  // Store the new timestamp for the gas estimate
  handleStorage.save("lastGasEstimateTime", currentTimestamp.toString());

  return gasLimit;
};

// L1 Relay txn submission
const submitRelayTransactionStep = async (
  txHash: string,
  txResponse: ethers.providers.TransactionResponse,
  l1Provider: ethers.providers.JsonRpcProvider,
  managementContract: ethers.Contract,
  valueTransferEventData: any,
  tree: any,
  proof: any,
  gasLimit: ethers.BigNumber
) => {
  let responseL1 = txResponse;
  if (!responseL1) {
    const gasPrice = await l1Provider.getGasPrice();
    const txL1 =
      await managementContract.populateTransaction.ExtractNativeValue(
        valueTransferEventData?.args,
        proof,
        tree.root,
        { gasPrice, gasLimit }
      );

    if (!txL1) {
      throw new Error("Failed to populate transaction");
    }

    responseL1 = await l1Provider.getSigner().sendTransaction(txL1);
    console.log("🚀 ~ responseL1:", responseL1);

    showToast(
      ToastType.INFO,
      "Transaction sent to L1; waiting for confirmation"
    );

    updatePendingBridgeTransaction(txHash, {
      resumeStep: TransactionStep.RelaySubmission,
      txResponse: responseL1,
    });
  }

  const receiptL1 = await getTransactionReceipt(responseL1);
  console.log("🚀 ~ receiptL1:", receiptL1);

  // finalize and rm txn after L1 confirmation
  removePendingBridgeTransaction(txHash);

  showToast(ToastType.SUCCESS, "Transaction processed successfully");

  return receiptL1;
};

const getTransactionReceipt = async (
  txResponse: ethers.providers.TransactionResponse
) => {
  return await Promise.race([
    txResponse.wait(),
    new Promise((resolve, reject) => {
      const timeoutId = setTimeout(() => {
        reject(new Error("Transaction confirmation timed out"));
      }, 60000 * 2); // 2 minutes

      txResponse
        .wait()
        .then(resolve)
        .catch(reject)
        .finally(() => clearTimeout(timeoutId));
    }),
  ]);
};

const constructMerkleTree = (
  leafEntries: [string, string][],
  msgHash: string
) => {
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
  gasPrice: ethers.BigNumber,
  bridgeContract: ethers.Contract
) => {
  try {
    showToast(ToastType.INFO, "Estimating gas for the transaction");
    const estimatedGas = await estimateGas(receiver, value, bridgeContract);

    showToast(ToastType.INFO, "Populating transaction with estimated gas");
    const populatTxResp = await bridgeContract?.populateTransaction?.sendNative(
      receiver,
      {
        value: ethers.utils.parseEther(value),
        gasPrice,
        gasLimit: estimatedGas,
      }
    );
    return populatTxResp;
  } catch (error) {
    return handleError(error, "Error estimating and populating transaction");
  }
};

const extractAndProcessValueTransfer = async (
  txReceipt: ethers.providers.TransactionReceipt,
  messageBusContract: ethers.Contract,
  provider: ethers.providers.JsonRpcProvider
): Promise<void | {
  valueTransferEventData: ethers.utils.LogDescription;
  block: any;
}> => {
  try {
    showToast(
      ToastType.INFO,
      "Extracting and processing value transfer event data"
    );

    const valueTransferEvent = txReceipt.logs.find(
      (log: ethers.providers.Log) =>
        log.topics[0] ===
        ethers.utils.id("ValueTransfer(address,address,uint256,uint64)")
    );

    if (!valueTransferEvent) {
      throw new Error("ValueTransfer event not found in the logs");
    }

    showToast(ToastType.INFO, "ValueTransfer event found; parsing event data");

    const valueTransferEventData =
      messageBusContract.interface.parseLog(valueTransferEvent);

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
    throw handleError(error, "Error processing value transfer");
  }
};

const estimateGasWithTimeout = async (
  managementContract: ethers.Contract,
  msg: ethers.utils.Result,
  proof: string[],
  root: string,
  timeout = 60000 * 60, // 1 hour
  interval = 5000
): Promise<ethers.BigNumber> => {
  let gasLimit: ethers.BigNumber | null = null;
  const startTime = Date.now();

  while (!gasLimit) {
    showToast(ToastType.INFO, "Estimating gas for value transfer");

    try {
      gasLimit = await managementContract.estimateGas.ExtractNativeValue(
        msg,
        proof,
        root,
        {}
      );
    } catch (error: any) {
      console.error(`Estimate gas threw error: ${error?.reason}`);
    }

    const currentTime = Date.now();
    if (currentTime - startTime >= timeout) {
      console.log("Timed out waiting for gas estimate, using default");
      showToast(
        ToastType.INFO,
        "Timed out waiting for gas estimate; using default"
      );
      return ethers.BigNumber.from(2000000);
    }

    await new Promise((resolve) => setTimeout(resolve, interval));
  }

  return gasLimit as ethers.BigNumber;
};

export {
  constructMerkleTree,
  estimateGasWithTimeout,
  extractAndProcessValueTransfer,
  estimateAndPopulateTx,
  estimateGas,
  sendTransactionStep,
  confirmTransactionStep,
  extractEventDataStep,
  constructMerkleTreeStep,
  estimateGasStep,
  submitRelayTransactionStep,
};
