// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std-1.9.7/src/Script.sol";
import {ERC721Token} from "../src/ERC721.sol";
import {ProxyDeployer} from "../src/ProxyDeployer.sol";

contract Deploy is Script {
    ERC721Token public erc721;
    ProxyDeployer public proxyDeployer;

    function setUp() public {}

    function run() public returns (ERC721Token, ProxyDeployer) {
        vm.startBroadcast();

        erc721 = new ERC721Token();
        proxyDeployer = new ProxyDeployer();

        vm.stopBroadcast();

        return (erc721, proxyDeployer);
    }
}
