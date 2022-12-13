import '@nomiclabs/hardhat-ethers';
import {ethers} from "hardhat";
import { ObscuroBridge } from '../typechain-types/contracts/bridge/L1/L1_Bridge.sol';
import { ManagementContract } from '../typechain-types/contracts/management';
import { CrossChainMessenger } from '../typechain-types/contracts/messaging/messenger';
import * as hre from "hardhat";

interface L1DeploymentResult {
  management: ManagementContract,
  messenger: CrossChainMessenger,
  bridge: ObscuroBridge,
}

async function deployL1() {
    const ManagementContract = await ethers.getContractFactory("ManagementContract");

    const Messenger = await ethers.getContractFactory("CrossChainMessenger");
    const L1Bridge = await ethers.getContractFactory("ObscuroBridge");

    const mgmtContract : ManagementContract = await ManagementContract.deploy();
    const messenger : CrossChainMessenger = await Messenger.deploy(await mgmtContract.messageBus());
    const bridge : ObscuroBridge = await L1Bridge.deploy(messenger.address);

    return {
      management : mgmtContract,
      messenger : messenger,
      bridge: bridge
    }
}

async function deployL2(l1Result: L1DeploymentResult) {
    const L2Bridge = await ethers.getContractFactory("ObscuroL2Bridge");
    const Messenger = await ethers.getContractFactory("CrossChainMessenger");

    await Messenger.deploy()
}

async function finalizeL1(L2DeploymentREsult: {}) {

}

// @ts-ignore
async function sayHello(hello){
  //5. It will execute a function inside the contract, the `hello()` function literally the function that we create on our smart contract (requested out of the network)
  console.log("Say hello:", await hello.hello())
}

deployL1().then(deployL2)