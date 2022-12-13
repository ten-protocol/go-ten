import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l1Network = hre.companionNetworks.layer1;
    const l2Network = hre; 

    const l1Accounts = await l1Network.getNamedAccounts();
    const l2Accounts = await l2Network.getNamedAccounts();

    const layer2BridgeDeployment = await l2Network.deployments.get("ObscuroL2Bridge");
    const HOCDeployment = await l1Network.deployments.get("HOCERC20");
    const POCDeployment = await l1Network.deployments.get("POCERC20");

    const setResult = await l1Network.deployments.execute("ObscuroBridge", {
        from: l1Accounts.deployer, 
        log: true,
    }, "setRemoteBridgeAddress", layer2BridgeDeployment.address);
    if (setResult.status != 1) {
        console.error("Ops");
        throw Error("ops");
    }

    let hocResult = await l1Network.deployments.execute("ObscuroBridge", {
        from: l1Accounts.deployer, 
        log: true,
    }, "whitelistToken", HOCDeployment.address, "HOC", "HOC");

    if (hocResult.status != 1) {
        console.error("Ops");
        throw Error("ops");
    }

    const pocResult = await l1Network.deployments.execute("ObscuroBridge", {
        from: l1Accounts.deployer, 
        log: true,
    }, "whitelistToken", POCDeployment.address, "POC", "POC");
    
    if (pocResult.status != 1) {
        console.error("Ops");
        throw Error("ops");
    }

    const eventSignature = "LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)";
    const topic = hre.ethers.utils.id(eventSignature)
    let eventIface = new hre.ethers.utils.Interface([ `event ${eventSignature}`]);

    function getXChainMessages(result: Receipt) {
        const events = result.logs?.filter((x)=> { 
            return x.topics.find((t: string)=> t == topic) != undefined;
        });

        const messages = events!.map((event)=> {
            const decodedEvent = eventIface.parseLog({
                topics: event!.topics!,
                data: event!.data
            });
        
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

    const promises = messages.map(async (msg) => {
        return await l2Network.deployments.execute("CrossChainMessenger", {
            from: l2Accounts.deployer, 
            log: true,
        }, "relayMessage", msg)
    })

    const results = await Promise.all(promises);
    results.forEach(res=>{
        if (res.status != 1) {
            throw Error("Unable to relay messages...");
        } 
    });
};

export default func;
func.tags = ['Whitelist', 'Whitelist_deploy'];
func.dependencies = ['ObscuroL2Bridge'];
