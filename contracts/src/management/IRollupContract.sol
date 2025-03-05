// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "../common/Structs.sol";

interface IRollupContract {
    event RollupAdded(bytes32 rollupHash, bytes signature);

    function addRollup(Structs.MetaRollup calldata rollup) external;
    function getRollupByHash(bytes32 rollupHash) external view returns (bool, Structs.MetaRollup memory);
}