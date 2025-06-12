// SPDX-License-Identifier: MIT
// Compatible with OpenZeppelin Contracts ^5.0.0
pragma solidity ^0.8.27;

import {IERC20} from "@openzeppelin-contracts-5.2.0/token/ERC20/IERC20.sol";

contract EmergencyWithdraw {
    function emergencyERC20Withdraw(IERC20 token, address to) public {
        token.transfer(to, token.balanceOf(address(this)));
    }
    // TODO: add other tokens support
}
