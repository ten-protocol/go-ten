import React, { useEffect, useState } from "react";
import { ethers } from "ethers";
import Bridge from "../../artifacts/IBridge.sol/IBridge.json";
import {
  l1Bridge as l1BridgeAddress,
  l2Bridge as l2BridgeAddress,
  messageBusAddress,
} from "../lib/constants";
import { useWalletStore } from "../components/providers/wallet-provider";
import { showToast } from "../components/ui/use-toast";
import { ToastType } from "../types";

export const useContract = () => {
  const [contract, setContract] = useState<ethers.Contract>();
  const [wallet, setWallet] = useState<ethers.Wallet>();
  const { signer, isL1ToL2, provider } = useWalletStore();

  useEffect(() => {
    // var abi = JSON.parse(json)
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
    setContract(contract);
    setWallet(wallet);
  }, [provider, isL1ToL2]);

  // Send native currency to the other network.
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

      const tx = await contract.populateTransaction.sendNative(receiver, {
        value: ethers.utils.parseEther(value),
        gasPrice: gasPrice,
        gasLimit: estimatedGas,
      });

      const txResponse = await signer.sendTransaction(tx);
      const txReceipt = await txResponse.wait();
      console.log("🚀 ~ sendNative ~ txReceipt:", txReceipt);

      // Get the logs from the transaction receipt
      const logs = txReceipt.logs;
      console.log("🚀 ~ sendNative ~ logs:", logs);

      logs.forEach((log: any) => {
        console.log("🚀 ~ logs.forEach ~ log:", log);
        try {
          const parsedLog = contract.interface.parseLog(log);
          console.log(parsedLog);
        } catch (error) {
          console.error("🚀 ~ logs.forEach ~ error:", error);
          // Handle the case where the log does not match the contract's events
        }
      });

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
