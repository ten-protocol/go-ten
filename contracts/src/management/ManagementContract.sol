// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import "@openzeppelin/contracts/utils/Strings.sol";


import "./Structs.sol";
import * as MessageStructs from "../messaging/Structs.sol";
import * as MessageBus from "../messaging/MessageBus.sol";
import * as MerkleTreeMessageBus from "../messaging/MerkleTreeMessageBus.sol";

contract ManagementContract is Initializable, OwnableUpgradeable {

    using MessageHashUtils for bytes32;
    using MessageHashUtils for bytes;

    constructor() {
        _transferOwnership(msg.sender);
    }

    event LogManagementContractCreated(address messageBusAddress);
    event ImportantContractAddressUpdated(string key, address newAddress);
    event SequencerEnclaveGranted(address enclaveID);
    event SequencerEnclaveRevoked(address enclaveID);
    event RollupAdded(bytes32 rollupHash);
    event NetworkSecretRequested(address indexed requester, string requestReport);
    event NetworkSecretResponded(address indexed attester, address indexed requester);

    // mapping of enclaveID to whether it is attested
    mapping(address => bool) private attested;

    // mapping of enclaveID to whether it is permissioned as a sequencer enclave
    // note: the enclaveID which initialises the network secret is automatically permissioned as a sequencer.
    //       Beyond that, the contract owner can grant and revoke sequencer status.
    mapping(address => bool) private sequencerEnclave;

    // In the near-term it is convenient to have an accessible source of truth for important contract addresses
    // TODO - this is probably not appropriate long term but currently useful for testnets. Look to remove.
    // We store the keys as well as the mapping for the key-value store for important contract addresses for convenience
    string[] public importantContractKeys;
    mapping (string => address) public importantContractAddresses;

    // networkSecretNotInitialized marks if the network secret has been initialized
    bool private networkSecretInitialized ;

    // isWithdrawalAvailable marks if the contract allows withdrawals or not
    bool private isWithdrawalAvailable;

    uint256 public lastBatchSeqNo;

    Structs.RollupStorage private rollups;
    //The messageBus where messages can be sent to Obscuro
    MessageBus.IMessageBus public messageBus;
    MerkleTreeMessageBus.IMerkleTreeMessageBus public merkleMessageBus;
    mapping(bytes32 =>bool) public isWithdrawalSpent;
    mapping(bytes32 =>bool) public isBundleSaved;

    bytes32 public lastBatchHash;

    uint256 private challengePeriod;

    function initialize() public initializer {
        __Ownable_init(msg.sender);
        lastBatchSeqNo = 0;
        merkleMessageBus = new MerkleTreeMessageBus.MerkleTreeMessageBus();
        messageBus = MessageBus.IMessageBus(address(merkleMessageBus));

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

    function isBundleAvailable(bytes[] memory crossChainHashes) public view returns (bool) {
        bytes32 bundleHash = bytes32(0);

        for(uint256 i = 0; i < crossChainHashes.length; i++) {
            bundleHash = keccak256(abi.encode(bundleHash, bytes32(crossChainHashes[i])));
        }

        return isBundleSaved[bundleHash];
    }

    function pushCrossChainMessages(Structs.HeaderCrossChainData calldata crossChainData) internal {
        uint256 messagesLength = crossChainData.messages.length;
        for (uint256 i = 0; i < messagesLength; ++i) {
            messageBus.storeCrossChainMessage(crossChainData.messages[i], 1); //TODO - make finality depend on rollup challenge period
        }
    }

    modifier verifyRollupIntegrity(Structs.MetaRollup calldata r) {
        // Block binding checks
        require(block.number >= r.BlockBindingNumber, "Cannot bind to future block");
        require(block.number < (r.BlockBindingNumber + 255), "Block binding too old");

        bytes32 knownBlockHash = blockhash(r.BlockBindingNumber);

        require(knownBlockHash != 0x0, "Unknown block hash");
        require(knownBlockHash == r.BlockBindingHash, "Block binding mismatch");

        bytes32 compositeHash = keccak256(abi.encodePacked(
            r.LastSequenceNumber,
            r.BlockBindingHash,
            r.BlockBindingNumber,
            r.crossChainRoot,
            r.BlobHash
        ));

        // Verify the hash matches the one in the rollup
        require(compositeHash == r.CompositeHash, "Composite hash mismatch");

        // Verify the enclave signature
        address enclaveID = ECDSA.recover(compositeHash, r.Signature);
        require(attested[enclaveID], "enclaveID not attested");
        require(sequencerEnclave[enclaveID], "enclaveID not a sequencer");
        _;
    }

    // solc-ignore-next-line unused-param
    function AddRollup(Structs.MetaRollup calldata r) public verifyRollupIntegrity(r) {
        AppendRollup(r);

        if (r.crossChainRoot != bytes32(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff)) {
            merkleMessageBus.addStateRoot(r.crossChainRoot, block.timestamp);
        }
        emit RollupAdded(r.Hash);
    }

    // InitializeNetworkSecret kickstarts the network secret, can only be called once
    // solc-ignore-next-line unused-param
    function InitializeNetworkSecret(address _enclaveID, bytes calldata  _initSecret, string calldata _genesisAttestation) public {
        require(!networkSecretInitialized, "network secret already initialized");

        // network can no longer be initialized
        networkSecretInitialized = true;

        // enclave is now on the list of attested enclaves (and its host address is published for p2p)
        attested[_enclaveID] = true;

        // the enclave that starts the network with this call is implicitly a sequencer so doesn't need adding
        sequencerEnclave[_enclaveID] = true;
        emit SequencerEnclaveGranted(_enclaveID);
    }

    // Enclaves can request the Network Secret given an attestation request report
    function RequestNetworkSecret(string calldata requestReport) public {
        emit NetworkSecretRequested(msg.sender, requestReport);
    }

    function ExtractNativeValue(MessageStructs.Structs.ValueTransferMessage calldata _msg, bytes32[] calldata proof, bytes32 root) external {
        merkleMessageBus.verifyValueTransferInclusion(_msg, proof, root);
        bytes32 msgHash = keccak256(abi.encode(_msg));
        require(isWithdrawalSpent[msgHash] == false, "withdrawal already spent");
        isWithdrawalSpent[keccak256(abi.encode(_msg))] = true;

        messageBus.receiveValueFromL2(_msg.receiver, _msg.amount);
    }

    // An attested enclave will pickup the Network Secret Request
    // and, if valid, will respond with the Network Secret
    // and mark the requesterID as attested
    // @param verifyAttester Whether to ask the attester to complete a challenge (signing a hash) to prove their identity.
    function RespondNetworkSecret(address attesterID, address requesterID, bytes memory attesterSig, bytes memory responseSecret, bool verifyAttester) public {
        // only attested enclaves can respond to Network Secret Requests
        bool isEnclAttested = attested[attesterID];
        require(isEnclAttested, "responding attester is not attested");

        if (verifyAttester) {

            // the data must be signed with by the correct private key
            // signature = f(PubKey, PrivateKey, message)
            // address = f(signature, message)
            // valid if attesterID = address
            bytes32 calculatedHashSigned = abi.encodePacked(attesterID, requesterID, responseSecret).toEthSignedMessageHash();
            address recoveredAddrSignedCalculated = ECDSA.recover(calculatedHashSigned, attesterSig);

            require(recoveredAddrSignedCalculated == attesterID, "calculated address and attesterID dont match");
        }

        // mark the requesterID enclave as an attested enclave and store its host address
        attested[requesterID] = true;

        emit NetworkSecretResponded(attesterID, requesterID);
    }


    // Accessor to check if the contract is locked or not
    function IsWithdrawalAvailable() view public returns (bool) {
        return isWithdrawalAvailable;
    }

    // Accessor that checks if an address is attested or not
    function Attested(address _addr) view public returns (bool) {
        return attested[_addr];
    }

    // Accessor that checks if an address is permissioned as a sequencer
    function IsSequencerEnclave(address _addr) view public returns (bool) {
        return sequencerEnclave[_addr];
    }

    // Function to grant sequencer status for an enclave - contract owner only
    function GrantSequencerEnclave(address _addr) public onlyOwner {
        // require the enclave to be attested already
        require(attested[_addr], "enclaveID not attested");
        sequencerEnclave[_addr] = true;
        emit SequencerEnclaveGranted(_addr);
    }
    // Function to revoke sequencer status for an enclave - contract owner only
    function RevokeSequencerEnclave(address _addr) public onlyOwner {
        // require the enclave to be a sequencer already
        require(sequencerEnclave[_addr], "enclaveID not a sequencer");
        delete sequencerEnclave[_addr];
        emit SequencerEnclaveRevoked(_addr);
    }

    // Testnet function to allow the contract owner to retrieve **all** funds from the network bridge.
    function RetrieveAllBridgeFunds() public onlyOwner {
        messageBus.retrieveAllFunds(msg.sender);
    }

    // Function to set an important contract's address, only callable by owner
    function SetImportantContractAddress(string memory key, address newAddress) public onlyOwner {
        if (importantContractAddresses[key] == address(0)) {
            importantContractKeys.push(key);
        }
        importantContractAddresses[key] = newAddress;
        emit ImportantContractAddressUpdated(key, newAddress);
    }

    function GetImportantContractKeys() public view returns(string[] memory) {
        return importantContractKeys;
    }

    // Return the challenge period delay for message bus root
    function GetChallengePeriod() public view returns (uint256) {
        return challengePeriod;
    }

    // Sets the challenge period for message bus root (owner only)
    function SetChallengePeriod(uint256 _delay) public onlyOwner {
        challengePeriod = _delay;
    }
}
