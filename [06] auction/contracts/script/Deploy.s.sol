// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std-1.9.7/src/Script.sol";
import {ERC20Token} from "../src/ERC20Token.sol";

contract Deploy is Script {
    ERC20Token public erc20token;

    function setUp() public {}

    function run() public returns (ERC20Token) {
        vm.startBroadcast();

        erc20token = new ERC20Token();

        vm.stopBroadcast();

        return (erc20token);
    }
}
