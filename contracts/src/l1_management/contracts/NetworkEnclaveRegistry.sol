// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import "../interfaces/INetworkEnclaveRegistry.sol";
import "../../common/UnrenouncableOwnable2Step.sol";
/**
 * @title NetworkEnclaveRegistry
 * @dev Contract for managing network enclave registry
 * Implements a network secret initialization and attestation mechanism
 * Allows enclaves to request and respond to the network secret
 * Provides sequencer enclave status management
*/
contract NetworkEnclaveRegistry is INetworkEnclaveRegistry, Initializable, UnrenouncableOwnable2Step {
    
    using MessageHashUtils for bytes32;

    /**
     * @dev Flag to check if the network secret has been initialized
     */
    bool private networkSecretInitialized;

    /**
     * @dev Mapping of enclaveID to whether it is attested
     */
    mapping(address enclaveID => bool isAttested) private attested;

    /**
     * @dev Mapping of enclaveID to whether it is permissioned as a sequencer enclave. The enclaveID which initialises
     * the network secret is automatically permissioned as a sequencer. Beyond that, the contract owner can grant and revoke
     * sequencer status.
     */
    mapping(address sequencerID => bool isSequencer) private sequencerEnclave;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Initializes the contract with the owner
     * @param _owner Address of the contract owner
     */
    function initialize(address _owner) public initializer {
        __UnrenouncableOwnable2Step_init(_owner);  // Initialize UnrenouncableOwnable2Step
        networkSecretInitialized = false;
    }

    /**
     * @dev Initializes the network secret, can only be called once.
     * @param enclaveID The enclaveID of the enclave that is initializing the network secret
     * @param _initSecret The initial secret
     * @param _genesisAttestation The genesis attestation
     */
    // solc-ignore-next-line unused-param
    function initializeNetworkSecret(address enclaveID, bytes calldata _initSecret, string calldata _genesisAttestation) external {
        require(!networkSecretInitialized, "network secret already initialized");
        require(enclaveID != address(0), "invalid enclave address");

        // network can no longer be initialized
        networkSecretInitialized = true;

        // enclave is now on the list of attested enclaves (and its host address is published for p2p)
        attested[enclaveID] = true;

        // the enclave that starts the network with this call is implicitly a sequencer so doesn't need adding
        sequencerEnclave[enclaveID] = true;
        emit NetworkSecretInitialized(enclaveID);
    }

    /**
     * @dev Requests the network secret, can only be called by an attested enclave.
     * @param requestReport The request report
     */
    function requestNetworkSecret(string calldata requestReport) external {
        // once an enclave has been attested there is no need for them to request this again
        require(!attested[msg.sender], "already attested");
        emit NetworkSecretRequested(msg.sender, requestReport);
    }

    /**
     * @dev Responds to the network secret request, can only be called by an attested enclave.
     * @param attesterID The enclaveID of the enclave that is responding to the request
     * @param requesterID The enclaveID of the enclave that is requesting the network secret
     * @param attesterSig The signature of the attester
     * @param responseSecret The response secret
     * @param verifyAttester Whether to ask the attester to complete a challenge (signing a hash) to prove their identity.
     */
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
        require(responseSecret.length == 145, "invalid secret response lenght");

        if (verifyAttester) {
            // the data must be signed with by the correct private key
            // signature = f(PubKey, PrivateKey, message)
            // address = f(signature, message)
            // valid if attesterID = address
            bytes32 messageHash = keccak256(
                abi.encodePacked(
                    requesterID,
                    responseSecret
                )
            ).toEthSignedMessageHash();

            address recoveredAddr = ECDSA.recover(messageHash, attesterSig);
            require(recoveredAddr == attesterID, "invalid signature");
        }

        // mark the requesterID enclave as an attested enclave and store its host address
        attested[requesterID] = true;
        emit NetworkSecretResponded(attesterID, requesterID);
    }

    /**
     * @dev Checks if an enclave address has been attested
     * @param enclaveID The enclaveID of the enclave to check
     * @return bool True if the enclave is attested, false otherwise
     */
    function isAttested(address enclaveID) external view returns (bool) {
        return attested[enclaveID];
    }

    /**
     * @dev Checks if an enclave address is permissioned as a sequencer
     * @param enclaveID The enclaveID of the enclave to check
     * @return bool True if the enclave is a sequencer, false otherwise
     */
    function isSequencer(address enclaveID) external view returns (bool) {
        return sequencerEnclave[enclaveID];
    }

    /**
     * @dev Grants sequencer status to an enclave, can only be called by the contract owner.
     * @param _addr The enclaveID of the enclave to grant sequencer status to
     */
    function grantSequencerEnclave(address _addr) external onlyOwner {
        // require the enclave to be attested already
        require(attested[_addr], "enclaveID not attested");
        sequencerEnclave[_addr] = true;
        emit SequencerEnclaveGranted(_addr);
    }

    /**
     * @dev Revokes sequencer status from an enclave, can only be called by the contract owner.
     * @param _addr The enclaveID of the enclave to revoke sequencer status from
     */
    function revokeSequencerEnclave(address _addr) external onlyOwner {
        // require the enclave to be a sequencer already
        require(sequencerEnclave[_addr], "enclaveID not a sequencer");
        delete sequencerEnclave[_addr];
        emit SequencerEnclaveRevoked(_addr);
    }
}