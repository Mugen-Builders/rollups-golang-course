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

var payload = common.Hex2Bytes("deadbeef")
var msgSender = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

func TestMyApplicationSuite(t *testing.T) {
	suite.Run(t, new(MyApplicationSuite))
}

type MyApplicationSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *MyApplicationSuite) SetupTest() {
	app := new(Application)
	s.tester = rollmelette.NewTester(app)
}

func (s *MyApplicationSuite) TestDelegateCallVoucherSafeTransfer() {
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
	s.Len(result.DelegateCallVouchers, 1)
	s.Equal(safeERC20TransferAddress, result.DelegateCallVouchers[0].Destination)

	abiJSON := `[{"type":"function","name":"safeTransfer","inputs":[{"type":"address"},{"type":"address"},{"type":"uint256"}]}]`
	safeTransferABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := safeTransferABI.Methods["safeTransfer"].Inputs.Unpack(result.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(token, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
	s.Equal(amount, unpacked[2].(*big.Int))
}
