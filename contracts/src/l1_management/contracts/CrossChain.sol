// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.28;

import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../../common/PausableWithRoles.sol";

import "../interfaces/ICrossChain.sol";
import * as MessageBus from "../../cross_chain_messaging/common/MessageBus.sol";
import * as MerkleTreeMessageBus from "../../cross_chain_messaging/L1/MerkleTreeMessageBus.sol";
import "../../common/UnrenouncableOwnable2Step.sol";
/**
 * @title CrossChain
 * @dev Contract managing cross-chain value transfers and message verification
 * Implements reentrancy protection. Uses MerkleTreeMessageBus for message 
 * verification and value transfers
 */
contract CrossChain is ICrossChain, Initializable, UnrenouncableOwnable2Step, ReentrancyGuardUpgradeable {
     /**
     * @dev Mapping to track spent withdrawals and prevent double-spending
     */
    mapping(bytes32 withdrawalHash => bool isWithdrawalSpent) public isWithdrawalSpent;
    
    /**
     * @dev Mapping to track saved bundles and prevent double-spending
     */
    mapping(bytes32 bundleHash => bool isBundleSaved) public isBundleSaved;

    MessageBus.IMessageBus public messageBus;
    MerkleTreeMessageBus.IMerkleTreeMessageBus public merkleMessageBus;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Initializes the contract with an owner and message bus
     * @param owner Address that will own the contract
     * @param _messageBus Address of the message bus contract
     */
    function initialize(address owner, address _messageBus) public initializer {
        require(_messageBus != address(0), "Invalid message bus address");
        require(owner != address(0), "Owner cannot be 0x0");
        __UnrenouncableOwnable2Step_init(owner);  // This will initialize OwnableUpgradeable and Ownable2StepUpgradeable
        __ReentrancyGuard_init();
        merkleMessageBus = MerkleTreeMessageBus.IMerkleTreeMessageBus(_messageBus);
        messageBus = MessageBus.IMessageBus(_messageBus);
    }

    /**
     * @dev Checks if a bundle of cross-chain messages is available
     * @param crossChainHashes Array of cross-chain message hashes to verify
     * @return bool True if the bundle is available
     */
    function isBundleAvailable(bytes[] memory crossChainHashes) external view returns (bool) {
        bytes32 bundleHash = bytes32(0);
        for(uint256 i = 0; i < crossChainHashes.length; i++) {
            bundleHash = keccak256(abi.encode(bundleHash, bytes32(crossChainHashes[i])));
        }
        return isBundleSaved[bundleHash];
    }
}
