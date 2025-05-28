// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

contract ProxyDeployer {
    event Deploy(address);

    receive() external payable {}

    function deploy(bytes memory _code)
        external
        payable
        returns (address addr)
    {
        assembly {
            // create(v, p, n)
            // v = amount of ETH to send
            // p = pointer in memory to start of code
            // n = size of code
            addr := create(callvalue(), add(_code, 0x20), mload(_code))
        }
        // return address 0 on error
        require(addr != address(0), "deploy failed");

        emit Deploy(addr);
    }
}