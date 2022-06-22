// SPDX-License-Identifier: GPL-3.0
import "libs/openzeppelin/cryptography/ECDSA.sol";

pragma solidity >=0.7.0 <0.9.0;

contract ManagementContract {
    // TODO these should all be private, but are public for testing in integration/smartcontract
    mapping(uint256 => Rollup[]) public rollups;
    mapping(address => string) public attestationRequests;
    mapping(address => bool) public attested;
    string[] public hostAddresses; // The addresses of all the Obscuro hosts on the network.

    // networkSecretNotInitialized marks if the network secret has been initialized
    bool private networkSecretInitialized ;

    struct Rollup{
        bytes32 ParentHash;
        address AggregatorID;
        bytes32 L1Block;
        uint256 Number;
    }

    function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string calldata rollupData) public {
        // TODO How to ensure the sender without hashing the calldata?
        // bytes32 derp = keccak256(abi.encodePacked(ParentHash, AggregatorID, L1Block, Number, rollupData));

        // revert if the AggregatorID is not attested
        require(attested[AggregatorID]);

        Rollup memory r = Rollup(ParentHash, AggregatorID, L1Block, Number);
        rollups[block.number].push(r);
    }

    // InitializeNetworkSecret kickstarts the network secret, can only be called once
    function InitializeNetworkSecret(address aggregatorID, bytes calldata initSecret) public {
        require(!networkSecretInitialized);

        // network can no longer be initialized
        networkSecretInitialized = true;

        // aggregator is now on the list of attested aggregators
        attested[aggregatorID] = true;
    }

    // Aggregators can request the Network Secret given an attestation request report
    function RequestNetworkSecret(string calldata requestReport) public {
        // Attestations should only be allowed to produce once ?
        attestationRequests[msg.sender] = requestReport;
    }

    // Attested node will pickup on Network Secret Request
    // and if valid will respond with the Network Secret
    // and mark the requesterID as attested
    function RespondNetworkSecret(address attesterID, address requesterID, bytes memory attesterSig, bytes memory responseSecret, string hostAddress) public {
        // only attested aggregators can respond to Network Secret Requests
        bool isAggAttested = attested[attesterID];
        require(isAggAttested);

        // the data must be signed with by the correct private key
        // signature = f(PubKey, PrivateKey, message)
        // address = f(signature, message)
        // valid if attesterID = address
        bytes32 calculatedHashSigned = ECDSA.toEthSignedMessageHash(abi.encodePacked(attesterID, requesterID, responseSecret));
        address recoveredAddrSignedCalculated = ECDSA.recover(calculatedHashSigned, attesterSig);

        // todo remove this toAsciiString helper
        require(recoveredAddrSignedCalculated == attesterID,
            string.concat("recovered address and attesterID don't match ",
                "\n Expected:                         ", toAsciiString(attesterID),
                "\n / recoveredAddrSignedCalculated:  ", toAsciiString(recoveredAddrSignedCalculated)));

        // mark the requesterID aggregator as an attested aggregator
        attested[requesterID] = true;
        hostAddresses.push(hostAddress);
    }

    function GetHostAddresses() public view returns (string[] memory) {
        return hostAddresses;
    }

    // Taken from https://ethereum.stackexchange.com/a/8447.
    function toAsciiString(address x) internal pure returns (string memory) {
        bytes memory s = new bytes(40);
        for (uint i = 0; i < 20; i++) {
            bytes1 b = bytes1(uint8(uint(uint160(x)) / (2**(8*(19 - i)))));
            bytes1 hi = bytes1(uint8(b) / 16);
            bytes1 lo = bytes1(uint8(b) - 16 * uint8(hi));
            s[2*i] = toChar(hi);
            s[2*i+1] = toChar(lo);
        }
        return string(s);
    }

    // Taken from https://ethereum.stackexchange.com/a/8447.
    function toChar(bytes1 b) internal pure returns (bytes1 c) {
        if (uint8(b) < 10) return bytes1(uint8(b) + 0x30);
        else return bytes1(uint8(b) + 0x57);
    }
}