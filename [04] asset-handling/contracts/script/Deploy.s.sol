// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std-1.9.7/src/Script.sol";
import {MyERC20Token} from "../src/ERC20.sol";
import {MyERC721Token} from "../src/ERC721.sol";
import {MyERC1155Token} from "../src/ERC1155.sol";
import {SafeERC20Transfer} from "../src/SafeERC20Transfer.sol";

contract Deploy is Script {
    MyERC20Token public erc20;
    MyERC721Token public erc721;
    MyERC1155Token public erc1155;
    SafeERC20Transfer public simpleERC20transfer;

    function setUp() public {}

    function run() public returns (MyERC20Token, MyERC721Token, MyERC1155Token, SafeERC20Transfer) {
        vm.startBroadcast();

        erc20 = new MyERC20Token();
        erc721 = new MyERC721Token();
        erc1155 = new MyERC1155Token();
        simpleERC20transfer = new SafeERC20Transfer();

        vm.stopBroadcast();

        return (erc20, erc721, erc1155, simpleERC20transfer);
    }
}
