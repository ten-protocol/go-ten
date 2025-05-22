// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.28;

import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "../interfaces/ICrossChain.sol";
import * as MessageBus from "../../cross_chain_messaging/common/MessageBus.sol";
import * as MerkleTreeMessageBus from "../../cross_chain_messaging/L1/MerkleTreeMessageBus.sol";
import "../../common/UnrenouncableOwnable2Step.sol";
/**
 * @title CrossChain
 * @dev Contract managing cross-chain value transfers and message verification
 * Implements reentrancy protection and pausable withdrawals for security
 * Uses MerkleTreeMessageBus for message verification and value transfers
 */
contract CrossChain is ICrossChain, Initializable, UnrenouncableOwnable2Step, ReentrancyGuardUpgradeable {

    /**
     * @dev Flag to control withdrawal functionality
     */
    bool private paused;

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
     * @dev Initializes the contract with an owner
     * @param owner Address that will own the contract
     */
    function initialize(address owner) public initializer {
        __UnrenouncableOwnable2Step_init(owner);
        __ReentrancyGuard_init();
        merkleMessageBus = new MerkleTreeMessageBus.MerkleTreeMessageBus();
        merkleMessageBus.initialize(owner, address(this));
        messageBus = MessageBus.IMessageBus(address(merkleMessageBus));
        paused = false; // Default to withdrawals enabled
    }

    /**
     * @dev Pauses or resumes withdrawals
     * @param _pause True to pause withdrawals, false to resume
     */
    function pauseWithdrawals(bool _pause) external onlyOwner {
        paused = _pause;
        emit WithdrawalsPaused(_pause);
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
