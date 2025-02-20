// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

interface INetworkEnclaveRegistry {
    event NetworkSecretInitialized(address initializer);
    event NetworkSecretRequested(address indexed requester, string requestReport);
    event NetworkSecretResponded(address indexed attester, address indexed requester);
    event SequencerEnclaveGranted(address enclaveID);
    event SequencerEnclaveRevoked(address enclaveID);

    // FIXME - can we delete these unused params?
    // solc-ignore-next-line unused-param
    function initializeNetworkSecret(address enclaveID, bytes calldata _initSecret, string calldata _genesisAttestation) external;
    function requestNetworkSecret(string calldata requestReport) external;
    function respondNetworkSecret(
        address attesterID,
        address requesterID,
        bytes memory attesterSig,
        bytes memory responseSecret,
        bool verifyAttester
    ) external;
    function isAttested(address enclaveID) external view returns (bool);
    function isSequencer(address enclaveID) external view returns (bool);
    function grantSequencerEnclave(address _addr) external;
    function revokeSequencerEnclave(address _addr) external;
}