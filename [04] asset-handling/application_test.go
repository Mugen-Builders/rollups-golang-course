package main

import (
	"testing"

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

func (s *MyApplicationSuite) TestAdvance() {
	result := s.tester.Advance(msgSender, payload)
	s.Nil(result.Err)
}

func (s *MyApplicationSuite) TestInspect() {
	result := s.tester.Inspect(payload)
	s.Nil(result.Err)
}