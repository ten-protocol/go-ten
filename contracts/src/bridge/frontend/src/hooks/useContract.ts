import { useEffect, useMemo } from "react";
import { ethers } from "ethers";
import useWalletStore from "../stores/wallet-store";
import useContractStore from "../stores/contract-store";
import { toast } from "../components/ui/use-toast";
import { ToastType } from "../types";
import { useGeneralService } from "../services/useGeneralService";
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { privateKey } from "../lib/constants";
import Bridge from "../../artifacts/IBridge.sol/IBridge.json";
import ManagementContractAbi from "../../artifacts/ManagementContract.sol/ManagementContract.json";
import IMessageBusAbi from "../../artifacts/IMessageBus.sol/IMessageBus.json";

export const useContract = () => {
  const { signer, isL1ToL2, provider } = useWalletStore();
  const { obscuroConfig, isObscuroConfigLoading } = useGeneralService();
  const {
    setContractState,
    bridgeContract,
    managementContract,
    messageBusContract,
    wallet,
    messageBusAddress,
  } = useContractStore();

  const memoizedConfig = useMemo(() => {
    if (isObscuroConfigLoading || !obscuroConfig || !obscuroConfig.result) {
      return null;
    }
    return obscuroConfig.result;
  }, [obscuroConfig, isObscuroConfigLoading]);

  useEffect(() => {
    const setupContracts = async () => {
      if (!memoizedConfig || !provider || !signer) {
        return;
      }

      const l1BridgeAddress = memoizedConfig.ImportantContracts.L1Bridge;
      const l2BridgeAddress = memoizedConfig.ImportantContracts.L2Bridge;
      const l1MessageBusAddress = memoizedConfig.MessageBusAddress;
      const l2MessageBusAddress = memoizedConfig.L2MessageBusAddress;
      const managementContractAddress =
        memoizedConfig.ManagementContractAddress;

      const wallet = new ethers.Wallet(privateKey as string, provider);

      const address = isL1ToL2 ? l1BridgeAddress : l2BridgeAddress;
      const messageBusAddress = isL1ToL2
        ? l1MessageBusAddress
        : l2MessageBusAddress;

      const bridgeContract = new ethers.Contract(address, Bridge.abi, wallet);
      const messageBusContract = new ethers.Contract(
        messageBusAddress,
        IMessageBusAbi,
        wallet
      );
      const managementContract = new ethers.Contract(
        managementContractAddress,
        ManagementContractAbi,
        wallet
      );

      setContractState({
        bridgeContract,
        managementContract,
        messageBusContract,
        wallet,
        messageBusAddress,
      });
    };

    setupContracts();
  }, [memoizedConfig, provider, isL1ToL2, signer, setContractState]);

  const sendNative = async (receiver: string, value: string) => {
    const { bridgeContract, managementContract, messageBusContract, wallet } =
      useContractStore.getState();

    if (!bridgeContract || !signer || !wallet || !managementContract) {
      console.error("Contract or signer not found");
      throw new Error("Contract or signer not found");
    }

    try {
      if (!ethers.utils.isAddress(receiver)) {
        throw new Error("Invalid receiver address");
      }

      const gasPrice = await signer.provider.getGasPrice();
      const estimatedGas = await bridgeContract.estimateGas.sendNative(
        receiver,
        { value: ethers.utils.parseEther(value) }
      );

      const tx = await bridgeContract.populateTransaction.sendNative(receiver, {
        value: ethers.utils.parseEther(value),
        gasPrice,
        gasLimit: estimatedGas,
      });

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

      const abiTypes = ["address", "address", "uint256", "uint64"];
      const _msg = [
        valueTransferEventData.args.sender,
        valueTransferEventData.args.receiver,
        valueTransferEventData.args.amount.toString(),
        valueTransferEventData.args.sequence.toString(),
      ];

      const abiCoder = new ethers.utils.AbiCoder();
      const encodedMsg = abiCoder.encode(abiTypes, _msg);
      const processedValueTransfer = [_msg, ethers.utils.keccak256(encodedMsg)];
      const msg = processedValueTransfer[0];
      const msgHash = processedValueTransfer[1];

      const base64CrossChainTree = block.result.crossChainTree;
      const decodedCrossChainTree = atob(base64CrossChainTree);

      const leafEntries = JSON.parse(decodedCrossChainTree);

      if (leafEntries[0][1] === msgHash) {
        toast({
          description: "Value transfer hash is in the xchain tree",
          variant: ToastType.INFO,
        });
      }

      toast({
        description: "Constructing merkle tree",
        variant: ToastType.INFO,
      });

      const tree = StandardMerkleTree.of(leafEntries, ["string", "bytes32"]);
      toast({
        description: "Merkle tree constructed",
        variant: ToastType.INFO,
      });

      const proof = tree.getProof(["v", msgHash]);
      const root = tree.root;

      if (block.result.crossChainTreeHash === tree.root) {
        toast({
          description:
            "Constructed merkle root matches block crossChainTreeHash",
          variant: ToastType.INFO,
        });
      }

      let gasLimit: ethers.BigNumber | null = null;
      toast({
        description: "Estimating gas",
        variant: ToastType.INFO,
      });

      const estimateGasWithTimeout = async (
        timeout = 30000,
        interval = 5000
      ) => {
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
            toast({
              description: "Timed out waiting for gas estimate, using default",
              variant: ToastType.INFO,
            });
            return ethers.BigNumber.from(2000000);
          }
          await new Promise((resolve) => setTimeout(resolve, interval));
        }
        return gasLimit;
      };

      gasLimit = await estimateGasWithTimeout();
      toast({
        description: "Sending value transfer to L2",
        variant: ToastType.INFO,
      });

      const txL1: ethers.PopulatedTransaction =
        await managementContract.populateTransaction.ExtractNativeValue(
          msg,
          proof,
          root,
          { gasPrice, gasLimit }
        );

      console.log("🚀 ~ sendNative ~ txL1:", txL1);

      const responseL1 = await wallet.sendTransaction(txL1);
      console.log("🚀 ~ sendNative ~ responseL1:", responseL1);
      toast({
        description: "Value transfer sent to L2; waiting for confirmation",
        variant: ToastType.INFO,
      });

      const receiptL1 = await responseL1.wait();
      console.log("🚀 ~ sendNative ~ receiptL1:", receiptL1);
      return receiptL1;
    } catch (error) {
      console.error("Error sending native currency:", error);
      throw error;
    }
  };

  const sendERC20 = async (
    receiver: string,
    amount: string,
    tokenContractAddress: string
  ) => {
    const { bridgeContract } = useContractStore.getState();

    if (!bridgeContract) {
      console.error("Contract not found");
      return null;
    }
    return bridgeContract.sendERC20(tokenContractAddress, amount, receiver);
  };

  const getNativeBalance = async (provider: any, walletAddress: string) => {
    if (!provider || !walletAddress) {
      console.error("Provider or wallet address not found");
      return null;
    }

    try {
      const balance = await provider.getBalance(walletAddress);
      return ethers.utils.formatEther(balance);
    } catch (error) {
      console.error("Error checking Ether balance:", error);
      throw error;
    }
  };

  const getTokenBalance = async (
    tokenAddress: string,
    walletAddress: string,
    provider: any
  ) => {
    if (!provider || !walletAddress) {
      console.error("Provider or wallet address not found");
      return null;
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
      console.error("Error checking token balance:", error);
      throw error;
    }
  };

  const getBridgeTransactions = async (provider: any, userAddress: string) => {
    const { messageBusAddress } = useContractStore.getState();

    if (!provider || !userAddress || !messageBusAddress) {
      console.error("Provider, user address, or message bus address not found");
      return null;
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
      console.error("Error fetching transactions:", error);
      throw error;
    }
  };

  return {
    bridgeContract,
    managementContract,
    messageBusContract,
    wallet,
    messageBusAddress,
    sendNative,
    sendERC20,
    getNativeBalance,
    getTokenBalance,
    getBridgeTransactions,
  };
};
