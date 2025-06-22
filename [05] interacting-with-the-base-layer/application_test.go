package main

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

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

func (s *MyApplicationSuite) TestMintNFT() {
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	to := common.HexToAddress("0x0000000000000000000000000000000000000000")
	uri := "https://example.com"
	input := map[string]interface{}{
		"path": "mint_path",
		"data": map[string]interface{}{
			"token": token,
			"to":    to,
			"uri":   uri,
		},
	}
	payload, err := json.Marshal(input)
	s.Require().NoError(err)

	result := s.tester.Advance(msgSender, payload)
	s.Len(result.Vouchers, 1)
	s.Nil(result.Err)

	s.Equal(common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa"), result.Vouchers[0].Destination)
	s.Equal(big.NewInt(0), result.Vouchers[0].Value)

	abiJSON := `[{"type":"function","name":"safeMint","inputs":[{"type":"address"},{"type":"string"}]}]`
	safeMintABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := safeMintABI.Methods["safeMint"].Inputs.Unpack(result.Vouchers[0].Payload[4:])
	s.Require().NoError(err)

	expectedTo := to
	expectedURI := "https://example.com"
	s.Equal(expectedTo, unpacked[0])
	s.Equal(expectedURI, unpacked[1])

	s.Equal(token, result.Vouchers[0].Destination)
	s.Equal(big.NewInt(0), result.Vouchers[0].Value)
}

func (s *MyApplicationSuite) TestDeployContract() {
	bytecodeHex := "6080604052348015600e575f5ffd5b5060405161028f38038061028f8339818101604052810190602e9190607b565b603b81604060201b60201c565b5060a1565b805f8190555050565b5f5ffd5b5f819050919050565b605d81604d565b81146066575f5ffd5b50565b5f815190506075816056565b92915050565b5f60208284031215608d57608c6049565b5b5f6098848285016069565b91505092915050565b6101e1806100ae5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c80633fb5c1cb146100435780638381f58a1461005f578063d09de08a1461007d575b5f5ffd5b61005d600480360381019061005891906100e4565b610087565b005b610067610090565b604051610074919061011e565b60405180910390f35b610085610095565b005b805f8190555050565b5f5481565b5f5f8154809291906100a690610164565b9190505550565b5f5ffd5b5f819050919050565b6100c3816100b1565b81146100cd575f5ffd5b50565b5f813590506100de816100ba565b92915050565b5f602082840312156100f9576100f86100ad565b5b5f610106848285016100d0565b91505092915050565b610118816100b1565b82525050565b5f6020820190506101315f83018461010f565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61016e826100b1565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036101a05761019f610137565b5b60018201905091905056fea264697066735822122085ac6d83dc5434c93a779a3d0ba32601e836a7655e245b670b88a4f291fbdc7664736f6c634300081e0033"
	bytecode, err := hex.DecodeString(bytecodeHex)
	s.Require().NoError(err)

	uint256Type, err := abi.NewType("uint256", "", nil)
	s.Require().NoError(err)
	
	args := abi.Arguments{{Type: uint256Type}}
	packedArgs, err := args.Pack(big.NewInt(1596))
	s.Require().NoError(err)

	combinedBytecode := append(bytecode, packedArgs...)

	input := map[string]interface{}{
		"path": "deploy_contract",
		"data": map[string]interface{}{
			"deployer": common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa"),
			"bytecode": combinedBytecode,
		},
	}
	payload, err := json.Marshal(input)
	s.Require().NoError(err)

	result := s.tester.Advance(msgSender, payload)
	s.Len(result.Vouchers, 1)
	s.Nil(result.Err)

	s.Equal(common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa"), result.Vouchers[0].Destination)
	s.Equal(big.NewInt(0), result.Vouchers[0].Value)

	abiJSON := `[{"type":"function","name":"deploy","inputs":[{"type":"bytes"}]}]`
	deployABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := deployABI.Methods["deploy"].Inputs.Unpack(result.Vouchers[0].Payload[4:])
	s.Require().NoError(err)

	s.Equal(combinedBytecode, unpacked[0].([]byte))
}
