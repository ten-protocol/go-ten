import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    // Get the deployer account which should be prefunded with enough funds in order to create smart contracts
    const {deployer} = await getNamedAccounts();

    // We pull out the CrossChainMessenger which should be deployed as a core contract.
    const messengerDeployment = await deployments.get("CrossChainMessenger");

    // We deploy the obscuro bridge and give it the address of the cross chain messenger
    // which will be used in order to enable to cross chain functionality of sending/receiving assets.
    await deployments.deploy('ObscuroBridge', {
        from: deployer,
        args: [ messengerDeployment.address ],
        log: true,
    });
};

export default func;
func.tags = ['ObscuroBridge', 'ObscuroBridge_deploy'];
// This script will always be deployed after the script deploying the CrossChainMessenger
func.dependencies = ['CrossChainMessenger'];
