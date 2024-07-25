import { useEffect, useMemo, useState } from "react";
import { ethers } from "ethers";
import Bridge from "../../artifacts/IBridge.sol/IBridge.json";
import ManagementContractAbi from "../../artifacts/ManagementContract.sol/ManagementContract.json";
import IMessageBusAbi from "../../artifacts/IMessageBus.sol/IMessageBus.json";
import { useWalletStore } from "../components/providers/wallet-provider";
import { showToast } from "../components/ui/use-toast";
import { ToastType } from "../types";
import { useGeneralService } from "../services/useGeneralService";
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";

interface ContractState {
  bridgeContract?: ethers.Contract;
  managementContract?: ethers.Contract;
  messageBusContract?: ethers.Contract;
  wallet?: ethers.Wallet;
  messageBusAddress: string;
}

const useContract = () => {
  const [contractState, setContractState] = useState<ContractState>({
    messageBusAddress: "",
  });

  const { signer, isL1ToL2, provider } = useWalletStore();
  const { obscuroConfig, isObscuroConfigLoading } = useGeneralService();

  const memoizedConfig = useMemo(() => {
    if (isObscuroConfigLoading || !obscuroConfig || !obscuroConfig.result) {
      return null;
    }
    return obscuroConfig.result;
  }, [obscuroConfig, isObscuroConfigLoading]);

  useEffect(() => {
    if (isObscuroConfigLoading) {
      return;
    }
    if (!obscuroConfig) {
      showToast(ToastType.DESTRUCTIVE, "Config not found");
      return;
    }

    if (!memoizedConfig) {
      showToast(ToastType.DESTRUCTIVE, "Config not found");
      return;
    }

    const l1BridgeAddress = memoizedConfig.ImportantContracts.L1Bridge;
    const l2BridgeAddress = memoizedConfig.ImportantContracts.L2Bridge;
    const l1MessageBusAddress = memoizedConfig.MessageBusAddress;
    const l2MessageBusAddress = memoizedConfig.L2MessageBusAddress;
    const managementContractAddress = memoizedConfig.ManagementContractAddress;

    if (!provider) {
      throw new Error("Provider not found");
    }
    const p = new ethers.providers.Web3Provider(provider);
    const wallet = new ethers.Wallet(
      process.env.NEXT_PUBLIC_PRIVATE_KEY as string,
      p
    );
    const address = isL1ToL2 ? l1BridgeAddress : l2BridgeAddress;
    const messageBusAddress = isL1ToL2
      ? l1MessageBusAddress
      : l2MessageBusAddress;
    const bridgeContract = new ethers.Contract(
      address as string,
      Bridge.abi,
      wallet
    );
    const messageBusContract = new ethers.Contract(
      messageBusAddress as string,
      IMessageBusAbi,
      wallet
    );
    const managementContract = new ethers.Contract(
      managementContractAddress as string,
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

    console.log("🚀 ~ config", memoizedConfig);
  }, [provider, isL1ToL2]);
  // }, [obscuroConfig, isObscuroConfigLoading, memoizedConfig]);

  const sendNative = async (receiver: string, value: string) => {
    const { bridgeContract, managementContract, messageBusContract, wallet } =
      contractState;

    if (!bridgeContract || !signer || !wallet) {
      console.error("Required contract or wallet not found");
      return null;
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

      const txReceipt = await txResponse.wait();
      console.log("Transaction receipt:", txReceipt);

      if (isL1ToL2) {
        return txReceipt;
      }

      const valueTransferEvent = txReceipt.logs.find(
        (log: any) =>
          log.topics[0] ===
          ethers.utils.id("ValueTransfer(address,address,uint256,uint64)")
      );

      if (!valueTransferEvent) {
        throw new Error("ValueTransfer event not found in the logs");
      }

      const valueTransferEventData =
        messageBusContract?.interface.parseLog(valueTransferEvent);

      if (!valueTransferEventData) {
        throw new Error("ValueTransfer event data not found");
      }

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

      const block = await provider.send("eth_getBlockByHash", [
        ethers.utils.hexValue(txReceipt.blockHash),
        true,
      ]);

      const base64CrossChainTree = block.result.crossChainTree;
      const decodedCrossChainTree = atob(base64CrossChainTree);
      console.log(
        "🚀 ~ sendNative ~ decodedCrossChainTree:",
        decodedCrossChainTree
      );

      const leafEntries = JSON.parse(decodedCrossChainTree);

      if (decodedCrossChainTree[0][1] === msgHash) {
        console.log("Value transfer hash is in the xchain tree");
      }

      const tree = StandardMerkleTree.of(leafEntries, ["string", "bytes32"]);
      const proof = tree.getProof(["v", msgHash]);

      const root = tree.root;

      if (block.result.crossChainTreeHash === tree.root) {
        console.log("Constructed merkle root matches block crossChainTreeHash");
      }

      let gasLimit: ethers.BigNumber | undefined;

      const estimateGasWithTimeout = async (
        timeout = 30000,
        interval = 5000
      ) => {
        const startTime = Date.now();
        while (!gasLimit) {
          try {
            gasLimit = await managementContract?.estimateGas.ExtractNativeValue(
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

      gasLimit = await estimateGasWithTimeout();

      const txL1: ethers.PopulatedTransaction | undefined =
        await managementContract?.populateTransaction.ExtractNativeValue(
          msg,
          proof,
          root,
          { gasPrice, gasLimit }
        );

      if (!txL1) {
        console.error("Transaction not found");
        return null;
      }
      const responseL1 = await wallet.sendTransaction(txL1);
      const receiptL1 = await responseL1.wait();

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
    const { bridgeContract } = contractState;

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
      const p = new ethers.providers.Web3Provider(provider);
      const balance = await p.getBalance(walletAddress);
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
      const p = new ethers.providers.Web3Provider(provider);
      const tokenContract = new ethers.Contract(
        tokenAddress,
        [
          "function balanceOf(address owner) view returns (uint256)",
          "function decimals() view returns (uint8)",
        ],
        p
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
    const { messageBusAddress } = contractState;

    if (!provider || !userAddress || !messageBusAddress) {
      console.error("Provider, user address, or message bus address not found");
      return null;
    }

    try {
      const p = new ethers.providers.Web3Provider(provider);

      const topics = [
        ethers.utils.id("ValueTransfer(address,address,uint256)"),
        ethers.utils.hexZeroPad(userAddress, 32),
      ];

      const filter = {
        address: messageBusAddress,
        topics,
        fromBlock: 5868682,
      };

      return await p.getLogs(filter);
    } catch (error) {
      console.error("Error fetching transactions:", error);
      throw error;
    }
  };

  return {
    ...contractState,
    sendNative,
    sendERC20,
    getNativeBalance,
    getTokenBalance,
    getBridgeTransactions,
  };
};

export default useContract;
