import { ethers } from "ethers";
import ImageGuessGameJson from "@/assets/contract/artifacts/contracts/ImageGuessGame.sol/ImageGuessGame.json";
import ContractAddress from "@/assets/contract/address.json";
import { handleMetaMaskError, bigNumberToNumber } from "./utils";

export default class Web3listener {
  contract: ethers.Contract;

  constructor(signer: any) {
    this.contract = new ethers.Contract(
      ContractAddress.address,
      ImageGuessGameJson.abi,
      signer
    );
    // ElNotification({
    //   message: `[ImageGuessGame Contract] Contract Address: ${ContractAddress.address}`,
    //   type: 'success'
    // })
  }
}
