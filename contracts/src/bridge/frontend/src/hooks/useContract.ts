import { useEffect, useMemo } from "react";
import { ethers } from "ethers";
import useWalletStore from "../stores/wallet-store";
import useContractStore from "../stores/contract-store";
import { toast } from "../components/ui/use-toast";
import { ToastType } from "../types";
import { useGeneralService } from "../services/useGeneralService";
import { privateKey } from "../lib/constants";
import Bridge from "../../artifacts/IBridge.sol/IBridge.json";
import ManagementContractAbi from "../../artifacts/ManagementContract.sol/ManagementContract.json";
import IMessageBusAbi from "../../artifacts/IMessageBus.sol/IMessageBus.json";
import {
  constructMerkleTree,
  estimateGasWithTimeout,
  handleError,
} from "../lib/utils/contractUtils";
import { isAddress } from "ethers/lib/utils";

export const useContract = () => {
  const { signer, isL1ToL2, provider } = useWalletStore();
  const { networkConfig, isNetworkConfigLoading } = useGeneralService();
  const { setContractState, messageBusAddress } = useContractStore();

  const memoizedConfig = useMemo(() => {
    if (isNetworkConfigLoading || !networkConfig) {
      return null;
    }
    return networkConfig;
  }, [networkConfig, isNetworkConfigLoading]);

  const initializeContracts = async () => {
    if (!memoizedConfig || !provider || !signer) return;

    const {
      ImportantContracts: { L1Bridge, L2Bridge },
      MessageBusAddress,
      L2MessageBusAddress,
      ManagementContractAddress,
    } = memoizedConfig;

    let prov = new ethers.providers.Web3Provider(provider);
    const walletInstance = new ethers.Wallet(privateKey as string, prov);
    const isL1 = isL1ToL2;
    const bridgeAddress = isL1 ? L1Bridge : L2Bridge;
    const messageBusAddress = isL1 ? MessageBusAddress : L2MessageBusAddress;

    const bridgeContract = new ethers.Contract(
      bridgeAddress,
      Bridge.abi,
      walletInstance
    );
    const messageBusContract = new ethers.Contract(
      messageBusAddress,
      IMessageBusAbi,
      walletInstance
    );
    const managementContract = new ethers.Contract(
      ManagementContractAddress,
      ManagementContractAbi,
      walletInstance
    );

    setContractState({
      bridgeContract,
      managementContract,
      messageBusContract,
      wallet: walletInstance,
      messageBusAddress,
    });
  };

  useEffect(() => {
    initializeContracts();
  }, [memoizedConfig, provider, isL1ToL2, signer, setContractState]);

  const sendNative = async (receiver: string, value: string) => {
    const { bridgeContract, managementContract, messageBusContract, wallet } =
      useContractStore.getState();

    if (
      !bridgeContract ||
      !signer ||
      !wallet ||
      !managementContract ||
      !messageBusContract
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
        gasPrice,
        bridgeContract
      );
      const txResponse = await signer.sendTransaction(tx);
      console.log("Transaction response:", txResponse);

      toast({
        description: "Transaction sent; waiting for confirmation",
        variant: ToastType.INFO,
      });

      const txReceipt = await txResponse.wait();
      console.log("Transaction receipt:", txReceipt);

      if (isL1ToL2) {
        return txReceipt;
      }

      const { valueTransferEventData, block } =
        await extractAndProcessValueTransfer(txReceipt, messageBusContract);

      const { tree, proof } = constructMerkleTree(
        JSON.parse(atob(block.result.crossChainTree)),
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
        valueTransferEventData?.args,
        proof,
        tree.root
      );

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

      const responseL1 = await wallet.sendTransaction(txL1);
      console.log("Transaction sent to L2:", responseL1);

      toast({
        description: "Value transfer sent to L2; waiting for confirmation",
        variant: ToastType.INFO,
      });

      const receiptL1 = await responseL1.wait();

      toast({
        description: "Value transfer completed",
        variant: ToastType.SUCCESS,
      });

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

    try {
      const balance = await signer.provider?.getBalance(walletAddress);
      return ethers.utils.formatEther(balance);
    } catch (error) {
      return handleError(error, "Error checking Ether balance");
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

  const getBridgeTransactions = async (userAddress: string) => {
    if (!userAddress) {
      return handleError(null, "User address not found");
    }

    try {
      const topics = [
        ethers.utils.id("ValueTransfer(address,address,uint256)"),
        ethers.utils.hexZeroPad(userAddress, 32),
      ];

      const filter = {
        address: messageBusAddress,
        topics,
        fromBlock: 5868682,
      };

      return await provider.getLogs(filter);
    } catch (error) {
      return handleError(error, "Error fetching transactions");
    }
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
      const estimatedGas = await estimateGas(receiver, value, bridgeContract);
      return await bridgeContract?.populateTransaction.sendNative(receiver, {
        value: ethers.utils.parseEther(value),
        gasPrice,
        gasLimit: estimatedGas,
      });
    } catch (error) {
      return handleError(error, "Error populating transaction");
    }
  };

  const extractAndProcessValueTransfer = async (
    txReceipt: any,
    messageBusContract: ethers.Contract
  ) => {
    try {
      toast({
        description: "Extracting logs from the transaction",
        variant: ToastType.INFO,
      });

      const valueTransferEvent = txReceipt.logs.find(
        (log: any) =>
          log.topics[0] ===
          ethers.utils.id("ValueTransfer(address,address,uint256,uint64)")
      );

      if (!valueTransferEvent) {
        throw new Error("ValueTransfer event not found in the logs");
      }

      toast({
        description: "ValueTransfer event found in the logs; processing data",
        variant: ToastType.INFO,
      });

      const valueTransferEventData =
        messageBusContract?.interface.parseLog(valueTransferEvent);

      if (!valueTransferEventData) {
        throw new Error("ValueTransfer event data not found");
      }

      toast({
        description: "ValueTransfer event data found; fetching block",
        variant: ToastType.INFO,
      });

      const block = await provider.send("eth_getBlockByHash", [
        ethers.utils.hexValue(txReceipt.blockHash),
        true,
      ]);

      if (!block) {
        throw new Error("Block not found");
      }

      toast({
        description: "Block found; processing value transfer",
        variant: ToastType.INFO,
      });

      return { valueTransferEventData, block };
    } catch (error) {
      return handleError(error, "Error processing value transfer");
    }
  };

  return {
    sendNative,
    getNativeBalance,
    getTokenBalance,
    sendERC20,
    getBridgeTransactions,
    estimateGas,
  };
};
