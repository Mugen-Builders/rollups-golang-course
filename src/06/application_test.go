package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

func TestDelegateCallVoucher(t *testing.T) {
	suite.Run(t, new(DelegateCallVoucher))
}

type DelegateCallVoucher struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *DelegateCallVoucher) SetupTest() {
	app := new(Application)
	s.tester = rollmelette.NewTester(app)
}

func (s *DelegateCallVoucher) TestVoucherDeployNFT() {
	application := common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e")

	deployNFTInput := []byte(`{"path":"deploy_nft","data":{"name":"MyToken","symbol":"MTK"}}`)
	deployNFTOutput := s.tester.Advance(common.HexToAddress("0x0000000000000000000000000000000000000000"), deployNFTInput)
	s.Nil(deployNFTOutput.Err)
	s.Len(deployNFTOutput.Vouchers, 1)
	s.Equal(nftFactoryAddress, deployNFTOutput.Vouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"newNFT",
		"inputs":[
			{"type":"address"},
			{"type":"bytes32"}
		]
	}]`
	newNFTABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := newNFTABI.Methods["newNFT"].Inputs.Unpack(deployNFTOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(application, unpacked[0].(common.Address))
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(0)), common.BytesToHash(saltBytes[:]))
}

func (s *DelegateCallVoucher) TestDelegateCallVoucherMintNFT() {
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")
	uri := "https://example.com/token/1"

	mintNFTInput := []byte(fmt.Sprintf(`{"path":"mint_nft","data":{"to":"%s","uri":"%s"}}`, to, uri))
	mintNFTOutput := s.tester.Advance(to, mintNFTInput)
	s.Nil(mintNFTOutput.Err)
	s.Len(mintNFTOutput.DelegateCallVouchers, 1)
	s.Equal(safeERC721MintAddress, mintNFTOutput.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"safeMint",
		"inputs":[
			{"type":"address"},
			{"type":"address"},
			{"type":"string"}
		]
	}]`
	safeMintABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := safeMintABI.Methods["safeMint"].Inputs.Unpack(mintNFTOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(nftAddress, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
	s.Equal(uri, unpacked[2].(string))
}

func (s *DelegateCallVoucher) TestDelegateCallVoucherSafeTransfer() {
	amount := big.NewInt(10000)
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

	safeTransferInput := []byte(fmt.Sprintf(`{"path":"safe_transfer","data":{"token":"%s","to":"%s","amount":"%s"}}`, token, to, amount))
	safeTransferOutput := s.tester.Advance(to, safeTransferInput)
	s.Nil(safeTransferOutput.Err)
	s.Len(safeTransferOutput.DelegateCallVouchers, 1)
	s.Equal(safeERC20TransferAddress, safeTransferOutput.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"safeTransfer",
		"inputs":[
			{"type":"address"},
			{"type":"address"},
			{"type":"uint256"}
		]
	}]`
	safeTransferABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := safeTransferABI.Methods["safeTransfer"].Inputs.Unpack(safeTransferOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(token, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
	s.Equal(amount, unpacked[2].(*big.Int))
}

func (s *DelegateCallVoucher) TestDelegateCallVoucherSafeTransferTargeted() {
	amount := big.NewInt(10000)
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

	safeTransferTargetedInput := []byte(fmt.Sprintf(`{"path":"safe_transfer_targeted","data":{"token":"%s","to":"%s","amount":"%s"}}`, token, to, amount))
	safeTransferTargetedOutput := s.tester.Advance(to, safeTransferTargetedInput)
	s.Nil(safeTransferTargetedOutput.Err)
	s.Len(safeTransferTargetedOutput.DelegateCallVouchers, 1)
	s.Equal(safeERC20TransferAddress, safeTransferTargetedOutput.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"safeTransferTargeted",
		"inputs":[
			{"type":"address"},
			{"type":"address"},
			{"type":"address"},
			{"type":"uint256"}
		]
	}]`
	safeTransferTargetedABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := safeTransferTargetedABI.Methods["safeTransferTargeted"].Inputs.Unpack(safeTransferTargetedOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(token, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
	s.Equal(to, unpacked[2].(common.Address))
	s.Equal(amount, unpacked[3].(*big.Int))
}

func (s *DelegateCallVoucher) TestDelegateCallVoucherEmergencyERC20Withdraw() {
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

	emergencyERC20WithdrawInput := []byte(fmt.Sprintf(`{"path":"emergency_erc20_withdraw","data":{"token":"%s","to":"%s"}}`, token, to))
	emergencyERC20WithdrawOutput := s.tester.Advance(to, emergencyERC20WithdrawInput)
	s.Nil(emergencyERC20WithdrawOutput.Err)
	s.Len(emergencyERC20WithdrawOutput.DelegateCallVouchers, 1)
	s.Equal(emergencyWithdrawAddress, emergencyERC20WithdrawOutput.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"emergencyERC20Withdraw",
		"inputs":[
			{"type":"address"},
			{"type":"address"}
		]
	}]`
	emergencyWithdrawABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := emergencyWithdrawABI.Methods["emergencyERC20Withdraw"].Inputs.Unpack(emergencyERC20WithdrawOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(token, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
}

func (s *DelegateCallVoucher) TestDelegateCallVoucherEmergencyETHWithdraw() {
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")

	emergencyETHWithdrawInput := []byte(fmt.Sprintf(`{"path":"emergency_eth_withdraw","data":{"to":"%s"}}`, to))
	emergencyETHWithdrawOutput := s.tester.Advance(to, emergencyETHWithdrawInput)
	s.Nil(emergencyETHWithdrawOutput.Err)
	s.Len(emergencyETHWithdrawOutput.DelegateCallVouchers, 1)
	s.Equal(emergencyWithdrawAddress, emergencyETHWithdrawOutput.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"emergencyETHWithdraw",
		"inputs":[
			{"type":"address"}
		]
	}]`
	emergencyETHWithdrawABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := emergencyETHWithdrawABI.Methods["emergencyETHWithdraw"].Inputs.Unpack(emergencyETHWithdrawOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(to, unpacked[0].(common.Address))
}

func (s *DelegateCallVoucher) TestInspectContracts() {
	contractsInput := []byte(`{"path":"contracts"}`)
	contractsOutput := s.tester.Inspect(contractsInput)
	s.Nil(contractsOutput.Err)
	s.Len(contractsOutput.Reports, 1)

	expectedContractsOutput := fmt.Sprintf(`[{"name":"NFT","address":"%s"},{"name":"NFTFactory","address":"%s"},{"name":"EmergencyWithdraw","address":"%s"},{"name":"SafeERC20Transfer","address":"%s"}]`, nftAddress, nftFactoryAddress, emergencyWithdrawAddress, safeERC20TransferAddress)
	s.Equal(string(contractsOutput.Reports[0].Payload), expectedContractsOutput)
}
