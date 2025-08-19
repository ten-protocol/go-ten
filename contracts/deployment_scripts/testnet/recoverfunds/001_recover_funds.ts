// Requires: npm install axios
import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import {TenBridgeTestnet} from "../../../typechain-types";
import axios from 'axios';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { deployer } = await hre.getNamedAccounts();
    const { ethers } = hre;

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

    // --- Original checks ---
    if (!deployer) {
        console.log("ERROR: No deployer account found");
        return;
    }

    // ensure we attach with the DEPLOYER signer
    const signer = await ethers.getSigner(deployer);
    const bridgeContract = (await hre.ethers.getContractFactory('TenBridgeTestnet', signer))
        .attach(bridgeContractAddress) as TenBridgeTestnet;

    // --- Diagnostics block (kept lightweight; script still executes) ---
    console.log("Deployer address (intended):", deployer);
    console.log("Actual tx signer:", await signer.getAddress());
    console.log("Bridge contract address:", bridgeContractAddress);

    // Chain ID sanity
    const net = await ethers.provider.getNetwork();
    console.log("Provider chainId:", net.chainId.toString());

    // Is there code at the proxy?
    const proxyCode = await ethers.provider.getCode(bridgeContractAddress);
    console.log("Proxy has runtime code:", proxyCode !== '0x', "len:", proxyCode.length);

    // EIP-1967 implementation slot: bytes32(uint256(keccak256('eip1967.proxy.implementation')) - 1)
    try {
        const implSlot =
            BigInt(ethers.keccak256(ethers.toUtf8Bytes('eip1967.proxy.implementation'))) - 1n;
        const implSlotHex = ethers.toBeHex(implSlot, 32);

        // read storage via raw RPC (works across providers)
        const rawImpl: string = await ethers.provider.send('eth_getStorageAt', [
            bridgeContractAddress,
            implSlotHex,
            'latest'
        ]);
        const implAddress = ethers.getAddress('0x' + rawImpl.slice(-40));
        console.log("Implementation address (EIP-1967):", implAddress);

        const implCode = await ethers.provider.getCode(implAddress);
        console.log("Impl has runtime code:", implCode !== '0x', "len:", implCode.length);

        // Check the function selector is present in impl bytecode (heuristic)
        const selector = ethers.id("recoverTestnetFunds(address)").slice(0, 10);
        const selectorInImpl = implCode.toLowerCase().includes(selector.slice(2).toLowerCase());
        console.log("Selector", selector, "present in impl bytecode:", selectorInImpl);
    } catch (e) {
        console.log("EIP-1967 impl-slot read failed (non-proxy or different pattern?):", (e as any)?.message ?? e);
    }

    // Bridge balance (helps to know if value transfer will actually happen)
    const bridgeBal = await ethers.provider.getBalance(bridgeContractAddress);
    console.log("Bridge balance (wei):", bridgeBal.toString());

    const ADMIN_ROLE = await bridgeContract.ADMIN_ROLE();
    const hasAdminRole = await bridgeContract.hasRole(ADMIN_ROLE, deployer);
    console.log("Deployer has ADMIN_ROLE:", hasAdminRole);
    if (!hasAdminRole) {
        console.log("ERROR: Deployer does not have ADMIN_ROLE on the bridge contract");
        return;
    }

    // Optional preflight: static call to surface a revert reason (won't stop execution)
    try {
        // ethers v6: use .staticCall on the function
        // @ts-ignore (typechain may not declare .staticCall yet)
        await (bridgeContract.recoverTestnetFunds as any).staticCall("0x9f7b0CDB121Af3923A98771c326b1aAC03A0D717");
        console.log("Preflight staticCall succeeded (no immediate revert).");
    } catch (e) {
        const err = e as any;
        console.log("Preflight staticCall reverted (still proceeding):", err?.shortMessage || err?.reason || err?.message || e);
    }

    // --- Execute recovery (unchanged intention; now guaranteed to use deployer signer) ---
    try {
        const tx = await bridgeContract.recoverTestnetFunds("0x9f7b0CDB121Af3923A98771c326b1aAC03A0D717");
        console.log("Sent tx:", tx.hash);
        const receipt = await tx.wait();
        if (receipt && receipt.status === 1) {
            console.log("Successfully recovered funds from the bridge.");
        } else {
            console.log("Recovery transaction failed (receipt status not 1).");
        }
    } catch (e) {
        const err = e as any;
        console.log("Recovery tx threw:", err?.shortMessage || err?.reason || err?.message || e);
        // Sometimes estimateGas hides the reason; this may help:
        if (err?.data) console.log("Revert data:", err.data);
    }
};

export default func;
// No dependencies