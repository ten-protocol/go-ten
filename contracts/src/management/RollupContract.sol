// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "./INetworkEnclaveRegistry.sol";
import "./IRollupContract.sol";
import "../messaging/IMerkleTreeMessageBus.sol";
import "../common/Structs.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract RollupContract is IRollupContract, Initializable, OwnableUpgradeable {
    Structs.RollupStorage private rollups;
    uint256 public lastBatchSeqNo;

    IMerkleTreeMessageBus public merkleMessageBus;
    INetworkEnclaveRegistry public enclaveRegistry;

    constructor() {
        _transferOwnership(msg.sender);
    }

    function initialize(
        address _merkleMessageBus,
        address _enclaveRegistry,
        address _owner
    ) public initializer {
        __Ownable_init(_owner);
        merkleMessageBus = IMerkleTreeMessageBus(_merkleMessageBus);
        enclaveRegistry = INetworkEnclaveRegistry(_enclaveRegistry);
        lastBatchSeqNo = 0;
    }

    function AppendRollup(Structs.MetaRollup calldata _r) internal {
        rollups.byHash[_r.Hash] = _r;

        if (_r.LastSequenceNumber > lastBatchSeqNo) {
            lastBatchSeqNo = _r.LastSequenceNumber;
        }
    }


    modifier verifyRollupIntegrity(Structs.MetaRollup calldata r) {
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

    // TODO can we make it so only attested sequencer enclaves can call this? can pass the requester ID as a param?
    function addRollup(Structs.MetaRollup calldata r) external verifyRollupIntegrity(r) {
        AppendRollup(r);

        if (r.crossChainRoot != bytes32(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff)) {
            merkleMessageBus.addStateRoot(r.crossChainRoot, block.timestamp);
        }

        emit RollupAdded(blobhash(0), r.Signature);
    }


    function getRollupByHash(bytes32 rollupHash) external view returns (bool, Structs.MetaRollup memory) {
        Structs.MetaRollup memory rol = rollups.byHash[rollupHash];
        return (rol.Hash == rollupHash , rol);
    }
}