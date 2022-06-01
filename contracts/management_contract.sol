// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ManagementContract {

    mapping(uint256 => string[]) public rollups;
    mapping(address => string) public attestationRequests;

    mapping(string => string) public attestations;
    mapping(string => string) public attested;

    string networkSecret;

    function AddRollup(string calldata rollupData) public {
        rollups[block.number].push(rollupData);
    }


    // Aggregators can request the Network Secret given an attestation request report
    function RequestNetworkSecret(string calldata requestReport) public {
        // Attestations should only be allowed to produce once ?
        attestationRequests[msg.sender] = requestReport;
    }

    // Genesis node ( for now ) will pickup on Network Secret Request
    // and if valid will respond with the Network Secret
    function RespondNetworkSecret(string memory requester, string memory pubKey, string memory inputSecret, string calldata requestReport) public {
        attested[requester] = pubKey;
    }
}