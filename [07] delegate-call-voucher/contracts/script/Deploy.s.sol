// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std-1.9.7/src/Script.sol";
import {ERC20Token} from "../src/ERC20Token.sol";
import {SafeERC20Transfer} from "../src/SafeERC20Transfer.sol";

contract Deploy is Script {
    ERC20Token public erc20token;
    SafeERC20Transfer public safeERC20Transfer;

    function setUp() public {}

    function run() public returns (ERC20Token, SafeERC20Transfer) {
        vm.startBroadcast();

        erc20token = new ERC20Token();
        safeERC20Transfer = new SafeERC20Transfer();

        vm.stopBroadcast();

        return (erc20token, safeERC20Transfer);
    }
}
