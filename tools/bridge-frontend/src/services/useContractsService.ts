import { useEffect, useMemo } from "react";
import { ethers } from "ethers";
import useContractStore from "../stores/contract-store";
import { ToastType, TransactionStatus } from "../types";
import { useGeneralService } from "./useGeneralService";
import Bridge from "../../artifacts/IBridge.sol/IBridge.json";
import ManagementContractAbi from "../../artifacts/ManagementContract.sol/ManagementContract.json";
import IMessageBusAbi from "../../artifacts/IMessageBus.sol/IMessageBus.json";
import {
  constructMerkleTree,
  estimateAndPopulateTx,
  estimateGasWithTimeout,
  extractAndProcessValueTransfer,
  handleError,
} from "../lib/utils/contractUtils";
import { isAddress } from "ethers/lib/utils";
import { showToast } from "../components/ui/use-toast";
import useWalletStore from "../stores/wallet-store";
import { currentNetwork } from "../lib/utils";

export const useContractsService = () => {
  const { signer, isL1ToL2, provider, address } = useWalletStore();
  const { networkConfig, isNetworkConfigLoading } = useGeneralService();
  const {
    setContractState,
    messageBusAddress,
    bridgeAddress,
    bridgeContract,
    managementContract,
    messageBusContract,
  } = useContractStore();

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
    const isL1 = isL1ToL2;
    const bridgeAddress = isL1 ? L1Bridge : L2Bridge;
    const messageBusAddress = isL1 ? MessageBusAddress : L2MessageBusAddress;

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

  const sendNative = async (receiver: string, value: string) => {
    if (
      !bridgeContract ||
      !signer ||
      !managementContract ||
      !messageBusContract ||
      !provider
    ) {
      return handleError(null, "Contract or signer not found");
    }

    if (!ethers.utils.isAddress(receiver)) {
      return handleError(null, "Invalid receiver address");
    }

    try {
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

      showToast(ToastType.INFO, "Transaction sent; waiting for confirmation");

      const txReceipt = await txResponse.wait();
      console.log("Transaction receipt:", txReceipt);

      if (isL1ToL2) {
        return txReceipt;
      }

      const { valueTransferEventData, block } =
        (await extractAndProcessValueTransfer(
          txReceipt,
          messageBusContract,
          provider
        )) || { valueTransferEventData: null, block: null };

      if (!valueTransferEventData || !block) {
        throw new Error("Failed to extract value transfer event data");
      }

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

      const gasLimit = await estimateGasWithTimeout(
        managementContract,
        valueTransferEventData?.args as (string | ethers.BigNumber)[],
        proof,
        tree.root
      );

      console.log("🚀 ~ sendNative ~ gasLimit:", gasLimit);
      const txL1 =
        await managementContract.populateTransaction.ExtractNativeValue(
          valueTransferEventData?.args,
          proof,
          tree.root,
          { gasPrice, gasLimit }
        );
      console.log("🚀 ~ sendNative ~ txL1:", txL1);

      if (!txL1) {
        throw new Error("Failed to populate transaction");
      }
      const responseL1 = await l1Provider.getSigner().sendTransaction(txL1);
      console.log("🚀 ~ sendNative ~ responseL1:", responseL1);

      showToast(
        ToastType.INFO,
        "Transaction sent to L1; waiting for confirmation"
      );

      const receiptL1 = await responseL1.wait();
      console.log("L1 txn receipt:", receiptL1);

      showToast(ToastType.SUCCESS, "Transaction processed successfully");

      return receiptL1;
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

  return {
    sendNative,
    getNativeBalance,
    getTokenBalance,
    sendERC20,
    getBridgeTransactions,
  };
};
