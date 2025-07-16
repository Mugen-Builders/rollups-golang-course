// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {ERC1155} from "@openzeppelin-contracts-5.2.0/token/ERC1155/ERC1155.sol";
import {ERC1155Burnable} from "@openzeppelin-contracts-5.2.0/token/ERC1155/extensions/ERC1155Burnable.sol";

contract MyERC1155Token is ERC1155, ERC1155Burnable {
    constructor() ERC1155("") {}

    function setURI(string memory newuri) public {
        _setURI(newuri);
    }

    function mint(address account, uint256 id, uint256 amount, bytes memory data) public {
        _mint(account, id, amount, data);
    }

    function mintBatch(address to, uint256[] memory ids, uint256[] memory amounts, bytes memory data) public {
        _mintBatch(to, ids, amounts, data);
    }
}
