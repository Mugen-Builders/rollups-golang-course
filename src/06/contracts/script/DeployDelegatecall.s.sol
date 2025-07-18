// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script, console} from "forge-std-1.9.7/src/Script.sol";
import {SafeERC721Mint} from "../src/delegatecall/SafeERC721Mint.sol";
import {EmergencyWithdraw} from "../src/delegatecall/EmergencyWithdraw.sol";
import {SafeERC20Transfer} from "../src/delegatecall/SafeERC20Transfer.sol";

contract DeployDelegatecall is Script {
    SafeERC721Mint public safeERC721Mint;
    EmergencyWithdraw public emergencyWithdraw;
    SafeERC20Transfer public safeERC20Transfer;

    function run() public returns (SafeERC721Mint, EmergencyWithdraw, SafeERC20Transfer) {
        vm.startBroadcast();
        safeERC721Mint = new SafeERC721Mint();
        emergencyWithdraw = new EmergencyWithdraw();
        safeERC20Transfer = new SafeERC20Transfer();
        vm.stopBroadcast();

        _saveDeploymentInfo();

        return (safeERC721Mint, emergencyWithdraw, safeERC20Transfer);
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"deployer":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"safeERC721Mint":"',
            vm.toString(address(safeERC721Mint)),
            '"emergencyWithdraw":"',
            vm.toString(address(emergencyWithdraw)),
            '","safeERC20Transfer":"',
            vm.toString(address(safeERC20Transfer)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/deployer.json");
    }
}