// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ManagementContract {

    mapping(uint256 => Rollup[]) public rollups;

    mapping(address => string) public attestationRequests;
    mapping(bytes20 => string) public attested;

    string networkSecret;

    struct Rollup{
        bytes32 ParentHash;
        bytes20 AggregatorID;
        bytes32 L1Block;
        uint256 Number;
    }

    function AddRollup(bytes32 ParentHash, bytes20 AggregatorID, bytes32 L1Block, uint256 Number, string calldata rollupData) public {
        // TODO How to ensure the sender without hashing the calldata ?
        // bytes32 derp = keccak256(abi.encodePacked(ParentHash, AggregatorID, L1Block, Number, rollupData));

        // revert if the AggregatorID is not attested
        require(bytes(attested[AggregatorID]).length > 0);

        Rollup memory r = Rollup(ParentHash, AggregatorID, L1Block, Number);
        rollups[block.number].push(r);
    }


    // Aggregators can request the Network Secret given an attestation request report
    function RequestNetworkSecret(string calldata requestReport) public {
        // Attestations should only be allowed to produce once ?
        attestationRequests[msg.sender] = requestReport;
    }

    // Genesis node ( for now ) will pickup on Network Secret Request
    // and if valid will respond with the Network Secret
    function RespondNetworkSecret(bytes20 requesterID, string memory pubKey, string memory inputSecret, string calldata requestReport) public {
        attested[requesterID] = pubKey;
    }
}