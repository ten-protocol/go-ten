// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import "./INetworkEnclaveRegistry.sol";

contract NetworkEnclaveRegistry is INetworkEnclaveRegistry, Initializable, OwnableUpgradeable {
    using MessageHashUtils for bytes32;

    bool private networkSecretInitialized;

    // mapping of enclaveID to whether it is attested
    mapping(address => bool) private attested;

    // mapping of enclaveID to whether it is permissioned as a sequencer enclave
    // note: the enclaveID which initialises the network secret is automatically permissioned as a sequencer.
    // Beyond that, the contract owner can grant and revoke sequencer status.
    mapping(address => bool) private sequencerEnclave;

    function initialize() public initializer {
        __Ownable_init(msg.sender);
        networkSecretInitialized = false;
    }

    function initializeNetworkSecret(
        address enclaveID,
        bytes calldata initSecret,
        string calldata genesisAttestation
    ) external {
        require(!networkSecretInitialized, "network secret already initialized");
        require(enclaveID != address(0), "invalid enclave address");

        networkSecretInitialized = true;
        attested[enclaveID] = true;

        sequencerEnclave[enclaveID] = true;
        emit NetworkSecretInitialized(enclaveID);
        emit EnclaveAttested(enclaveID);
    }

    function isInitialized() external view returns (bool) {
        return networkSecretInitialized;
    }

    function requestNetworkSecret(string calldata requestReport) external {
        require(!attested[msg.sender], "already attested");
        emit NetworkSecretRequested(msg.sender, requestReport);
    }

    function respondNetworkSecret(
        address attesterID,
        address requesterID,
        bytes memory attesterSig,
        bytes memory responseSecret,
        bool verifyAttester
    ) external {
        require(attested[attesterID], "responding attester is not attested");
        require(!attested[requesterID], "requester already attested");
        require(requesterID != address(0), "invalid requester address");

        if (verifyAttester) {
            bytes32 messageHash = keccak256(
                abi.encodePacked(
                    requesterID,
                    responseSecret
                )
            ).toEthSignedMessageHash();

            address recoveredAddr = ECDSA.recover(messageHash, attesterSig);
            require(recoveredAddr == attesterID, "invalid signature");
        }

        attested[requesterID] = true;
        emit NetworkSecretResponded(_attesterID, requesterID);
        emit EnclaveAttested(requesterID);
    }

    function isAttested(address enclaveID) external view returns (bool) {
        return attested[enclaveID];
    }

    function isSequencer(address enclaveID) external view returns (bool) {
        return sequencerEnclave[enclaveID];
    }

    // Function to grant sequencer status for an enclave - contract owner only
    function grantSequencerEnclave(address _addr) external onlyOwner {
        // require the enclave to be attested already
        require(attested[_addr], "enclaveID not attested");
        sequencerEnclave[_addr] = true;
        emit SequencerEnclaveGranted(_addr);
    }
    
    // Function to revoke sequencer status for an enclave - contract owner only
    function revokeSequencerEnclave(address _addr) external onlyOwner {
        // require the enclave to be a sequencer already
        require(sequencerEnclave[_addr], "enclaveID not a sequencer");
        delete sequencerEnclave[_addr];
        emit SequencerEnclaveRevoked(_addr);
    }
}