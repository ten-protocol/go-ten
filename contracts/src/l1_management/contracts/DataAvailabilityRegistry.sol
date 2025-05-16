// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "../../common/Structs.sol";
import "../../cross_chain_messaging/L1/IMerkleTreeMessageBus.sol";
import "../interfaces/IDataAvailabilityRegistry.sol";
import "../interfaces/INetworkEnclaveRegistry.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "../../common/UnrenouncableOwnable2Step.sol";
/**
 * @title DataAvailabilityRegistry
 * @dev Contract for managing data availability and rollup validation
 * Implements a challenge period for state root disputes
 * Uses MerkleTreeMessageBus for message verification and value transfers
 */
contract DataAvailabilityRegistry is IDataAvailabilityRegistry, Initializable, UnrenouncableOwnable2Step {

    // RollupStorage: A storage structure to manage and organize MetaRollup instances in a mapping by their hash.
    struct RollupStorage {
        mapping(bytes32 rollupHash => MetaRollup rollup) byHash;
        uint256 nextFreeSequenceNumber;
    }

    /**
     * @dev Storage for rollups
     */
    RollupStorage private rollups;

    /**
     * @dev Last batch sequence number
     */
    uint256 public lastBatchSeqNo;

    /**
     * @dev Rollup challenge period for state root disputes
     */
    uint256 private challengePeriod;
    
    IMerkleTreeMessageBus public merkleMessageBus;
    INetworkEnclaveRegistry public enclaveRegistry;

    constructor() {
        _transferOwnership(msg.sender);
    }

    /**
     * @dev Initializes the contract with the owner and sets up the message bus
     * @param _merkleMessageBus Address of the MerkleTreeMessageBus
     * @param _enclaveRegistry Address of the NetworkEnclaveRegistry
     * @param _owner Address of the contract owner
     */
    function initialize(
        address _merkleMessageBus,
        address _enclaveRegistry,
        address _owner
    ) public initializer {
        __Ownable_init(_owner);
        merkleMessageBus = IMerkleTreeMessageBus(_merkleMessageBus);
        enclaveRegistry = INetworkEnclaveRegistry(_enclaveRegistry);
        lastBatchSeqNo = 0;
        challengePeriod = 0;
    }

    /**
     * @dev Appends a rollup to the registry
     * @param _r The rollup to append
     */
    function AppendRollup(MetaRollup calldata _r) internal {
        rollups.byHash[_r.Hash] = _r;

        if (_r.LastSequenceNumber > lastBatchSeqNo) {
            lastBatchSeqNo = _r.LastSequenceNumber;
        }
    }


    /**
     * @dev Modifier to verify the integrity of a rollup
     * @param r The rollup to verify
     */
    modifier verifyRollupIntegrity(MetaRollup calldata r) {
        // Block binding checks
        require(block.number > r.BlockBindingNumber, "Cannot bind to future or current block");
        require(block.number < (r.BlockBindingNumber + 255), "Block binding too old");

        bytes32 knownBlockHash = blockhash(r.BlockBindingNumber);

        require(knownBlockHash != 0x0, "Unknown block hash");
        require(knownBlockHash == r.BlockBindingHash, "Block binding mismatch");
        require(blobhash(0) != bytes32(0), "Blob hash is not set");

        bytes32 compositeHash = keccak256(abi.encodePacked(
            r.LastSequenceNumber,
            r.LastBatchHash,
            r.BlockBindingHash,
            r.BlockBindingNumber,
            r.crossChainRoot,
            blobhash(0)
        ));

        // Verify the enclave signature using the registry
        address enclaveID = ECDSA.recover(compositeHash, r.Signature);
        require(enclaveRegistry.isAttested(enclaveID), "enclaveID not attested");
        require(enclaveRegistry.isSequencer(enclaveID), "enclaveID not a sequencer");
        _;
    }
    /**
     * @dev Adds a rollup to the registry
     * @param r The rollup to add
     * 
     * TODO can we make it so only attested sequencer enclaves can call this? can pass the requester ID as a param?
     */
    function addRollup(MetaRollup calldata r) external verifyRollupIntegrity(r) {
        AppendRollup(r);

        if (r.crossChainRoot != bytes32(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff)) {
            uint256 activationTime = block.timestamp + challengePeriod;
            merkleMessageBus.addStateRoot(r.crossChainRoot, activationTime);
        }

        emit RollupAdded(blobhash(0), r.Signature);
    }

    /**
     * @dev Gets a rollup by hash
     * @param rollupHash The hash of the rollup to get
     * @return bool True if the rollup exists, false otherwise
     * @return Structs.MetaRollup The rollup
     */
    function getRollupByHash(bytes32 rollupHash) external view returns (bool, MetaRollup memory) {
        MetaRollup memory rol = rollups.byHash[rollupHash];
        return (rol.Hash == rollupHash , rol);
    }

    /**
     * @dev Gets the rollup challenge period
     * @return uint256 The challenge period
     */
    function getChallengePeriod() external view returns (uint256) {
        return challengePeriod;
    }

    /**
     * @dev Sets the challenge period
     * @param _delay The delay in seconds
     */
    function setChallengePeriod(uint256 _delay) external onlyOwner {
        challengePeriod = _delay;
    }
}