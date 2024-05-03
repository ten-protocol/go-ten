import { ethers } from "ethers";
import { BridgeInterface } from "../../contracts/IBridge.sol";
import { l1Bridge as l1BridgeAddress } from "../lib/constants";

// // SPDX-License-Identifier: Apache 2

// pragma solidity >=0.7.0 <0.9.0;

// // The ERC20 token bridge interface.
// // Calling functions on it will result in assets being bridged over to the other layer automatically.
// interface BridgeInterface {
//     enum Topics {
//         TRANSFER,
//         MANAGEMENT,
//         VALUE
//     }

//     // Sends the native currency to the other layer. On Layer 1 the native currency is ETH, while on Layer 2 it is OBX.
//     // When it arrives on the other side it will be wrapped as a token.
//     // receiver - the L2 address that will receive the assets on the other network.
//     function sendNative(address receiver) external payable;

//     // Sends ERC20 assets over to the other network. The user must grant allowance to the bridge
//     // before calling this function for more or equal to the amount being bridged over.
//     // This can be done using IERC20(asset).approve(bridge, amount);
//     // asset - the address of the smart contract of the ERC20 token.
//     // amount - the number of tokens being transferred.
//     // receiver - the L2 address receiving the assets.
//     function sendERC20(
//         address asset,
//         uint256 amount,
//         address receiver
//     ) external;

//     // This function is called to retrieve assets that have been sent on the other layer.
//     // In the basic implementation it is only callable from the CrossChainMessenger when a message is
//     // being relayed.
//     function receiveAssets(
//         address asset,
//         uint256 amount,
//         address receiver
//     ) external;

//     struct ValueTransfer {
//         uint256 amount;
//         address recipient;
//     }
// }

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
      BridgeInterface,
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
