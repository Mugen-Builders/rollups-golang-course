package main

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

func TestDelegateCallVoucherExample(t *testing.T) {
	suite.Run(t, new(DelegateCallVoucherExample))
}

type DelegateCallVoucherExample struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *DelegateCallVoucherExample) SetupTest() {
	app := new(Application)
	s.tester = rollmelette.NewTester(app)
}

func (s *DelegateCallVoucherExample) TestDelegateCallVoucherSafeTransfer() {
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	to := common.HexToAddress("0x0000000000000000000000000000000000000000")
	amount := big.NewInt(10000)

	input := map[string]interface{}{
		"path": "safe_transfer",
		"data": map[string]interface{}{
			"token":  token,
			"to":     to,
			"amount": amount,
		},
	}
	payload, err := json.Marshal(input)
	s.Require().NoError(err)

	result := s.tester.DepositERC20(token, to, amount, payload)
	s.Nil(result.Err)
	s.Len(result.DelegateCallVouchers, 2)
	s.Equal(safeERC20TransferAddress, result.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"safeTransfer",
		"inputs":[
			{"type":"address"},
			{"type":"address"},
			{"type":"uint256"}
		]
	},
	{
		"type":"function",
		"name":"safeTransferTargeted",
		"inputs":[
			{"type":"address"},
			{"type":"address"},
			{"type":"address"},
			{"type":"uint256"}
		]
	}]`
	safeTransferABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	halfAmount := new(big.Int).Div(amount, big.NewInt(2))
	unpacked, err := safeTransferABI.Methods["safeTransfer"].Inputs.Unpack(result.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(token, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
	s.Equal(halfAmount, unpacked[2].(*big.Int))

	unpackedTargeted, err := safeTransferABI.Methods["safeTransferTargeted"].Inputs.Unpack(result.DelegateCallVouchers[1].Payload[4:])
	s.Require().NoError(err)
	s.Equal(token, unpackedTargeted[0].(common.Address))
	s.Equal(anyone, unpackedTargeted[1].(common.Address))
	s.Equal(to, unpackedTargeted[2].(common.Address))
	s.Equal(halfAmount, unpackedTargeted[3].(*big.Int))
}
