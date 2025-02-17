import {EthereumProvider, HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { keccak256 } from 'ethers';
import { HardhatEthersProvider } from '@nomicfoundation/hardhat-ethers/internal/hardhat-ethers-provider';

function process_value_transfer(ethers, value_transfer) {
    const abiTypes = ['address', 'address', 'uint256', 'uint64'];
    const msg = [
      value_transfer['args'].sender, value_transfer['args'].receiver,
      value_transfer['args'].amount.toString(), value_transfer['args'].sequence.toString()
    ]

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
async function waitForRootPublished(management, msg, proof, root, provider: EthereumProvider, interval = 20000, timeout = 12000000) {
    var gas_estimate = null
    const l1Ethers = new HardhatEthersProvider(provider, "layer1")    

    console.log(`balance of management contract = ${await l1Ethers.getBalance(management.getAddress())}`)

    const startTime = Date.now();
    while (gas_estimate === null) {
        try {
            console.log(`Extracting native value from cross chain message for root ${root}`)
            const tx = await management.getFunction('ExtractNativeValue').populateTransaction(msg, proof, root, {} ) 
            gas_estimate = await l1Ethers.estimateGas(tx)
        } catch (error) {
            console.log(`Elapsed: ${Date.now() - startTime}ms Estimate gas threw error : ${error}`)
        }
        if (Date.now() - startTime >= timeout) {
            console.log(`Timed out waiting for the estimate gas to return`)
            break
        }
        await sleep(interval)
    }
    console.log(`Estimation took ${Date.now() - startTime} ms`)
    return gas_estimate
}

async function retryLoop(callback: ()=>Promise<boolean>, interval = 20000, timeout = 12000000) {
    const startTime = Date.now();
    while (true) {
      try {
        const result = await callback();
        if (result) {
          return result;
        }
      } catch (error) {
        console.log(`Error in retry loop: ${error}`);
      }
      await sleep(interval);
    }
}
    

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l2Network = hre; 
    const {deployer} = await hre.getNamedAccounts();

    const networkConfig : any = await hre.network.provider.request({method: 'net_config'});
    console.log(`Network config = ${JSON.stringify(networkConfig, null, 2)}`);

    const mgmtContractAddress = networkConfig.ManagementContractAddress;
    const messageBusAddress = networkConfig.MessageBusAddress;
    const l2MessageBusAddress = networkConfig.L2MessageBusAddress;

    var mbusBase = await hre.ethers.getContractAt("MessageBus", l2MessageBusAddress);
    const mbus = mbusBase.connect(await hre.ethers.provider.getSigner(deployer)); 
    const tx = await mbus.getFunction("sendValueToL2").send(deployer, 1000, { value: 1000});
    const receipt = await tx.wait()
    console.log(`003_simple_withdrawal: Cross Chain send ${tx.hash} receipt status = ${receipt.status}`);

    const block = await hre.ethers.provider.send('eth_getBlockByHash', [receipt.blockHash, true]);
    console.log(`Block received:       ${block.number.toString(10)}`);
  

    const value_transfer = mbus.interface.parseLog(receipt.logs[0]);
    const _processed_value_transfer = process_value_transfer(hre.ethers, value_transfer)
    const msg = _processed_value_transfer[0]
    const msgHash = _processed_value_transfer[1]
    const decoded = decode_base64(block.crossChainTree)

    console.log(`Decoded = ${JSON.stringify(decoded, null, 2)}`)
    console.log(`Getting cross chain proof for 'v' and msgHash = ${msgHash}`)

    const proof : any = await retryLoop(() => hre.ethers.provider.send('ten_getCrossChainProof', ['v', msgHash]))

    console.log(`Proof = ${JSON.stringify(proof, null, 2)}`)

    var managementContract = await hre.ethers.getContractAt("ManagementContract", mgmtContractAddress);

    const decoded_proof = hre.ethers.decodeRlp(proof.Proof)

    const estimation = await waitForRootPublished(managementContract, msg, decoded_proof, proof.Root, hre.companionNetworks.layer1.provider)
    console.log(`Estimation for native value extraction = ${estimation}`)
};


export default func;
func.tags = ['GasDebug'];
