// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

/**
 * @title INetworkEnclaveRegistry
 * @dev Interface for managing registration process of enclaves in the network
 */
interface INetworkEnclaveRegistry {
    event NetworkSecretInitialized(address enclaveID);
    event NetworkSecretRequested(address indexed requester, string requestReport);
    event NetworkSecretResponded(address indexed attester, address indexed requester);
    event SequencerEnclaveGranted(address enclaveID);
    event SequencerEnclaveRevoked(address enclaveID);

    /**
     * @dev Initializes the network's secret, can only be called once
     * @param enclaveID Address of the initializing enclave
     * @param _initSecret Initial secret data for the network
     * @param _genesisAttestation Attestation proof for the genesis enclave (must be 116 bytes)
     * @notice The initializing enclave automatically becomes the first sequencer
     */
    function initializeNetworkSecret(address enclaveID, bytes calldata _initSecret, string calldata _genesisAttestation) external;

    /**
     * @dev Allows an unattested enclave to request the network secret
     * @param requestReport Attestation report from the requesting enclave
     * @notice Can only be called by enclaves that have not been attested yet
     * @notice Emits NetworkSecretRequested event for attested enclaves to respond
     */
    function requestNetworkSecret(string calldata requestReport) external;

    /**
     * @dev Processes a response to a network secret request from an attested enclave
     * @param attesterID Address of the attested enclave providing the secret
     * @param requesterID Address of the enclave that requested the secret
     * @param attesterSig Signature from the attesting enclave (if verification required)
     * @param responseSecret Encrypted network secret (must be 145 bytes)
     * @param verifyAttester If true, validates attester's signature of requesterID + responseSecret
     * @notice Attester must be already attested
     * @notice Requester must not be already attested
     */
    function respondNetworkSecret(
        address attesterID,
        address requesterID,
        bytes memory attesterSig,
        bytes memory responseSecret,
        bool verifyAttester
    ) external;

    /**
     * @dev Checks if an enclave has been attested
     * @param enclaveID Address of the enclave to check
     * @return bool True if the enclave has been successfully attested
     */
    function isAttested(address enclaveID) external view returns (bool);

    /**
     * @dev Checks if an enclave has sequencer permissions
     * @param enclaveID Address of the enclave to check
     * @return bool True if the enclave is permitted to act as a sequencer
     */
    function isSequencer(address enclaveID) external view returns (bool);

    /**
     * @dev Grants sequencer permissions to an attested enclave
     * @param _addr Address of the enclave to grant sequencer status
     * @notice Can only be called by contract owner
     * @notice Enclave must be attested before it can be granted sequencer status
     */
    function grantSequencerEnclave(address _addr) external;

    /**
     * @dev Revokes sequencer permissions from an enclave
     * @param _addr Address of the enclave to revoke sequencer status
     * @notice Can only be called by contract owner
     * @notice Enclave must currently have sequencer status to be revoked
     */
    function revokeSequencerEnclave(address _addr) external;
}