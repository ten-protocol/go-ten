// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "../messaging/IMerkleTreeMessageBus.sol";
import "../messaging/IMessageBus.sol";
import "../messaging/Structs.sol";
import "../messaging/messenger/ICrossChainMessenger.sol";
import "./ICrossChain.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

contract CrossChain is ICrossChain, Initializable, OwnableUpgradeable {
    bool private paused;
    uint256 private challengePeriod;
    IMessageBus public messageBus;
    MerkleTreeMessageBus.IMerkleTreeMessageBus public merkleMessageBus;
    mapping(bytes32 => bool) public isWithdrawalSpent;
    mapping(bytes32 =>bool) public isBundleSaved;

    event LogManagementContractCreated(address messageBusAddress);

    constructor() {
        _transferOwnership(msg.sender);
    }

    function initialize(address _messageBus) public initializer {
        __Ownable_init(msg.sender);
        messageBus = IMessageBus(_messageBus);
        paused = false; // Default to withdrawals enabled
        emit LogManagementContractCreated(address(messageBus));
    }

    function extractNativeValue(
        MessageStructs.Structs.ValueTransferMessage calldata msg,
        bytes32[] calldata proof,
        bytes32 root
    ) external {
        require(!paused, "withdrawals are paused");
        merkleMessageBus.verifyValueTransferInclusion(_msg, proof, root);
        bytes32 msgHash = keccak256(abi.encode(_msg));
        require(isWithdrawalSpent[msgHash] == false, "withdrawal already spent");
        isWithdrawalSpent[keccak256(abi.encode(_msg))] = true;

        messageBus.receiveValueFromL2(_msg.receiver, _msg.amount);
    }

    // Testnet function to allow the contract owner to retrieve **all** funds from the network bridge.
    function RetrieveAllBridgeFunds() external onlyOwner {
        messageBus.retrieveAllFunds(msg.sender);
    }

    function pauseWithdrawals(bool _pause) external onlyOwner {
        paused = _pause;
        emit WithdrawalsPaused(_pause);
    }

    function isBundleAvailable(bytes[] memory crossChainHashes) external view returns (bool) {
        bytes32 bundleHash = bytes32(0);
        for(uint256 i = 0; i < crossChainHashes.length; i++) {
            bundleHash = keccak256(abi.encode(bundleHash, bytes32(crossChainHashes[i])));
        }
        return isBundleSaved[bundleHash];
    }

    function getChallengePeriod() external view returns (uint256) {
        return challengePeriod;
    }

    function setChallengePeriod(uint256 _delay) external onlyOwner {
        challengePeriod = _delay;
    }
}
