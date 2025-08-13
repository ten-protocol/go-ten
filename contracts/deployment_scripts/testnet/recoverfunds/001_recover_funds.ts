// Requires: npm install axios
import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import {TenBridgeTestnet} from "../../../typechain-types";
import axios from 'axios';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployer} = await hre.getNamedAccounts();

    const rpcUrl = 'http://erpc.dev-testnet.ten.xyz:80';
    const rpcPayload = {
        jsonrpc: '2.0',
        method: 'ten_config',
        params: [],
        id: 1
    };
    const response = await axios.post(rpcUrl, rpcPayload, {
        headers: { 'Content-Type': 'application/json' }
    });

    // recover funds from the new bridge contract
    const bridgeContractAddress = response.data.result.L1Bridge;
    const bridgeContract = (await hre.ethers.getContractFactory('TenBridgeTestnet')).attach(bridgeContractAddress) as TenBridgeTestnet;

    // Add debugging before the transaction
    console.log("Deployer address:", deployer);
    console.log("Bridge contract address:", bridgeContractAddress);

    // Check if deployer exists
    if (!deployer) {
        console.log("ERROR: No deployer account found");
        return;
    }

    // Check if deployer has ADMIN_ROLE
    const ADMIN_ROLE = await bridgeContract.ADMIN_ROLE();
    const hasAdminRole = await bridgeContract.hasRole(ADMIN_ROLE, deployer);
    console.log("Deployer has ADMIN_ROLE:", hasAdminRole);

    if (!hasAdminRole) {
        console.log("ERROR: Deployer does not have ADMIN_ROLE on the bridge contract");
        return;
    }
    // recover all bridged funds to the gnosis safe
    const tx = await bridgeContract.recoverTestnetFunds("0x9f7b0CDB121Af3923A98771c326b1aAC03A0D717");

    const receipt = await tx.wait();
    if (receipt && receipt.status === 1) {
        console.log("Successfully recovered funds from the bridge.");
    } else {
        console.log("Recovery transaction failed");
    }
};

export default func;
// No dependencies