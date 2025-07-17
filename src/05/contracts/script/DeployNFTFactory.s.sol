// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {NFT} from "../src/NFT.sol";
import {NFTFactory} from "../src/NFTFactory.sol";
import {Script, console} from "forge-std-1.9.7/src/Script.sol";

contract DeployNFTFactory is Script {
    NFT public nft;
    NFTFactory public nftFactory;

    function setUp() public {}

    function run() public returns (NFT, NFTFactory) {
        vm.startBroadcast();

        nft = new NFT(address(this));
        nftFactory = new NFTFactory();

        vm.stopBroadcast();

        return (nft, nftFactory);
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"deployer":{',
            '"chainId":',
            vm.toString(block.chainid),
            '},"nft":{',
            '"address":"',
            vm.toString(address(nft)),
            '"},"nftFactory":{',
            '"address":"',
            vm.toString(address(nftFactory)),
            '"',
            "}",
            "}}"
        );
        vm.writeJson(deploymentInfo, "./deployments/deployer.json");
    }
}
