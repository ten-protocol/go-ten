// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "../messaging/IMerkleTreeMessageBus.sol";
import "../messaging/MessageBus.sol";
import "./ICrossChain.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";

contract CrossChain is ICrossChain, Initializable, OwnableUpgradeable, ReentrancyGuardUpgradeable {
    bool private paused;
    uint256 private challengePeriod;
    IMessageBus public messageBus;
    IMerkleTreeMessageBus public merkleMessageBus;
    mapping(bytes32 => bool) public isWithdrawalSpent;
    mapping(bytes32 =>bool) public isBundleSaved;

    event LogManagementContractCreated(address messageBusAddress);

    constructor() {
        _transferOwnership(msg.sender);
    }

    function initialize(address _messageBus) public initializer {
        __Ownable_init(msg.sender);
        __ReentrancyGuard_init();
        messageBus = IMessageBus(_messageBus);
        paused = false; // Default to withdrawals enabled
        emit LogManagementContractCreated(address(messageBus));
    }

    function extractNativeValue(
        Structs.ValueTransferMessage calldata _msg,
        bytes32[] calldata proof,
        bytes32 root
    ) external nonReentrant {
        require(!paused, "withdrawals are paused");
        merkleMessageBus.verifyValueTransferInclusion(_msg, proof, root);
        bytes32 msgHash = keccak256(abi.encode(_msg));
        require(isWithdrawalSpent[msgHash] == false, "withdrawal already spent");
        isWithdrawalSpent[msgHash] = true;  // Use stored msgHash

        messageBus.receiveValueFromL2(_msg.receiver, _msg.amount);
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
