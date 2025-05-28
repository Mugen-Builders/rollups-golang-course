// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std-1.9.7/src/Script.sol";
import {MyERC20Token} from "../src/ERC20.sol";

contract Deploy is Script {
    MyERC20Token public erc20;

    function setUp() public {}

    function run() public returns (MyERC20Token) {
        vm.startBroadcast();
        erc20 = new MyERC20Token();
        vm.stopBroadcast();

        return (erc20);
    }
}
