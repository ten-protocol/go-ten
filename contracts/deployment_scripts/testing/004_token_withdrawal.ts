import { EthereumProvider, HardhatRuntimeEnvironment } from 'hardhat/types';
import { DeployFunction } from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { ContractTransactionReceipt, decodeBase64, keccak256, Log } from 'ethers';
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

    const token = await l1Network.deployments.get("HOCERC20");

    await bridgeTokenToL2(hre, "1000");
    await withdrawTokenFromL2(hre, token.address, "1000");
};

async function bridgeTokenToL2(hre: HardhatRuntimeEnvironment, amount: string) {
    const l1Network = hre.companionNetworks.layer1!;
    const deployerL1 = (await l1Network.getNamedAccounts()).deployer;
    const deployerL2 = (await hre.getNamedAccounts()).deployer;
    const token = await l1Network.deployments.get("HOCERC20");
    console.log(`Token address = ${token.address}`);
    const tenBridge = await l1Network.deployments.get("TenBridge");
    
    const l1Provider = new HardhatEthersProvider(l1Network.provider, "layer1")
    const signer = await l1Provider.getSigner(deployerL1);
    const tokenContract = await hre.ethers.getContractAt("ERC20", token.address, signer);
    {
        const tx = await tokenContract.connect(signer).approve(tenBridge.address, amount);
        const receipt = await tx.wait();
        if (receipt!.status !== 1) {
            throw new Error("Token approval failed");
        }
        console.log(`Token approval successful for l1 to l2`);
    }
    {
        const tenBridgeContract = await hre.ethers.getContractAt("TenBridge", tenBridge.address, signer);
        const bridgeTx = await tenBridgeContract.sendERC20(token.address, amount, deployerL2!);
        const receipt = await bridgeTx.wait();
        if (receipt!.status !== 1) {
            throw new Error("Token bridge failed");
        }
        console.log(`Layer1.SendERC20 successful`);
        await finalizeTokenBridgeToL2(hre, receipt!, token.address);
    }
}

async function getMessage(hre: HardhatRuntimeEnvironment, receipt: ContractTransactionReceipt) {
    async function submitMessagesFromTx(hre: HardhatRuntimeEnvironment, receipt: ContractTransactionReceipt) {

        const eventSignature = "LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)";
        const topic = hre.ethers.id(eventSignature)
        let eventIface = new hre.ethers.Interface([ `event LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)`]);
        console.log(`Event topic = ${topic}`);

        const events = receipt.logs!.filter((x)=> { 
          return x.topics.find((t: string)=> t == topic) != undefined;
        }) || [];
  
        if (events.length == 0) {
          throw new Error("No messages found");
        }
       
        const promises = events.map(async (event) => {
            const decodedEvent = eventIface.parseLog({
              topics: event!.topics!,
              data: event!.data
            })!!;
                
            const xchainMessage = {
              sender: decodedEvent.args[0],
              sequence: decodedEvent.args[1],
              nonce: decodedEvent.args[2],
              topic: decodedEvent.args[3],
              payload: decodedEvent.args[4],
              consistencyLevel: decodedEvent.args[5]
            };
            return xchainMessage;
        });
        return await Promise.all(promises);
    }
    return await submitMessagesFromTx(hre, receipt);
}

async function finalizeTokenBridgeToL2(hre: HardhatRuntimeEnvironment, receipt: ContractTransactionReceipt, tokenAddress: string) {
    await sleep(5000);
    
    const l2Network = hre;
    const { deployer } = await hre.getNamedAccounts();
    const signer = await hre.ethers.getSigner(deployer!);

    const l2Messenger = await l2Network.deployments.get("CrossChainMessenger");
    const l2MessengerContract = await hre.ethers.getContractAt("CrossChainMessenger", l2Messenger.address, signer);
    const messages = await getMessage(hre, receipt);
    const tx = await l2MessengerContract.relayMessage(messages![0]!);
    const receiptForRelay = await tx.wait();
    if (receiptForRelay!.status !== 1) {
        throw new Error("Relay message failed");
    }
    const bridgeDeployment = await l2Network.deployments.get("EthereumBridge")
    const bridge = await hre.ethers.getContractAt("EthereumBridge", bridgeDeployment.address, signer)
    const correspondingToken = await bridge.remoteToLocalToken(tokenAddress);
    if (!correspondingToken) {
        throw new Error("No corresponding wrapped token found");
    }

    console.log(`Corresponding token = ${correspondingToken}`);

    const tokenContract = await hre.ethers.getContractAt("WrappedERC20", correspondingToken, signer);
    const myAddress = await signer.getAddress();
    console.log(`My address = ${myAddress}`);
    const balance = await tokenContract.balanceOf(myAddress);
    console.log(`Balance = ${balance}`);
}

async function withdrawTokenFromL2(hre: HardhatRuntimeEnvironment, tokenAddress: string, amount: string) {
    const l2Network = hre;
    const l1Network = hre.companionNetworks.layer1!;
    const { deployer } = await hre.getNamedAccounts();
    const deployerL1 = (await l1Network.getNamedAccounts()).deployer;
    const signer = await hre.ethers.getSigner(deployer!);

    const bridgeDeployment = await l2Network.deployments.get("EthereumBridge")
    const bridge = await hre.ethers.getContractAt("EthereumBridge", bridgeDeployment.address, signer)


    const wrappedToken = await bridge.remoteToLocalToken(tokenAddress);
    const wrappedTokenContract = await hre.ethers.getContractAt("WrappedERC20", wrappedToken, signer);
    {
        const tx = await wrappedTokenContract.approve(bridgeDeployment.address, amount);
        const receipt = await tx.wait();
        if (receipt!.status !== 1) {
            throw new Error("Token approval failed");
        }
    }
    {
        const tx = await bridge.sendERC20(wrappedToken, amount, deployerL1!);
        const receipt = await tx.wait();
        console.log(`Receipt = ${JSON.stringify(receipt, null, 2)}`);
        if (receipt!.status !== 1) {
        throw new Error("Token withdrawal failed");
        }
        await finalizeTokenWithdrawalFromL2(hre, receipt!, tokenAddress);
    }
}

async function finalizeTokenWithdrawalFromL2(hre: HardhatRuntimeEnvironment, receipt: ContractTransactionReceipt, tokenAddress: string) {
    const l1Network = hre.companionNetworks.layer1!;
    const l1Provider = new HardhatEthersProvider(l1Network.provider, "layer1")
    const deployerL1 = (await l1Network.getNamedAccounts()).deployer;
    const signer = await l1Provider.getSigner(deployerL1!);
    
    const block = await hre.ethers.provider.send("eth_getBlockByHash", [receipt.blockHash, true]);

    const decoded = decode_base64(block["crossChainTree"]);
    console.log(`Tree = ${JSON.stringify(decoded, null, 2)}`);
    
    const messages = await getMessage(hre, receipt);

    await sleep(5000);

    const l1Messenger = await l1Network.deployments.get("CrossChainMessenger");
    const l1MessengerContract = await hre.ethers.getContractAt("CrossChainMessenger", l1Messenger.address, signer);
   
    const tree = StandardMerkleTree.of(decoded, ["string", "bytes32"]);
    const proof = tree.getProof(0)
   
    console.log(`Root calculated = ${tree.root}; block root = ${block["crossChainTreeHash"]}`)
    console.log(`Proof = ${JSON.stringify(proof, null, 2)}`);

    const message = messages![0]!;
    // Finally do a normal ABI encode
    const encodedMessage = hre.ethers.AbiCoder.defaultAbiCoder().encode(
        ["address", "uint64", "uint32", "uint32", "bytes", "uint8"],
        [
            message.sender,
            message.sequence,
            message.nonce, 
            message.topic,
            message.payload,
            message.consistencyLevel
        ]
    );

    const hashedMessage2 = hre.ethers.keccak256(encodedMessage);
    console.log(`Hashed message = ${hashedMessage2}`);

    // Encode packed "m" and the hashed message, then hash it
    const packedLeaf = hre.ethers.AbiCoder.defaultAbiCoder().encode(
        ["string", "bytes32"],
        ["m", hashedMessage2]
    );
    const hashedLeaf = hre.ethers.keccak256(packedLeaf);
    const packedLeaf2 = hre.ethers.AbiCoder.defaultAbiCoder().encode(
        ["bytes32"],
        [hashedLeaf]
    );
    const hashedLeaf2 = hre.ethers.keccak256(packedLeaf2);
    console.log(`Hashed leaf = ${hashedLeaf2}`);
    console.log(`Tree root = ${tree.root}`);
    console.log(`Leaf matches root: ${hashedLeaf2 === tree.root}`);

    while (true) {
        try {
            const tx = await l1MessengerContract.relayMessageWithProof(message, proof, tree.root);
            const receiptForRelay = await tx.wait();
            if (receiptForRelay!.status !== 1) {
                break;
            } 
            break;
        } catch (e) {
            console.log(`Error = ${e}`);
            await sleep(5000);
        }
    }
}

export default func;
func.tags = ['GasDebug', 'TokenWithdrawal'];