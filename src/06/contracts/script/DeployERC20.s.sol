// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Token} from "../src/token/ERC20/Token.sol";
import {Script, console} from "forge-std-1.9.7/src/Script.sol";

contract DeployERC20 is Script {
    Token public token;

    function run() public returns (Token) {
        vm.startBroadcast();
        token = new Token();
        vm.stopBroadcast();

        _saveDeploymentInfo();

        return token;
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
            '"token":"',
            vm.toString(address(token)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/deployer.json");
    }
}