import { ethers } from "ethers";
import L1Bridge from "../../artifacts/IBridge.sol/IBridge.json";
import { l1Bridge as l1BridgeAddress } from "../lib/constants";

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
    return this.contract.sendNative(receiver, {
      value: ethers.utils.parseEther(value),
    });
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

  // Get the balance of an ERC20 asset.
  async getERC20Balance(asset: string) {
    const erc20 = new ethers.Contract(
      asset,
      ["function balanceOf(address)"],
      this.signer
    );
    return erc20.balanceOf(this.signer.getAddress());
  }

  // Get the allowance of an ERC20 asset.
  async getERC20Allowance(asset: string) {
    const erc20 = new ethers.Contract(
      asset,
      ["function allowance(address,address)"],
      this.signer
    );
    return erc20.allowance(this.signer.getAddress(), l1BridgeAddress);
  }

  // Approve an ERC20 asset.
  async approveERC20(asset: string, amount: string) {
    const erc20 = new ethers.Contract(
      asset,
      ["function approve(address,uint256)"],
      this.signer
    );
    return erc20.approve(l1BridgeAddress, amount);
  }

  // Get the transaction count of the signer.
  async getTransactionCount() {
    return this.signer.getTransactionCount();
  }
}
