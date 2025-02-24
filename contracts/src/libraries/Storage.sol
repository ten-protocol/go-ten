// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0 <0.9.0;

// TODO
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

//    function getUint(bytes32 _slot) internal view returns (uint256 value_) {
//        assembly {
//            value_ := sload(_slot)
//        }
//    }
//
//    function setUint(bytes32 _slot, uint256 _value) internal {
//        assembly {
//            sstore(_slot, _value)
//        }
//    }
//
//    function getBytes32(bytes32 _slot) internal view returns (bytes32 value_) {
//        assembly {
//            value_ := sload(_slot)
//        }
//    }
//
//    function setBytes32(bytes32 _slot, bytes32 _value) internal {
//        assembly {
//            sstore(_slot, _value)
//        }
//    }
//
//    function setBool(bytes32 _slot, bool _value) internal {
//        assembly {
//            sstore(_slot, _value)
//        }
//    }
//
//    function getBool(bytes32 _slot) internal view returns (bool value_) {
//        assembly {
//            value_ := sload(_slot)
//        }
//    }
}