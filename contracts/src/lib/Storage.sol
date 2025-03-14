// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0 <0.9.0;

library Storage {
    function getAddress(bytes32 _slot) internal view returns (address addr_) {
        assembly {
            addr_ := sload(_slot)
        }
    }

    function setAddress(bytes32 _slot, address _address) internal {
        assembly {
            sstore(_slot, _address)
        }
    }
}