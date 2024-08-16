import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { keccak256 } from 'ethers';

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


  function decode_base64(base64String) {
    let jsonString = atob(base64String);
    return JSON.parse(jsonString);
  }
  

async function sleep(ms: number) {
    return new Promise((resolve) => {
      setTimeout(resolve, ms);
    });
}
async function waitForRootPublished(management, msg, proof, root, interval = 5000, timeout = 90000) {
    var gas_estimate = null
    const startTime = Date.now();
    while (gas_estimate === null) {
        try {
            gas_estimate = await management.estimateGas.ExtractNativeValue(msg, proof, root, {} )
        } catch (error) {
            console.log(`Estimate gas threw error : ${error.reason}`)
        }
        if (Date.now() - startTime >= timeout) {
            console.log(`Timed out waiting for the estimate gas to return`)
            break
        }
        await sleep(1_000)
    }
    return gas_estimate
}
    

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l2Network = hre; 
    const {deployer} = await hre.getNamedAccounts();
    
    var mbusBase = await hre.ethers.getContractAt("MessageBus", "0x526c84529b2b8c11f57d93d3f5537aca3aecef9b");
    const mbus = mbusBase.connect(await hre.ethers.provider.getSigner(deployer)); 
    const tx = await mbus.getFunction("sendValueToL2").send(deployer, 1000, { value: 1000});
    const receipt = await tx.wait()
    console.log(`003_simple_withdrawal: Cross Chain send receipt status = ${receipt.status}`);

    const block = await hre.ethers.provider.send('eth_getBlockByHash', [receipt.blockHash, true]);
    console.log(`Block received:       ${block.number}`)
  

    const value_transfer = mbus.interface.parseLog(receipt.logs[0]);
    const _processed_value_transfer = process_value_transfer(hre.ethers, value_transfer)
    const msg = _processed_value_transfer[0]
    const msgHash = _processed_value_transfer[1]
    const decoded = decode_base64(block.crossChainTree)

    console.log(`  Sender:        ${value_transfer['args'].sender}`)
    console.log(`  Receiver:      ${value_transfer['args'].receiver}`)
    console.log(`  Amount:        ${value_transfer['args'].amount}`)
    console.log(`  Sequence:      ${value_transfer['args'].sequence}`)
    console.log(`  VTrans Hash:   ${msgHash}`)
    console.log(`  XChain tree:   ${decoded}`)
    
    if (decoded[0][1] != msgHash) {
        console.error('Value transfer hash is not in the xchain tree!');
        return;
    }

    const tree = StandardMerkleTree.of(decoded, ["string", "bytes32"]);
    const proof = tree.getProof(['v',msgHash])
    console.log(`  Merkle root:   ${tree.root}`)
    console.log(`  Merkle proof:  ${JSON.stringify(proof, null,2)}`)
  
    if (block.crossChainTreeHash != tree.root) {
      console.error('Constructed merkle root matches block crossChainTreeHash');
      return
    }

    var managementContract = await hre.ethers.getContractAt("ManagementContract", "0x526c84529b2b8c11f57d93d3f5537aca3aecef9b");
  
};


export default func;
func.tags = ['GasDebug'];
