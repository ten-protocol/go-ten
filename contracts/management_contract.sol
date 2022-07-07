// SPDX-License-Identifier: GPL-3.0
import "libs/openzeppelin/cryptography/ECDSA.sol";

pragma solidity >=0.7.0 <0.9.0;

contract ManagementContract {

    mapping(address => string) private attestationRequests;
    mapping(address => bool) private attested;
    // TODO - Revisit the decision to store the host addresses in the smart contract.
    string[] private hostAddresses; // The addresses of all the Obscuro hosts on the network.

    // tree holds a tree of rollups
    Tree private tree;

    // networkSecretNotInitialized marks if the network secret has been initialized
    bool private networkSecretInitialized ;

    // isWithdrawalAvailable marks if the contract allows withdrawals or not
    bool private isWithdrawalAvailable;

    // MetaRollup is a rollup meta data
    struct MetaRollup{
        bytes32 ParentHash;
        bytes32 Hash;
        address AggregatorID;
        bytes32 L1Block;
        uint256 Number;
    }

    // TreeElement is an element of the Tree structure
    struct TreeElement{
        uint256 ElementID;
        uint256 ParentID;
        MetaRollup rollup;
    }

    // NonExisting - 0 (Constant)
    // Tail - 1 (Constant)
    // Head - X (Variable)
    // Does not use rollup hashes as a storing ID as they can be compromised
    struct Tree {
        // rollups stores the Elements using incremental IDs
        mapping(uint256 => TreeElement) rollups;
        // map a rollup hash to a storage ID
        mapping(bytes32 => uint256) rollupsHashes;
        // map the children of a node
        mapping(uint256 => uint256[]) rollupChildren;

        uint256 _TAIL; // tail is always 1
        uint256 _HEAD;
        uint256 _nextID; // TODO use openzeppelin counters
        bool initialized;
    }

    //
    //  -- Start of Tree element list Library
    //

    // InitializeTree starts the list and sets the initial values
    function InitializeTree(MetaRollup memory r) public {
        require(!tree.initialized, "cannot be initialized again");
        tree.initialized = true;

        // TreeElement starts a 1 and has no parent ( ParentID: 0 )
        tree.rollups[1] = TreeElement(1, 0, r);
        tree._HEAD = 1;
        tree._nextID = 2;
        tree.rollupsHashes[r.Hash] = 1;

        // withdrawals are available at the start
        isWithdrawalAvailable = true;
    }

    function GetRollupByID(uint256 rollupID) view public returns(bool, TreeElement memory) {
        TreeElement memory rol = tree.rollups[rollupID];
        return (rol.ElementID != 0 , rol);
    }

    function GetRollupByHash(bytes32 rollupHash) view public returns (bool, TreeElement memory) {
        return GetRollupByID(tree.rollupsHashes[rollupHash]);
    }

    function GetHeadRollup() internal view returns ( TreeElement memory ) {
        return tree.rollups[tree._HEAD];
    }

    function GetParentRollup(TreeElement memory element) view public returns( bool, TreeElement memory) {
        return GetRollupByID(element.ParentID);
    }

    function AppendRollup(uint256 _parentID, MetaRollup memory _r) public {
        // guarantee the storage ids are not compromised
        uint rollupID = tree._nextID;
        tree._nextID++;

        // cannot append to non-existing parent rollups
        (bool found, TreeElement memory parent) = GetRollupByID(_parentID);
        require(found, "parent not found");

        // store the rollup in an element
        tree.rollups[rollupID] = TreeElement(rollupID, _parentID, _r);

        // mark the element as a child of parent
        tree.rollupChildren[_parentID].push(rollupID);

        // store the hashpointer
        tree.rollupsHashes[_r.Hash] = rollupID;

        // mark this as the head
        if (parent.ElementID == tree._HEAD) {
            tree._HEAD = rollupID;
        }
    }

    // HasSecondCousinFork returns whether there is a fork in the current view of the rollups
    // It works by:
    // - Traversing up two levels ( from the HEAD to the grand father element )
    // - Checking if there are siblings ( at the grand father level )
    // - Checking if the siblings have children ( meaning that a fork expanded )
    //
    // Will return true when a rollup 6 or 6' with parent 5 or 5' is inserted
    // 0 -> 1 -> 2 -> 3 -> 4 -> 5
    //                  -> 4'-> 5'
    //
    function HasSecondCousinFork() view public returns (bool) {
        TreeElement memory currentElement = GetHeadRollup();

        // traverse up to the grandpa ( 2 levels up )
        (bool foundParent, TreeElement memory parentElement) = GetParentRollup(currentElement);
        require(foundParent, "no parent");
        (bool foundGrandpa, TreeElement memory grandpaElement) = GetParentRollup(parentElement);
        require(foundGrandpa, "no grand parent");

        // follow each of the grandpa children until it's two levels deep
        uint256[] memory childrenIDs = tree.rollupChildren[grandpaElement.ElementID];
        for (uint256 i = 0; i < childrenIDs.length ; i++) {
            (bool foundChild, TreeElement memory child) = GetRollupByID(childrenIDs[i]);

            // no more children
            if (!foundChild) {
                return false;
            }

            // ignore the current tree
            if (child.ElementID == parentElement.ElementID ) {
                continue;
            }

            // if child has children then it's bad ( fork of depth 2 )
            if (tree.rollupChildren[child.ElementID].length > 0) {
                return true;
            }
        }

        return false;
    }

    //
    //  -- End of Tree element list Library
    //

    function AddRollup(bytes32 _parentHash, bytes32 _hash, address _aggregatorID, bytes32 _l1Block, uint256 _number, string calldata _rollupData) public {
        // TODO How to ensure the sender without hashing the calldata ?
        // bytes32 derp = keccak256(abi.encodePacked(ParentHash, AggregatorID, L1Block, Number, rollupData));

        // revert if the AggregatorID is not attested
        require(attested[_aggregatorID], "aggregator not attested");

        MetaRollup memory r = MetaRollup(_parentHash, _hash, _aggregatorID, _l1Block, _number);

        // if this is the first element initialize the tree structure
        // TODO this should be moved to the network initialization
        if (!tree.initialized) {
            InitializeTree(r);
            return;
        }

        (bool found, TreeElement memory parent) = GetRollupByHash(_parentHash);
        require(found, "unable to find parent hash");

        // don't check for forks at the start
        if (tree._HEAD > 2) {
            bool forkFound = HasSecondCousinFork();
            if (forkFound) {
                isWithdrawalAvailable = false;
                // We keep accepting rollups just locks the contract
                // require(!found, "detected a fork");
            }
        }

        AppendRollup(parent.ElementID, r);
    }

    // InitializeNetworkSecret kickstarts the network secret, can only be called once
    function InitializeNetworkSecret(address _aggregatorID, bytes calldata _initSecret, string memory _hostAddress, string calldata genesisAttestation) public {
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
}
