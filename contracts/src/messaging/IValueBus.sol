// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

interface IValueBus {
    function sendValue() external;
    function retrieveValue() external;
}