package main

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

var payload = common.Hex2Bytes("deadbeef")

func TestApplicationSuite(t *testing.T) {
	suite.Run(t, new(ApplicationSuite))
}

type ApplicationSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *ApplicationSuite) SetupTest() {
	app := new(Application)
	s.tester = rollmelette.NewTester(app)
}

func (s *ApplicationSuite) TestERC20Deposit() {
	user := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	erc20token := common.HexToAddress("0xa0Ee7A142d267C1f36714E4a8F75612F20a79720")
	result := s.tester.DepositERC20(erc20token, user, big.NewInt(10000), payload)
	s.Len(result.Notices, 4)
	s.Len(result.Vouchers, 1)
	s.Nil(result.Err)

	s.Equal("1 - ERC20 balance of 0x70997970C51812dc3A010C7d01b50e0d17dc79C8: 10000 before transfer to 0x0000000000000000000000000000000000000000", string(result.Notices[0].Payload))

	s.Equal("2 - Balance of 0x0000000000000000000000000000000000000000: 10000 before transfer to 0x70997970C51812dc3A010C7d01b50e0d17dc79C8", string(result.Notices[1].Payload))

	s.Equal("3 - ERC20 balance of 0x70997970C51812dc3A010C7d01b50e0d17dc79C8: 10000 before withdraw", string(result.Notices[2].Payload))

	s.Equal("4 - ERC20 balance of 0x70997970C51812dc3A010C7d01b50e0d17dc79C8: 0 after withdraw", string(result.Notices[3].Payload))

	expectedWithdrawVoucherPayload := make([]byte, 0, 4+32+32)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, user[:]...)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(10000).FillBytes(make([]byte, 32))...)
	s.Equal(expectedWithdrawVoucherPayload, result.Vouchers[0].Payload)
	s.Equal(erc20token, result.Vouchers[0].Destination)
	s.Equal(big.NewInt(0), result.Vouchers[0].Value)
}

func (s *ApplicationSuite) TestEtherDeposit() {
	user := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	amount := big.NewInt(10000)
	result := s.tester.DepositEther(user, amount, payload)
	s.Len(result.Notices, 4)
	s.Len(result.Vouchers, 1)
	s.Nil(result.Err)

	s.Equal("1 - Ether balance of 0x70997970C51812dc3A010C7d01b50e0d17dc79C8: 10000, before transfer to 0x0000000000000000000000000000000000000000", string(result.Notices[0].Payload))

	s.Equal("2 - Balance of 0x0000000000000000000000000000000000000000: 10000 before transfer to 0x70997970C51812dc3A010C7d01b50e0d17dc79C8", string(result.Notices[1].Payload))

	s.Equal("3 - Ether balance of 0x70997970C51812dc3A010C7d01b50e0d17dc79C8: 10000 before withdraw", string(result.Notices[2].Payload))

	s.Equal("4 - Ether balance of 0x70997970C51812dc3A010C7d01b50e0d17dc79C8: 0 after withdraw", string(result.Notices[3].Payload))

	s.Equal(common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e"), result.Vouchers[0].Destination)
	s.Equal(big.NewInt(10000), result.Vouchers[0].Value)
	s.Empty(result.Vouchers[0].Payload)
}
