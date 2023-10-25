// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

import "./Structs.sol";
import * as MessageBus from "../messaging/MessageBus.sol";

contract ManagementContract is Ownable {

    event LogManagementContractCreated(address messageBusAddress);
    // Event to log changes to important contract addresses
    event ImportantContractAddressUpdated(string key, address newAddress);

    mapping(address => string) private attestationRequests;
    mapping(address => bool) private attested;
    // TODO - Revisit the decision to store the host addresses in the smart contract.
    string[] private hostAddresses; // The addresses of all the Obscuro hosts on the network.

    // In the near-term it is convenient to have an accessible source of truth for important contract addresses
    // TODO - this is probably not appropriate long term but currently useful for testnets. Look to remove.
    // We store the keys as well as the mapping for the key-value store for important contract addresses for convenience
    string[] public importantContractKeys;
    mapping (string => address) public importantContractAddresses;

    // networkSecretNotInitialized marks if the network secret has been initialized
    bool private networkSecretInitialized ;

    // isWithdrawalAvailable marks if the contract allows withdrawals or not
    bool private isWithdrawalAvailable;

    uint256 public lastBatchSeqNo = 0;

    Structs.RollupStorage private rollups;
    //The messageBus where messages can be sent to Obscuro
    MessageBus.IMessageBus public messageBus;
    constructor() {
        messageBus = new MessageBus.MessageBus();
        emit LogManagementContractCreated(address(messageBus));
    }

    function GetRollupByHash(bytes32 rollupHash) view public returns(bool, Structs.MetaRollup memory) {
        Structs.MetaRollup memory rol = rollups.byHash[rollupHash];
        return (rol.Hash == rollupHash , rol);
    }

    function AppendRollup(Structs.MetaRollup calldata _r) internal {
        rollups.byHash[_r.Hash] = _r;
        if (_r.LastSequenceNumber > lastBatchSeqNo) {
            lastBatchSeqNo = _r.LastSequenceNumber;
        }
    }
    //
    //  -- End of Tree element list Library
    //

// TODO: ensure challenge period is added on top of block timestamp.
    function pushCrossChainMessages(Structs.HeaderCrossChainData calldata crossChainData) internal {
        uint256 messagesLength = crossChainData.messages.length;
        for (uint256 i = 0; i < messagesLength; ++i) {
            messageBus.storeCrossChainMessage(crossChainData.messages[i], 1); //TODO - make finality depend on rollup challenge period
        }
    }

    // solc-ignore-next-line unused-param
    function AddRollup(Structs.MetaRollup calldata r, string calldata  _rollupData, Structs.HeaderCrossChainData calldata crossChainData) public {
        // TODO How to ensure the sender without hashing the calldata ?
        // bytes32 derp = keccak256(abi.encodePacked(ParentHash, AggregatorID, L1Block, Number, rollupData));

        // TODO: Add a check that ensures the cross messages are coming from the correct fork using block hashes.

        // revert if the AggregatorID is not attested
        require(attested[r.AggregatorID], "aggregator not attested");

        AppendRollup(r);
        pushCrossChainMessages(crossChainData);
    }

    // InitializeNetworkSecret kickstarts the network secret, can only be called once
    // solc-ignore-next-line unused-param
    function InitializeNetworkSecret(address _aggregatorID, bytes calldata  _initSecret, string memory _hostAddress, string calldata _genesisAttestation) public {
        require(!networkSecretInitialized);

        // network can no longer be initialized
        networkSecretInitialized = true;

        // aggregator is now on the list of attested aggregators and its host address is available
        attested[_aggregatorID] = true;
        hostAddresses.push(_hostAddress);
    }

    // Aggregators can request the Network Secret given an attestation request report
    function RequestNetworkSecret(string calldata requestReport) public {
        // Attestations should only be allowed to produce once ?
        attestationRequests[msg.sender] = requestReport;
    }

    // Attested node will pickup on Network Secret Request
    // and if valid will respond with the Network Secret
    // and mark the requesterID as attested
    // @param verifyAttester Whether to ask the attester to complete a challenge (signing a hash) to prove their identity.
    function RespondNetworkSecret(address attesterID, address requesterID, bytes memory attesterSig, bytes memory responseSecret, string memory hostAddress, bool verifyAttester) public {
        // only attested aggregators can respond to Network Secret Requests
        bool isAggAttested = attested[attesterID];
        require(isAggAttested);

        if (verifyAttester) {
            // the data must be signed with by the correct private key
            // signature = f(PubKey, PrivateKey, message)
            // address = f(signature, message)
            // valid if attesterID = address
            bytes32 calculatedHashSigned = ECDSA.toEthSignedMessageHash(abi.encodePacked(attesterID, requesterID, hostAddress, responseSecret));
            address recoveredAddrSignedCalculated = ECDSA.recover(calculatedHashSigned, attesterSig);

        require(recoveredAddrSignedCalculated == attesterID, "calculated address and attesterID dont match");
        }

        // mark the requesterID aggregator as an attested aggregator and store its host address
        attested[requesterID] = true;
        // TODO - Consider whether to remove duplicates.
        hostAddresses.push(hostAddress);
    }

    function GetHostAddresses() public view returns (string[] memory) {
        return hostAddresses;
    }


    // Accessor to check if the contract is locked or not
    function IsWithdrawalAvailable() view public returns (bool) {
        return isWithdrawalAvailable;
    }

    // Accessor that checks if an address is attested or not
    function Attested(address _addr) view public returns (bool) {
        return attested[_addr];
    }

    // Testnet function to allow the contract owner to retrieve **all** funds from the network bridge.
    function RetrieveAllBridgeFunds() public onlyOwner {
        messageBus.retrieveAllFunds(msg.sender);
    }

    // Function to set an important contract's address, only callable by owner
    function setImportantContractAddress(string memory key, address newAddress) public onlyOwner {
        if (importantContractAddresses[key] == address(0)) {
            importantContractKeys.push(key);
        }
        importantContractAddresses[key] = newAddress;
        emit ImportantContractAddressUpdated(key, newAddress);
    }

    function getImportantContractKeys() public view returns(string[] memory) {
        return importantContractKeys;
    }
}
