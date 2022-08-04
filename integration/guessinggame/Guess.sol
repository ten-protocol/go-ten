//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.0;

interface IERC20 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function allowance(address owner, address spender) external view returns (uint256);

    function transfer(address recipient, uint256 amount) external returns (bool);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

contract Guess {
    address payable owner;
    uint8 private target;
    uint16 public guesses;
    IERC20 public erc20;

    uint8 public size = 100;
    address public tokenAddress = 0xf3a8bd422097bFdd9B3519Eaeb533393a1c561aC;

    modifier onlyOwner {
        require(msg.sender == owner, "Only owner can call this function.");
        _;
    }

    // For now we use hardcoded values to simplify deploying.
    constructor(){
        owner = payable(msg.sender);
        erc20 = IERC20(tokenAddress);
        setNewTarget();
    }

    function attempt(uint8 guess) public payable {
        require(erc20.allowance(msg.sender, address(this)) >= 1, "Check the token allowance.");
        erc20.transferFrom(msg.sender, address(this), 1);
        guesses++;
        if (guess == target) {
            erc20.transfer(msg.sender, erc20.balanceOf(address(this)));
            setNewTarget();
        }
    }

    function close() public payable onlyOwner {
        selfdestruct(payable(owner));
    }

    function getBalance() public view returns (uint256) {
        return erc20.balanceOf(address(this));
    }

    function setNewTarget() private {
        // We call this in the constructor.
        //        require(erc20.balanceOf(address(this)) == 0, "Balance must be zero to set a new target.");
        target = uint8(uint256(keccak256(abi.encodePacked(block.timestamp, block.difficulty))) % size);
    }
}