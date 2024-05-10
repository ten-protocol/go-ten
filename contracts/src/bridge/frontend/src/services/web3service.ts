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
    const res = await this.contract.sendNative(receiver, {
      value: ethers.utils.parseEther(value),
    });
    console.log("🚀 ~ Web3Service ~ sendNative ~ res", res);
    const receipt = await res.wait();
    console.log("🚀 ~ Web3Service ~ sendNative ~ receipt", receipt);
    return receipt;
  }

  // Send ERC20 assets to the other network.
  async sendERC20(asset: string, amount: string, receiver: string) {
    return this.contract.sendERC20(asset, amount, receiver);
  }

  // Receive assets that have been sent on the other network.
  async receiveAssets(asset: string, amount: string, receiver: string) {
    return this.contract.receiveAssets(asset, amount, receiver);
  }

  // Get the balance of the signer.
  async getBalance() {
    return this.signer.getBalance();
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

  // Get the balance of an ERC20 asset.
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
      // a new instance of the ERC20 contract
      const tokenContract = new ethers.Contract(
        tokenAddress,
        [
          "function balanceOf(address owner) view returns (uint256)",
          "function decimals() view returns (uint8)",
        ],
        p
      );
      console.log("🚀 ~ Web3Service ~ tokenContract:", tokenContract);

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

  // Get the allowance of an ERC20 asset.
  async getERC20Allowance(asset: string) {}

  // Approve an ERC20 asset.
  async approveERC20(asset: string, amount: string) {}

  // Get the transaction count of the signer.
  async getTransactionCount() {
    return this.signer.getTransactionCount();
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

  // async getBridgeTransactions(provider: any, userAddress: string) {
  //   console.log(
  //     "🚀 ~ Web3Service ~ getBridgeTransactions ~ messageBusAddress:",
  //     messageBusAddress
  //   );
  //   try {
  //     // Convert the provider to Web3Provider
  //     const web3Provider = new ethers.providers.Web3Provider(provider);

  //     // Specify the topics for the logs
  //     const topics = [
  //       ethers.utils.id("ValueTransfer(address,address,uint256,uint256)"), // Event signature
  //       null, // Placeholder for user's address topic
  //       ethers.utils.hexZeroPad(messageBusAddress, 32), // Message bus address as a topic
  //     ];

  //     // Set the user's address as the second topic
  //     topics[1] = ethers.utils.hexZeroPad(userAddress, 32);

  //     // Filter logs based on the specified topics
  //     const filter = {
  //       address: messageBusAddress, // Contract address where the event occurred
  //       topics: topics.filter((topic) => topic !== null), // Remove null topics
  //     };

  //     // Get logs based on the filter
  //     const transactions = await web3Provider.getLogs(filter);

  //     // Process and display transactions
  //     console.log("User's transactions:", transactions);
  //   } catch (error) {
  //     console.error("Error fetching user's transactions:", error);
  //   }
  // }
}
