import { ethers } from "ethers";
import L1Bridge from "../../artifacts/IBridge.sol/IBridge.json";
import {
  l1Bridge as l1BridgeAddress,
  messageBusAddress,
} from "../lib/constants";

interface IWeb3Service {
  contract: ethers.Contract;
  signer: any;
}

export default class Web3Service implements IWeb3Service {
  contract: ethers.Contract;
  signer: any;

  constructor(signer: any) {
    this.contract = new ethers.Contract(
      l1BridgeAddress as string,
      L1Bridge.abi,
      signer
    );
    this.signer = signer;
  }

  // Send native currency to the other network.
  async sendNative(receiver: string, value: string) {
    if (!this.contract) {
      console.error("Contract not found");
      return null;
    }
    try {
      const res = await this.contract.sendNative(receiver, {
        value: ethers.utils.parseEther(value),
      });
      console.log("🚀 ~ Web3Service ~ sendNative ~ res:", res);
      const receipt = await res.wait();
      return receipt;
    } catch (error) {
      console.error("Error sending native currency:", error);
      return null;
    }
  }

  // Send ERC20 assets to the other network.
  async sendERC20(
    receiver: string,
    amount: string,
    tokenContractAddress: string
  ) {
    return this.contract.sendERC20(tokenContractAddress, amount, receiver);
  }

  async getNativeBalance(provider: any, walletAddress: string) {
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
      return null;
    }
  }

  // balance of an ERC20 asset.
  async getTokenBalance(
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
      return null;
    }
  }

  //get bridge transactions
  async getBridgeTransactions(provider: any, userAddress: string) {
    try {
      if (!provider) {
        console.error("Provider not found");
        return null;
      }
      if (!userAddress) {
        console.error("User address not found");
        return null;
      }
      console.log(
        "🚀 ~ Web3Service ~ getBridgeTransactions ~ adddress:",
        ethers.utils.hexZeroPad(userAddress, 32)
      );
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
    }
  }
}
