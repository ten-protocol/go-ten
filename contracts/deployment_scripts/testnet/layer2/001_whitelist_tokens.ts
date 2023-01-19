import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l1Network = hre.companionNetworks.layer1;
    const l2Network = hre; 

    const l1Accounts = await l1Network.getNamedAccounts();
    const l2Accounts = await l2Network.getNamedAccounts();

    // Get the HOC POC layer 1 deployments.
    const HOCDeployment = await l1Network.deployments.get("HOCERC20");
    const POCDeployment = await l1Network.deployments.get("POCERC20");

    // Tell the bridge to whitelist the address of HOC token. This generates a cross chain message.
    let hocResult = await l1Network.deployments.execute("ObscuroBridge", {
        from: l1Accounts.deployer, 
        log: true,
    }, "whitelistToken", HOCDeployment.address, "HOC", "HOC");

    if (hocResult.status != 1) {
        console.error("Unable to whitelist HOC token!");
        throw Error("Unable to whitelist HOC token!");
    }

    // Tell the bridge to whitelist POC. This also generates a cross chain message.
    const pocResult = await l1Network.deployments.execute("ObscuroBridge", {
        from: l1Accounts.deployer, 
        log: true,
    }, "whitelistToken", POCDeployment.address, "POC", "POC");
    
    if (pocResult.status != 1) {
        console.error("Unable to whitelist POC token!");
        throw Error("Unable to whitelist POC token!");
    }

    const eventSignature = "LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)";
    // Get the hash id of the event signature
    const topic = hre.ethers.utils.id(eventSignature)

    // Get the interface for the event in order to convert it to cross chain message.
    let eventIface = new hre.ethers.utils.Interface([ `event ${eventSignature}`]);

    // This function converts the logs from transaction receipts into cross chain messages
    function getXChainMessages(result: Receipt) {
        
        // Get events with matching topic as the event id of LogMessagePublished
        const events = result.logs?.filter((x)=> { 
            return x.topics.find((t: string)=> t == topic) != undefined;
        });

        const messages = events!.map((event)=> {
            // Parse the rlp encoded log into event.
            const decodedEvent = eventIface.parseLog({
                topics: event!.topics!,
                data: event!.data
            });
        
            //Construct the cross chain message.
            const xchainMessage = {
                sender: decodedEvent.args[0],
                sequence: decodedEvent.args[1],
                nonce: decodedEvent.args[2],
                topic: decodedEvent.args[3],
                payload: decodedEvent.args[4],
                consistencyLevel: decodedEvent.args[5]
            };

            return xchainMessage;
        })

        return messages;
    }

    let messages = getXChainMessages(hocResult);
    messages = messages.concat(getXChainMessages(pocResult));

    // Freeze until the enclave processes the blocks and picks up the messages that have been carried over.
    await new Promise(resolve=>setTimeout(resolve, 2_000));
    
    // Perform message relay. This will forward the whitelist command to the L2 subordinate bridge.
    console.log(`Relaying messages using account ${l2Accounts.deployer}`);
    const relayMsg = async (msg: any) => {
        return l2Network.deployments.execute("CrossChainMessenger", {
            from: l2Accounts.deployer, 
            log: true,
            gasLimit: 1024_000_000
        }, "relayMessage", msg);
    };

    console.log(`Relaying message - ${JSON.stringify(messages[0])}`);
    let results = [await relayMsg(messages[0])];

    console.log(`Relaying message - ${JSON.stringify(messages[1])}`);
    results = results.concat(await relayMsg(messages[1]))

    results.forEach(res=>{
        if (res.status != 1) {
            throw Error("Unable to relay messages...");
        } 
    });
};

export default func;
func.tags = ['Whitelist', 'Whitelist_deploy'];
func.dependencies = ['EthereumBridge'];
