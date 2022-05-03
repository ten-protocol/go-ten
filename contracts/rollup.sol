// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract RollupStorage {

    mapping(uint256 => string[]) public rollups;
    mapping(address => uint256) public deposits;
    mapping(address => string) public attestationRequests;
    string secret;

    function AddRollup(string calldata rollupData) public {
        rollups[block.number].push(rollupData);
    }


    function Rollup() public view returns (string[] memory){
        return rollups[block.number];
    }

    function Deposit(address dest) public payable {
        deposits[dest] = msg.value;
    }

    function StoreSecret(string memory inputSecret, string calldata requestReport) public {
        secret = inputSecret;
    }

    function RequestSecret(string calldata requestReport) public {
        // we probably don't need to persist this in state (at least not in its entirety)
        attestationRequests[msg.sender] = requestReport;
    }

    function Withdraw(uint256 withdrawAmount, address payable destination) public {
        destination.transfer(withdrawAmount);
    }

}