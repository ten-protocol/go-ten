// SPDX-License-Identifier: MIT
// Specify the Solidity version
pragma solidity ^0.8.0;

contract SimpleMessageContract {

    // State variable to store the message
    string public message;
    string public message2;

    // Event declaration
    event MessageUpdatedWithAddress(string newMessage, address indexed sender);
    event Message2Updated(string newMessage);

    // Constructor to initialize the message
    constructor() {
        message = "foo";
        message2 = "foo";
    }

    // Function to set a new message
    function setMessage(string memory newMessage) public {
        message = newMessage;
        emit MessageUpdatedWithAddress(newMessage, msg.sender);  // Emit the event (only sender can see it)
    }

    function setMessage2(string memory newMessage) public {
        message2 = newMessage;
        emit Message2Updated(newMessage);  // Emit the event (everyone can see it)
    }
}