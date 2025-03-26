// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

// ITokenFactory is the interface for the layer 2 bridge, which accepts commands from the L1 bridge.
// It needs to be implemented by the L2 bridge as a cross chain only callable method, with the calls coming from CrossChainMessenger.
interface ITokenFactory {
    // onCreateTokenCommand - Will instantiate an ERC20 token contract, owned by the bridge. This token contract allows
    // for minting and burning of assets.
    // crossChainAddress - the address on which the L1 token contract resides
    // name, symbol - ERC20 string representation of the token, for example "Obscuro Token", "OBX"
    function onCreateTokenCommand(
        address crossChainAddress,
        string calldata name,
        string calldata symbol
    ) external;
}
