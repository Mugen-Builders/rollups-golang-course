// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Test} from "forge-std-1.9.7/src/Test.sol";
import {NFT} from "../src/token/ERC721/NFT.sol";
import {MockApplication} from "./mock/MockApplication.sol";
import {NFTFactory} from "../src/token/ERC721/NFTFactory.sol";
import {Outputs} from "cartesi-rollups-contracts-2.0.0/src/common/Outputs.sol";
import {SafeERC721Mint, IERC721Mintable} from "../src/delegatecall/SafeERC721Mint.sol";

contract NFTTest is Test {
    NFTFactory public nftFactory;
    SafeERC721Mint public safeERC721Mint;
    MockApplication public mockApplication;

    address public user;

    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
    event NFTDeployed(address indexed nft, bytes32 salt);

    function setUp() public {
        user = makeAddr("user");

        nftFactory = new NFTFactory();
        safeERC721Mint = new SafeERC721Mint();
        mockApplication = new MockApplication();
    }

    function test_DeterministicDeploymentOfNFTViaFactoryThroughVoucherExecution() public {
        string memory symbol = "MTK";
        string memory name = "MyToken";
        bytes32 salt = keccak256("test-salt");

        bytes memory encodedDeployTx = abi.encodeCall(NFTFactory.newNFT, (address(mockApplication), salt, name, symbol));
        bytes memory voucher = abi.encodeCall(Outputs.Voucher, (address(nftFactory), 0, encodedDeployTx));

        address predictedAddress = nftFactory.computeAddress(address(mockApplication), salt, name, symbol);

        vm.expectEmit(true, true, false, true);
        emit NFTDeployed(predictedAddress, salt);
        mockApplication.executeOutput(voucher);

        assertEq(NFT(predictedAddress).name(), name);
        assertEq(NFT(predictedAddress).symbol(), symbol);
    }

    function test_MintNFTThroughDelegatecallVoucher() public {
        string memory symbol = "MTK";
        string memory name = "MyToken";
        bytes32 salt = keccak256("test-salt");

        bytes memory encodedDeployTx = abi.encodeCall(NFTFactory.newNFT, (address(mockApplication), salt, name, symbol));
        bytes memory voucher = abi.encodeCall(Outputs.Voucher, (address(nftFactory), 0, encodedDeployTx));

        address predictedAddress = nftFactory.computeAddress(address(mockApplication), salt, name, symbol);

        vm.expectEmit(true, true, false, true);
        emit NFTDeployed(predictedAddress, salt);
        mockApplication.executeOutput(voucher);

        bytes memory encodedMintTx =
            abi.encodeCall(SafeERC721Mint.safeMint, (IERC721Mintable(predictedAddress), user, "ipfs://test-uri"));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(safeERC721Mint), encodedMintTx));

        vm.expectEmit(true, true, false, true);
        emit Transfer(address(0), user, 0);
        mockApplication.executeOutput(delegateCallVoucher);
        assertEq(NFT(predictedAddress).ownerOf(0), user);
        assertEq(NFT(predictedAddress).tokenURI(0), "ipfs://test-uri");
        assertEq(NFT(predictedAddress).name(), name);
        assertEq(NFT(predictedAddress).symbol(), symbol);
    }
}