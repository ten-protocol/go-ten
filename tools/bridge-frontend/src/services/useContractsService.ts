import { useEffect, useMemo, useState } from "react";
import { ethers } from "ethers";
import useContractStore from "../stores/contract-store";
import {
  IPendingTx,
  ToastType,
  TransactionStatus,
  TransactionStep,
} from "../types";
import { useGeneralService } from "./useGeneralService";
import Bridge from "../../artifacts/IBridge.sol/IBridge.json";
import ManagementContractAbi from "../../artifacts/ManagementContract.sol/ManagementContract.json";
import IMessageBusAbi from "../../artifacts/IMessageBus.sol/IMessageBus.json";
import { isAddress } from "ethers/lib/utils";
import useWalletStore from "../stores/wallet-store";
import { currentNetwork } from "../lib/utils";
import {
  confirmTransactionStep,
  constructMerkleTreeStep,
  estimateGasStep,
  extractEventDataStep,
  sendTransactionStep,
  submitRelayTransactionStep,
} from "../lib/utils/contractUtils";
import { handleError } from "../lib/utils/walletUtils";
import {
  getPendingBridgeTransactions,
  removePendingBridgeTransaction,
} from "../lib/utils/txnUtils";
import { showToast } from "../components/ui/use-toast";
import { useQueryClient } from "@tanstack/react-query";

export const useContractsService = () => {
  const queryClient = useQueryClient();
  const { isL1ToL2, signer, provider, address } = useWalletStore();
  const { networkConfig, isNetworkConfigLoading } = useGeneralService();
  const {
    setContractState,
    messageBusAddress,
    bridgeAddress,
    bridgeContract,
    managementContract,
    messageBusContract,
  } = useContractStore();
  const [finalisingTxHashes, setFinalisingTxHashes] = useState<Set<string>>(
    new Set()
  );

  const memoizedConfig = useMemo(() => {
    if (isNetworkConfigLoading || !networkConfig) {
      return null;
    }
    return networkConfig;
  }, [networkConfig, isNetworkConfigLoading]);

  const l1Provider = new ethers.providers.JsonRpcProvider(currentNetwork.l1Rpc);
  const initializeContracts = async () => {
    if (!memoizedConfig || !provider) return;

    const {
      ImportantContracts: { L1Bridge, L2Bridge },
      MessageBusAddress,
      L2MessageBusAddress,
      ManagementContractAddress,
    } = memoizedConfig;

    const signer = provider.getSigner();
    const bridgeAddress = isL1ToL2 ? L1Bridge : L2Bridge;
    const messageBusAddress = isL1ToL2
      ? MessageBusAddress
      : L2MessageBusAddress;

    const bridgeContract = new ethers.Contract(
      bridgeAddress,
      Bridge.abi,
      signer
    );
    const messageBusContract = new ethers.Contract(
      messageBusAddress,
      IMessageBusAbi,
      signer
    );
    const managementContract = new ethers.Contract(
      ManagementContractAddress,
      ManagementContractAbi,
      l1Provider.getSigner()
    );

    setContractState({
      bridgeContract,
      managementContract,
      messageBusContract,
      messageBusAddress,
      bridgeAddress,
    });
  };
  useEffect(() => {
    initializeContracts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [memoizedConfig, provider, isL1ToL2, signer, setContractState]);

  // main entry function that calls each step, either from the beginning or resuming
  const sendNative = async (tx: IPendingTx) => {
    if (!provider) {
      return handleError(null, "Provider not found");
    }
    const signer = provider.getSigner();

    if (!bridgeContract || !managementContract || !messageBusContract) {
      return handleError(null, "Contract not found");
    }

    const { receiver, value } = tx;

    if (!receiver) {
      return handleError(null, "Receiver address not found");
    }

    if (!value) {
      return handleError(null, "Value not found");
    }

    if (!ethers.utils.isAddress(receiver)) {
      return handleError(null, "Invalid receiver address");
    }

    try {
      let {
        txHash,
        resumeStep,
        txResponse,
        txReceipt,
        valueTransferEventData,
        block,
        tree,
        root,
        proof,
        gasLimit,
      } = tx;

      let currentStep = resumeStep || TransactionStep.TransactionSubmission;

      while (true) {
        switch (currentStep) {
          case TransactionStep.TransactionSubmission:
            if (!txHash) {
              txResponse = await sendTransactionStep(
                bridgeContract,
                signer,
                receiver,
                value
              );

              if (!txResponse) {
                return handleError(null, "Error sending transaction");
              }

              txHash = txResponse.hash;
            } else {
              txResponse = await provider.getTransaction(txHash);

              if (!txResponse) {
                return handleError(
                  null,
                  "Transaction not found on the network"
                );
              }
            }

            currentStep = TransactionStep.TransactionConfirmation;
            showToast(ToastType.INFO, "Transaction submitted");
            break;

          case TransactionStep.TransactionConfirmation:
            if (!txReceipt) {
              if (txResponse) {
                txReceipt = await confirmTransactionStep(txResponse!);
              }

              // ...if we still don't have a txReceipt, fetch it w txHash
              if (!txReceipt && txHash) {
                txReceipt = await provider.getTransactionReceipt(txHash);

                if (!txReceipt) {
                  // if the receipt is still not confirmed, we can retry later or throw an error
                  return handleError(
                    null,
                    "Transaction is not yet confirmed. Please retry later."
                  );
                }
              }

              // if no receipt is found at all
              if (!txReceipt) {
                return handleError(null, "Error confirming transaction");
              }
            }

            // for L1 > L2, we can skip the rest of the steps
            if (isL1ToL2) {
              txHash && removePendingBridgeTransaction(txHash);
              queryClient.invalidateQueries({
                queryKey: ["bridgePendingTransactions", isL1ToL2 ? "l1" : "l2"],
              });
              return txReceipt;
            }

            currentStep = TransactionStep.EventDataExtraction;
            showToast(ToastType.INFO, "Transaction confirmed");
            break;

          case TransactionStep.EventDataExtraction:
            if (!txReceipt && txHash) {
              txReceipt = await provider.getTransactionReceipt(txHash);

              if (!txReceipt) {
                return handleError(null, "Transaction is not yet confirmed");
              }
            }
            const result = await extractEventDataStep(
              messageBusContract,
              provider,
              txReceipt!
            );
            valueTransferEventData = result.valueTransferEventData;
            block = result.block;

            if (!valueTransferEventData || !block) {
              return handleError(null, "Error extracting event data");
            }
            currentStep = TransactionStep.MerkleTreeConstruction;
            showToast(ToastType.INFO, "Event data extracted");
            break;

          case TransactionStep.MerkleTreeConstruction:
            ({ tree, proof } = await constructMerkleTreeStep(
              txHash!,
              valueTransferEventData!,
              block
            ));

            if (!tree || !proof) {
              return handleError(null, "Error constructing Merkle tree");
            }

            currentStep = TransactionStep.GasEstimation;
            showToast(ToastType.INFO, "Merkle tree constructed");
            break;

          case TransactionStep.GasEstimation:
            gasLimit = await estimateGasStep(
              managementContract,
              txHash!,
              valueTransferEventData!,
              proof,
              root || tree!.root
            );

            if (!gasLimit) {
              return handleError(null, "Error estimating gas");
            }

            currentStep = TransactionStep.RelaySubmission;
            showToast(ToastType.INFO, "Gas estimated for relay");
            break;

          case TransactionStep.RelaySubmission:
            txReceipt = (await submitRelayTransactionStep(
              txHash!,
              txResponse!,
              l1Provider,
              managementContract,
              valueTransferEventData,
              root || tree!.root,
              proof,
              gasLimit!
            )) as ethers.providers.TransactionReceipt;

            queryClient.invalidateQueries({
              queryKey: ["bridgePendingTransactions", isL1ToL2 ? "l1" : "l2"],
            });

            return txReceipt;

          default:
            return handleError(null, "Invalid transaction step");
        }
      }
    } catch (error) {
      return handleError(error, "Error sending native currency");
    }
  };

  const sendERC20 = async (
    receiver: string,
    amount: string,
    tokenContractAddress: string
  ) => {
    const { bridgeContract } = useContractStore.getState();

    if (!bridgeContract) {
      return handleError(null, "Contract not found");
    }
    return bridgeContract.sendERC20(tokenContractAddress, amount, receiver);
  };

  const getNativeBalance = async (walletAddress: string) => {
    if (!walletAddress || !isAddress(walletAddress)) {
      return handleError(null, "Invalid wallet address");
    }

    if (!provider) {
      return handleError(null, "Signer or provider not found");
    }

    try {
      const balance = await provider?.getBalance(walletAddress);
      return ethers.utils.formatEther(balance);
    } catch (error) {
      handleError(error, "Error checking Ether balance");
    }
  };

  const getTokenBalance = async (
    tokenAddress: string,
    walletAddress: string
  ) => {
    if (!provider || !walletAddress) {
      return handleError(null, "Provider or wallet address not found");
    }

    try {
      const tokenContract = new ethers.Contract(
        tokenAddress,
        [
          "function balanceOf(address owner) view returns (uint256)",
          "function decimals() view returns (uint8)",
        ],
        provider
      );

      const balance = await tokenContract.balanceOf(walletAddress);
      const decimals = await tokenContract.decimals();
      return ethers.utils.formatUnits(balance, decimals);
    } catch (error) {
      return handleError(error, "Error checking token balance");
    }
  };

  const getBridgeTransactions = async () => {
    if (!provider || !messageBusAddress || !bridgeAddress) {
      return handleError(null, "Provider or contract address not found");
    }
    try {
      const topics = [
        ethers.utils.id("ValueTransfer(address,address,uint256,uint64)"),
        ethers.utils.hexZeroPad(bridgeAddress, 32),
        ethers.utils.hexZeroPad(address, 32),
      ];

      const filter = {
        address: messageBusAddress,
        topics,
      };

      const logs = await provider.getLogs(filter);
      const transactions = await Promise.all(
        logs.map(async (log: ethers.providers.Log) => {
          const receipt = await provider.getTransactionReceipt(
            log.transactionHash
          );
          return {
            ...log,
            status: receipt
              ? receipt.status
                ? TransactionStatus.Success
                : TransactionStatus.Failure
              : TransactionStatus.Pending,
          };
        })
      );
      return transactions;
    } catch (error) {
      return handleError(error, "Error fetching transactions");
    }
  };

  const resumePendingTransactions = () => {
    const pendingTransactions = getPendingBridgeTransactions();

    pendingTransactions.forEach(async (tx: IPendingTx) => {
      try {
        await finaliseTransaction(tx);
      } catch (error) {
        handleError(error, `Error resuming transaction ${tx.txHash}:`);
      }
    });
  };

  const finaliseTransaction = async (tx: IPendingTx) => {
    console.log("🚀 ~ finalizeTransaction ~ tx:", tx);
    try {
      setFinalisingTxHashes((prevSet) =>
        new Set(prevSet).add(tx?.txHash || "")
      );
      showToast(ToastType.INFO, "Resuming transaction...");
      await sendNative({ ...tx });
      setFinalisingTxHashes((prevSet) => {
        const newSet = new Set(prevSet);
        newSet.delete(tx?.txHash || "");
        return newSet;
      });
    } catch (error) {
      handleError(error, `Error resuming transaction ${tx.txHash}:`);
      setFinalisingTxHashes((prevSet) => {
        const newSet = new Set(prevSet);
        newSet.delete(tx?.txHash || "");
        return newSet;
      });
    }
  };

  return {
    sendNative,
    getNativeBalance,
    getTokenBalance,
    sendERC20,
    getBridgeTransactions,
    finaliseTransaction,
    resumePendingTransactions,
    finalisingTxHashes,
  };
};
