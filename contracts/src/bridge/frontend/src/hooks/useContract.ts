import React from "react";
import { ethers } from "ethers";
import L1Bridge from "../../artifacts/IBridge.sol/IBridge.json";
import {
  l1Bridge as l1BridgeAddress,
  messageBusAddress,
} from "../lib/constants";
import { useWalletStore } from "../components/providers/wallet-provider";

export const useContract = () => {
  const [contract, setContract] = React.useState<ethers.Contract>();
  const { signer } = useWalletStore();

  React.useEffect(() => {
    if (signer) {
      const contract = new ethers.Contract(
        l1BridgeAddress as string,
        L1Bridge.abi,
        signer
      );
      setContract(contract);
    }
  }, [signer]);

  // Send native currency to the other network.
  async function sendNative(receiver: string, value: string) {
    if (!contract) {
      console.error("Contract not found");
      return null;
    }
    try {
      const res = await contract.sendNative(receiver, {
        value: ethers.utils.parseEther(value),
      });
      console.log("🚀 ~ sendNative ~ res:", res);
      const receipt = await res.wait();
      console.log("🚀 ~ sendNative ~ receipt:", receipt);
      return receipt;
    } catch (error) {
      console.error("Error sending native currency:", error);
      throw error;
    }
  }

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
