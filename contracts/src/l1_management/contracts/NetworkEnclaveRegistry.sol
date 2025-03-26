// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import "../interfaces/INetworkEnclaveRegistry.sol";

contract NetworkEnclaveRegistry is INetworkEnclaveRegistry, Initializable, OwnableUpgradeable {
    using MessageHashUtils for bytes32;

    bool private networkSecretInitialized;

    // mapping of enclaveID to whether it is attested
    mapping(address => bool) private attested;

    // mapping of enclaveID to whether it is permissioned as a sequencer enclave
    // note: the enclaveID which initialises the network secret is automatically permissioned as a sequencer.
    // Beyond that, the contract owner can grant and revoke sequencer status.
    mapping(address => bool) private sequencerEnclave;

    constructor() {
        _transferOwnership(msg.sender);
    }

    function initialize(address owner) public initializer {
        __Ownable_init(owner);
        networkSecretInitialized = false;
    }

    // initializeNetworkSecret kickstarts the network secret, can only be called once
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

    // Enclaves can request the Network Secret given an attestation request report
    function requestNetworkSecret(string calldata requestReport) external {
        // once an enclave has been attested there is no need for them to request this again
        require(!attested[msg.sender], "already attested");
        emit NetworkSecretRequested(msg.sender, requestReport);
    }

    // An attested enclave will pickup the Network Secret Request and, if valid, will respond with the Network Secret
    // and mark the requesterID as attested
    // @param verifyAttester Whether to ask the attester to complete a challenge (signing a hash) to prove their identity.
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

    // Accessor that checks if an enclave address has been attested
    function isAttested(address enclaveID) external view returns (bool) {
        return attested[enclaveID];
    }

    // Accessor that checks if an address is permissioned as a sequencer
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