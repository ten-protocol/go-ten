import { EthereumProvider, HardhatRuntimeEnvironment } from 'hardhat/types';
import { DeployFunction } from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { keccak256 } from 'ethers';
import { HardhatEthersProvider } from '@nomicfoundation/hardhat-ethers/internal/hardhat-ethers-provider';

function process_value_transfer(ethers, value_transfer) {
    const abiTypes = ['address', 'address', 'uint256', 'uint64'];
    const msg = [
      value_transfer['args'].sender, value_transfer['args'].receiver,
      value_transfer['args'].amount.toString(), value_transfer['args'].sequence.toString()
    ];

    const abiCoder = ethers.AbiCoder.defaultAbiCoder();
    const encodedMsg = abiCoder.encode(abiTypes, msg);
    return [msg, ethers.keccak256(encodedMsg)];
}

function decode_base64(base64String: string) {
    let jsonString = atob(base64String);
    return JSON.parse(jsonString);
}

async function sleep(ms: number) {
    return new Promise((resolve) => {
      setTimeout(resolve, ms);
    });
}

async function buildMerkleProof(tree: any, message: any) {
    const proof = tree.getProof(message);
    return proof;
}

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l2Network = hre; 
    const l1Network = hre.companionNetworks.layer1!;
    const { deployer } = await hre.getNamedAccounts();
    
    // Start Generation Here
    // Retrieve network configuration
    const networkConfig: any = await hre.network.provider.request({ method: 'net_config' });
    console.log(`Network config = ${JSON.stringify(networkConfig, null, 2)}`);

    const l2Messenger = l2Network.deployments.get("CrossChainMessenger");
    const l1Messenger = l1Network.deployments.get("CrossChainMessenger");
};

async function bridgeTokenToL2(l2Messenger: any, l1Messenger: any, tokenAddress: string, amount: string) {

}

export default func;
func.tags = ['GasDebug', 'TokenWithdrawal'];