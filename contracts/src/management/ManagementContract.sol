// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/utils/Strings.sol";


import "../common/Structs.sol";
import * as MessageBus from "../messaging/MessageBus.sol";
import * as MerkleTreeMessageBus from "../messaging/MerkleTreeMessageBus.sol";

contract ManagementContract is Initializable, OwnableUpgradeable {

    constructor() {
        _transferOwnership(msg.sender);
    }

    event LogManagementContractCreated(address messageBusAddress);
    event ImportantContractAddressUpdated(string key, address newAddress);

    // In the near-term it is convenient to have an accessible source of truth for important contract addresses
    // TODO - this is probably not appropriate long term but currently useful for testnets. Look to remove.
    // We store the keys as well as the mapping for the key-value store for important contract addresses for convenience
    string[] public importantContractKeys;
    mapping (string => address) public importantContractAddresses;

    function initialize() public initializer {
        __Ownable_init(msg.sender);
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
}