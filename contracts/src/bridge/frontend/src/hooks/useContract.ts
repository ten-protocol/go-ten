import React, { useEffect, useMemo, useState } from "react";
import axios from "axios";
import { ethers } from "ethers";
import Bridge from "../../artifacts/IBridge.sol/IBridge.json";
import ManagementContractAbi from "../../artifacts/ManagementContract.sol/ManagementContract.json";
import { useWalletStore } from "../components/providers/wallet-provider";
import { showToast } from "../components/ui/use-toast";
import { ToastType } from "../types";
import { useGeneralService } from "../services/useGeneralService";
import { environment } from "../lib/constants";
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";

export const useContract = () => {
  const [contract, setContract] = useState<ethers.Contract>();
  const [managementContract, setManagementContract] =
    useState<ethers.Contract>();
  const [wallet, setWallet] = useState<ethers.Wallet>();
  const { signer, isL1ToL2, provider } = useWalletStore();
  const { obscuroConfig, isObscuroConfigLoading } = useGeneralService();

  const [messageBusAddress, setMessageBusAddress] = useState<string>("");

  const memoizedConfig = useMemo(() => {
    if (isObscuroConfigLoading || !obscuroConfig || !obscuroConfig.result) {
      return null;
    }
    return obscuroConfig.result;
  }, [obscuroConfig, isObscuroConfigLoading]);

  useEffect(() => {
    // if (isObscuroConfigLoading) {
    //   return;
    // }
    // if (!obscuroConfig) {
    //   showToast(ToastType.DESTRUCTIVE, "Config not found");
    //   return;
    // }

    // if (!memoizedConfig) {
    //   showToast(ToastType.DESTRUCTIVE, "Config not found");
    //   return;
    // }

    // sepolia
    const l1BridgeAddress = "0x87B99D9709764b72C0c25Dffd1Eed5D1014b5F6C";
    const l2BridgeAddress = "0xFeD8Fc00f96d652244c6EE628da65Ea766CcEc81";
    const managementContractAddress =
      "0x1fc65b3639643d82848BfCc3a0D5a2D5881965Ab";

    // const l1BridgeAddress = memoizedConfig.ImportantContracts.L1Bridge;
    // const l2BridgeAddress = memoizedConfig.ImportantContracts.L2Bridge;
    // const messageBusAddress = memoizedConfig.MessageBusAddress;
    // setMessageBusAddress(messageBusAddress);

    // uat
    // const l1BridgeAddress = "0xb68B5B0Ec17AA746F30CF3f1E51b6ace56B7eCCF";
    // const l2BridgeAddress = "0xFeD8Fc00f96d652244c6EE628da65Ea766CcEc81";
    // const managementContractAddress =
    //   "0xc23D7eDdC53b235c30B02c05F753929a092B3872";

    if (!provider) {
      console.error("Provider not found");
      return;
    }
    const p = new ethers.providers.Web3Provider(provider);
    const wallet = new ethers.Wallet(
      process.env.NEXT_PUBLIC_PRIVATE_KEY as string,
      p
    );
    const address = isL1ToL2 ? l1BridgeAddress : l2BridgeAddress;
    const contract = new ethers.Contract(address as string, Bridge.abi, wallet);
    const managementContract = new ethers.Contract(
      managementContractAddress as string,
      ManagementContractAbi,
      wallet
    );
    console.log("🚀 ~ useEffect ~ contract:", contract);
    setContract(contract);
    setManagementContract(managementContract);
    setWallet(wallet);

    console.log("🚀 ~ config", memoizedConfig);
  }, [provider, isL1ToL2]);
  // }, [obscuroConfig, isObscuroConfigLoading, memoizedConfig]);

  const sendNative = async (receiver: string, value: string) => {
    if (!contract) {
      console.error("Contract not found");
      return null;
    }
    try {
      if (!ethers.utils.isAddress(receiver)) {
        console.error("Invalid address");
        return null;
      }

      const gasPrice = await signer.provider.getGasPrice();
      const estimatedGas = await contract.estimateGas.sendNative(receiver, {
        value: ethers.utils.parseEther(value),
      });

      // 1. Call into the Ethereum bridge sendNative method to initiate the transfer of funds from the L2 to the L1
      const tx = await contract.populateTransaction.sendNative(receiver, {
        value: ethers.utils.parseEther(value),
        gasPrice: gasPrice,
        gasLimit: estimatedGas,
      });

      const txResponse = await signer.sendTransaction(tx);
      const txReceipt = await txResponse.wait();
      console.log("🚀 ~ sendNative ~ txReceipt:", txReceipt);

      const valueTransferEvent = txReceipt.logs.find(
        (log: any) =>
          log.topics[0] ===
          ethers.utils.id("ValueTransfer(address,address,uint256,uint64)")
      );
      console.log("🚀 ~ sendNative ~ valueTransferEvent:", valueTransferEvent);
      if (!valueTransferEvent) {
        throw new Error("ValueTransfer event not found in the logs");
      }
      // Debug: Log all topics
      txReceipt.logs.forEach((log: any, index: number) => {
        console.log(`Log ${index}:`, log.topics[0]);
      });
      const valueTransferEventData =
        contract.interface.parseLog(valueTransferEvent);
      console.log(
        "🚀 ~ sendNative ~ valueTransferEventData:",
        valueTransferEventData
      );

      const hashedValueTransfer = ethers.utils.keccak256(
        ethers.utils.defaultAbiCoder.encode(
          ["address", "address", "uint256", "uint64"],
          [
            valueTransferEventData.args.sender,
            valueTransferEventData.args.receiver,
            valueTransferEventData.args.amount,
            valueTransferEventData.args.sequence,
          ]
        )
      );

      const blockHash = txReceipt.blockHash;
      const block = await signer.provider.getBlock(blockHash);

      const base64CrossChainTree = block.extraData; // Assuming crossChainTree is in extraData
      const decodedCrossChainTree = Buffer.from(
        base64CrossChainTree,
        "base64"
      ).toString("utf8");
      const leafEntries = JSON.parse(decodedCrossChainTree);

      const leaves = leafEntries.map((entry: any) =>
        ethers.utils.keccak256(entry)
      );
      const tree = new StandardMerkleTree(leaves);
      const root = tree.getHexRoot();
      const proof = tree.getHexProof(hashedValueTransfer);

      // const managementContractL1 = new ethers.Contract(MANAGEMENT_CONTRACT_ADDRESS_L1, MANAGEMENT_CONTRACT_ABI_L1, walletL1);

      let gasLimit;
      while (true) {
        try {
          gasLimit = await managementContract?.estimateGas.ExtractNativeValue(
            hashedValueTransfer,
            proof,
            root
          );
          break;
        } catch (error) {
          // Retry until it stops failing
        }
      }

      const txL1 = await managementContract?.ExtractNativeValue(
        hashedValueTransfer,
        proof,
        root,
        { gasLimit }
      );
      const receiptL1 = await txL1.wait();
      console.log("Funds released:", receiptL1);

      return txReceipt;
    } catch (error) {
      console.error("Error sending native currency:", error);
      throw error;
    }
  };

  // Send ERC20 assets to the other network.
  async function sendERC20(
    receiver: string,
    amount: string,
    tokenContractAddress: string
  ) {
    if (!contract) {
      console.error("Contract not found");
      return null;
    }
    return contract.sendERC20(tokenContractAddress, amount, receiver);
  }

  async function getNativeBalance(provider: any, walletAddress: string) {
    if (!provider) {
      console.error("Provider not found");
      return null;
    }
    if (!walletAddress) {
      console.error("Wallet address not found");
      return null;
    }
    try {
      const p = new ethers.providers.Web3Provider(provider);
      const balance = await p.getBalance(walletAddress);
      const readableBalance = ethers.utils.formatEther(balance);
      return readableBalance;
    } catch (error) {
      console.error("Error checking Ether balance:", error);
      throw error;
    }
  }

  // balance of an ERC20 asset.
  async function getTokenBalance(
    tokenAddress: string,
    walletAddress: string,
    provider: any
  ) {
    if (!provider) {
      console.error("Provider not found");
      return null;
    }
    if (!walletAddress) {
      console.error("Wallet address not found");
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

      // Call the balanceOf function of the ERC20 contract
      const balance = await tokenContract.balanceOf(walletAddress);

      // Call the decimals function of the ERC20 contract
      const decimals = await tokenContract.decimals();

      // Convert balance to a readable format
      const readableBalance = ethers.utils.formatUnits(balance, decimals);

      return readableBalance;
    } catch (error) {
      console.error("Error checking token balance:", error);
      throw error;
    }
  }

  //get bridge transactions
  async function getBridgeTransactions(provider: any, userAddress: string) {
    try {
      if (!provider) {
        console.error("Provider not found");
        return null;
      }
      if (!userAddress) {
        console.error("User address not found");
        return null;
      }
      if (!messageBusAddress) {
        console.error("Message bus address not found");
        return null;
      }
      const p = new ethers.providers.Web3Provider(provider);

      // using topics to filter logs and get only the user's transactions
      const topics = [
        ethers.utils.id("ValueTransfer(address,address,uint256)"), // Event signature
        ethers.utils.hexZeroPad(userAddress, 32),
        // ethers.utils.hexZeroPad(messageBusAddress, 32),
      ];

      const filter = {
        address: messageBusAddress,
        topics,
        fromBlock: 5868682,
      };

      const transactions = await p.getLogs(filter);

      console.log("Transactions on bridge:", transactions);
      return transactions;
    } catch (error) {
      console.error("Error fetching transactions:", error);
      throw error;
    }
  }

  return {
    contract,
    sendNative,
    sendERC20,
    getNativeBalance,
    getTokenBalance,
    getBridgeTransactions,
  };
};
