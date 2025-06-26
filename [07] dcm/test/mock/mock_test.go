package mock

import (
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/cmd/root"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository/factory"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

func TestDCMSystem(t *testing.T) {
	suite.Run(t, new(DCMSystemSuite))
}

type DCMSystemSuite struct {
	suite.Suite
	Tester                   *rollmelette.Tester
	EmergencyWithdrawAddress common.Address
}

func (s *DCMSystemSuite) SetupTest() {
	repo, err := factory.NewRepositoryFromConnectionString("sqlite://:memory:")
	if err != nil {
		slog.Error("Failed to setup in-memory SQLite database", "error", err)
		os.Exit(1)
	}

	dapp := root.NewDCMSystem(repo)
	s.Tester = rollmelette.NewTester(dapp)
}

func (s *DCMSystemSuite) TestCreateCampaign() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	issuer := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// create issuer user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"issuer"}}`, issuer))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":2,"role":"issuer","address":"%s","created_at":%d}`, issuer, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createCampaignOutput := s.Tester.DepositERC20(collateral, issuer, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))
}

func (s *DCMSystemSuite) TestCloseCampaign() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	anyone := common.HexToAddress("0x0000000000000000000000000000000000000001")
	issuer := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	investor01 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	investor02 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	investor03 := common.HexToAddress("0x0000000000000000000000000000000000000003")
	investor04 := common.HexToAddress("0x0000000000000000000000000000000000000004")
	investor05 := common.HexToAddress("0x0000000000000000000000000000000000000005")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// create issuer user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"issuer"}}`, issuer))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":2,"role":"issuer","address":"%s","created_at":%d}`, issuer, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create investors users
	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":3,"role":"investor","address":"%s","created_at":%d}`, investor01, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":4,"role":"investor","address":"%s","created_at":%d}`, investor02, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor03))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":5,"role":"investor","address":"%s","created_at":%d}`, investor03, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor04))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":6,"role":"investor","address":"%s","created_at":%d}`, investor04, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor05))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":7,"role":"investor","address":"%s","created_at":%d}`, investor05, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createCampaignOutput := s.Tester.DepositERC20(collateral, issuer, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	createOrderOutput := s.Tester.DepositERC20(token, investor01, big.NewInt(60000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"8"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor02, big.NewInt(28000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"4"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor03, big.NewInt(2000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"6"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor04, big.NewInt(5000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"4"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor05, big.NewInt(5500), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	time.Sleep(5 * time.Second)

	closeCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/close", "data":{"issuer":"%s"}}`, issuer))
	closeCampaignOutput := s.Tester.Advance(anyone, closeCampaignInput)
	s.Len(closeCampaignOutput.Notices, 1)

	expectedCloseCampaignOutput := fmt.Sprintf(`campaign closed - {"id":1,"token":"%s","issuer":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		issuer.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)
	s.Equal(expectedCloseCampaignOutput, string(closeCampaignOutput.Notices[0].Payload))

	// Verify final balances after campaign close
	// investor01: deposited 60000, partially accepted 59500, rejected 500
	// investor02: deposited 28000, fully accepted 28000
	// investor03: deposited 2000, fully accepted 2000
	// investor04: deposited 5000, fully accepted 5000
	// investor05: deposited 5500, fully accepted 5500
	// issuer: deposited 10000 collateral, received 100000 from investors

	// Verify investor01 balance (60000 - 59500 = 500 rejected should be returned)
	erc20BalanceInput := []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor01.Hex(), token.Hex()))
	erc20BalanceOutput := s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"500"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor02 balance (28000 - 28000 = 0)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor02.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor03 balance (2000 - 2000 = 0)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor03.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor04 balance (5000 - 5000 = 0)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor04.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor05 balance (5500 - 5500 = 0)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor05.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify issuer balance (should have received 100000 from investors)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, issuer.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"100000"`, string(erc20BalanceOutput.Reports[0].Payload))
}

func (s *DCMSystemSuite) TestSettleCampaign() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	anyone := common.HexToAddress("0x0000000000000000000000000000000000000001")
	issuer := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	investor01 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	investor02 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	investor03 := common.HexToAddress("0x0000000000000000000000000000000000000003")
	investor04 := common.HexToAddress("0x0000000000000000000000000000000000000004")
	investor05 := common.HexToAddress("0x0000000000000000000000000000000000000005")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// create issuer user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"issuer"}}`, issuer))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":2,"role":"issuer","address":"%s","created_at":%d}`, issuer, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create investors users
	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":3,"role":"investor","address":"%s","created_at":%d}`, investor01, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":4,"role":"investor","address":"%s","created_at":%d}`, investor02, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor03))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":5,"role":"investor","address":"%s","created_at":%d}`, investor03, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor04))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":6,"role":"investor","address":"%s","created_at":%d}`, investor04, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor05))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":7,"role":"investor","address":"%s","created_at":%d}`, investor05, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createCampaignOutput := s.Tester.DepositERC20(collateral, issuer, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	createOrderOutput := s.Tester.DepositERC20(token, investor01, big.NewInt(60000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"8"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor02, big.NewInt(28000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"4"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor03, big.NewInt(2000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"6"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor04, big.NewInt(5000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"4"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor05, big.NewInt(5500), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	time.Sleep(5 * time.Second)

	closeCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/close", "data":{"issuer":"%s"}}`, issuer))
	closeCampaignOutput := s.Tester.Advance(anyone, closeCampaignInput)
	s.Len(closeCampaignOutput.Notices, 1)

	expectedCloseCampaignOutput := fmt.Sprintf(`campaign closed - {"id":1,"token":"%s","issuer":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		issuer.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)
	s.Equal(expectedCloseCampaignOutput, string(closeCampaignOutput.Notices[0].Payload))

	// Withdraw raised amount
	withdrawRaisedAmountInput := []byte(fmt.Sprintf(`{"path":"user/erc20-withdraw","data":{"token":"%s","amount":"100000"}}`, token.Hex()))
	withdrawRaisedAmountOutput := s.Tester.Advance(issuer, withdrawRaisedAmountInput)
	s.Len(withdrawRaisedAmountOutput.Notices, 1)

	expectedWithdrawRaisedAmountOutput := fmt.Sprintf(`ERC20 withdrawn - token: %s, amount: 100000, user: %s`, token.Hex(), issuer.Hex())
	s.Equal(expectedWithdrawRaisedAmountOutput, string(withdrawRaisedAmountOutput.Notices[0].Payload))

	time.Sleep(5 * time.Second)

	settleCampaignInput := []byte(`{"path":"campaign/issuer/settle", "data":{"campaign_id":1}}`)
	settleCampaignOutput := s.Tester.DepositERC20(token, issuer, big.NewInt(108195), settleCampaignInput)
	s.Len(settleCampaignOutput.Notices, 1)

	settledAt := baseTime + 10 // baseTime

	expectedSettleCampaignOutput := fmt.Sprintf(`campaign settled - {"id":1,"token":"%s","issuer":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"settled","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		issuer.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, settledAt, // Order 1
		investor02.Hex(), baseTime, settledAt, // Order 2
		investor03.Hex(), baseTime, settledAt, // Order 3
		investor04.Hex(), baseTime, settledAt, // Order 4
		investor05.Hex(), baseTime, settledAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, settledAt,
	)
	s.Equal(expectedSettleCampaignOutput, string(settleCampaignOutput.Notices[0].Payload))

	// Verify final balances after campaign settlement
	// investor01: should receive 59500 + (59500 * 9% = 64855) = 64855
	// investor02: should receive 28000 + (28000 * 8% = 2240) = 30240
	// investor03: should receive 2000 + (2000 * 4% = 80) = 2080
	// investor04: should receive 5000 + (5000 * 6% = 300) = 5300
	// investor05: should receive 5500 + (5500 * 4% = 220) = 5720
	// issuer: paid 108195 to settle the campaign

	// Verify investor01 balance (received 64855 + rejected order amount = 65355)
	erc20BalanceInput := []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor01.Hex(), token.Hex()))
	erc20BalanceOutput := s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"65355"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor02 balance (received 30240)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor02.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"30240"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor03 balance (received 2080)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor03.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"2080"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor04 balance (received 5300)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor04.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"5300"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor05 balance (received 5720)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor05.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"5720"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify issuer balance (had 100000, paid 108195, so should be -8195)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, issuer.Hex(), token.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))
}

func (s *DCMSystemSuite) TestExecuteCampaignCollateral() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	anyone := common.HexToAddress("0x0000000000000000000000000000000000000001")
	issuer := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	investor01 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	investor02 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	investor03 := common.HexToAddress("0x0000000000000000000000000000000000000003")
	investor04 := common.HexToAddress("0x0000000000000000000000000000000000000004")
	investor05 := common.HexToAddress("0x0000000000000000000000000000000000000005")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// create issuer user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"issuer"}}`, issuer))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":2,"role":"issuer","address":"%s","created_at":%d}`, issuer, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create investors users
	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":3,"role":"investor","address":"%s","created_at":%d}`, investor01, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":4,"role":"investor","address":"%s","created_at":%d}`, investor02, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor03))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":5,"role":"investor","address":"%s","created_at":%d}`, investor03, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor04))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":6,"role":"investor","address":"%s","created_at":%d}`, investor04, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor05))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":7,"role":"investor","address":"%s","created_at":%d}`, investor05, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createCampaignOutput := s.Tester.DepositERC20(collateral, issuer, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	createOrderOutput := s.Tester.DepositERC20(token, investor01, big.NewInt(60000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"8"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor02, big.NewInt(28000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"4"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor03, big.NewInt(2000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"6"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor04, big.NewInt(5000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"4"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor05, big.NewInt(5500), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	time.Sleep(5 * time.Second)

	closeCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/close", "data":{"issuer":"%s"}}`, issuer))
	closeCampaignOutput := s.Tester.Advance(anyone, closeCampaignInput)
	s.Len(closeCampaignOutput.Notices, 1)

	expectedCloseCampaignOutput := fmt.Sprintf(`campaign closed - {"id":1,"token":"%s","issuer":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		issuer.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)
	s.Equal(expectedCloseCampaignOutput, string(closeCampaignOutput.Notices[0].Payload))

	// Withdraw raised amount
	withdrawRaisedAmountInput := []byte(fmt.Sprintf(`{"path":"user/erc20-withdraw","data":{"token":"%s","amount":"100000"}}`, token.Hex()))
	withdrawRaisedAmountOutput := s.Tester.Advance(issuer, withdrawRaisedAmountInput)
	s.Len(withdrawRaisedAmountOutput.Notices, 1)

	expectedWithdrawRaisedAmountOutput := fmt.Sprintf(`ERC20 withdrawn - token: %s, amount: 100000, user: %s`, token.Hex(), issuer.Hex())
	s.Equal(expectedWithdrawRaisedAmountOutput, string(withdrawRaisedAmountOutput.Notices[0].Payload))

	findCampaignByIdInput := []byte(fmt.Sprintf(`{"path":"campaign/id", "data":{"id":1}}`))

	findCampaignByIdOutput := s.Tester.Inspect(findCampaignByIdInput)
	s.Len(findCampaignByIdOutput.Reports, 1)

	expectedFindCampaignByIssuerOutput := fmt.Sprintf(`[{"id":1,"token":"%s","issuer":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}]`,
		token.Hex(),
		issuer.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)

	findCampaignsByIssuerInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer", "data":{"issuer":"%s"}}`, issuer))

	findCampaignsByIssuerOutput := s.Tester.Inspect(findCampaignsByIssuerInput)
	s.Len(findCampaignsByIssuerOutput.Reports, 1)
	s.Equal(expectedFindCampaignByIssuerOutput, string(findCampaignsByIssuerOutput.Reports[0].Payload))

	time.Sleep(6 * time.Second)

	executeCampaignCollateralInput := []byte(fmt.Sprintf(`{"path":"campaign/execute-collateral", "data":{"campaign_id":1}}`))
	executeCampaignCollateralOutput := s.Tester.Advance(issuer, executeCampaignCollateralInput)
	s.Len(executeCampaignCollateralOutput.Notices, 1)

	collateralExecutedAt := baseTime + 11 // baseTime

	expectedExecuteCampaignCollateralOutput := fmt.Sprintf(`campaign collateral executed - {"campaign_id":1,"token":"%s","issuer":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"collateral_executed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		issuer.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, collateralExecutedAt, // Order 1
		investor02.Hex(), baseTime, collateralExecutedAt, // Order 2
		investor03.Hex(), baseTime, collateralExecutedAt, // Order 3
		investor04.Hex(), baseTime, collateralExecutedAt, // Order 4
		investor05.Hex(), baseTime, collateralExecutedAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, collateralExecutedAt,
	)
	s.Equal(expectedExecuteCampaignCollateralOutput, string(executeCampaignCollateralOutput.Notices[0].Payload))

	// Verify final balances after campaign collateral execution
	// The collateral (10000) is distributed proportionally to accepted orders based on their final value
	// Total final value = 59500*1.09 + 28000*1.08 + 2000*1.04 + 5000*1.06 + 5500*1.04 = 64855 + 30240 + 2080 + 5300 + 5720 = 108195
	// investor01: 64855/108195 * 10000 = 5994 (rounded down)
	// investor02: 30240/108195 * 10000 = 2794 (rounded down)
	// investor03: 2080/108195 * 10000 = 192 (rounded down)
	// investor04: 5300/108195 * 10000 = 489 (rounded down)
	// investor05: 5720/108195 * 10000 = 528 (rounded down)
	// Total distributed: 5994 + 2794 + 192 + 489 + 528 = 9997
	// Remaining: 10000 - 9997 = 3 tokens remain in the application (not distributed)
	// Final distribution:
	// investor01: 5994, investor02: 2794, investor03: 192, investor04: 489, investor05: 528
	// issuer: no additional deposit, just execution of existing collateral

	// Verify investor01 balance (received 5994 collateral)
	erc20BalanceInput := []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor01.Hex(), collateral.Hex()))
	erc20BalanceOutput := s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"5994"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor02 balance (received 2794 collateral)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor02.Hex(), collateral.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"2794"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor03 balance (received 192 collateral)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor03.Hex(), collateral.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"192"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor04 balance (received 489 collateral)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor04.Hex(), collateral.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"489"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify investor05 balance (received 528 collateral)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, investor05.Hex(), collateral.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"528"`, string(erc20BalanceOutput.Reports[0].Payload))

	// Verify issuer balance (no additional deposit, just execution of existing collateral)
	erc20BalanceInput = []byte(fmt.Sprintf(`{"path":"user/erc20-balance","data":{"address":"%s","token":"%s"}}`, issuer.Hex(), collateral.Hex()))
	erc20BalanceOutput = s.Tester.Inspect(erc20BalanceInput)
	s.Len(erc20BalanceOutput.Reports, 1)
	s.Equal(`"0"`, string(erc20BalanceOutput.Reports[0].Payload))
}

func (s *DCMSystemSuite) TestFindAllCampaigns() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	issuer := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// create issuer user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"issuer"}}`, issuer))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":2,"role":"issuer","address":"%s","created_at":%d}`, issuer, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createCampaignOutput := s.Tester.DepositERC20(collateral, issuer, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	findAllCampaignsInput := []byte(`{"path":"campaign"}`)

	findAllCampaignsOutput := s.Tester.Inspect(findAllCampaignsInput)
	s.Len(findAllCampaignsOutput.Reports, 1)

	expectedFindAllCampaignsOutput := fmt.Sprintf(`[{"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}]`, baseTime, closesAt, maturityAt)
	s.Equal(expectedFindAllCampaignsOutput, string(findAllCampaignsOutput.Reports[0].Payload))
}

func (s *DCMSystemSuite) TestFindCampaignById() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	issuer := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// create issuer user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"issuer"}}`, issuer))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":2,"role":"issuer","address":"%s","created_at":%d}`, issuer, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createCampaignOutput := s.Tester.DepositERC20(collateral, issuer, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	findCampaignByIdInput := []byte(fmt.Sprintf(`{"path":"campaign/id", "data":{"id":1}}`))

	findCampaignByIdOutput := s.Tester.Inspect(findCampaignByIdInput)
	s.Len(findCampaignByIdOutput.Reports, 1)

	expectedFindCampaignByIdOutput := fmt.Sprintf(`{"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedFindCampaignByIdOutput, string(findCampaignByIdOutput.Reports[0].Payload))
}

func (s *DCMSystemSuite) TestFindCampaignsByIssuer() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	issuer := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// create issuer user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"issuer"}}`, issuer))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":2,"role":"issuer","address":"%s","created_at":%d}`, issuer, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createCampaignOutput := s.Tester.DepositERC20(collateral, issuer, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	findCampaignsByIssuerInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer", "data":{"issuer":"%s"}}`, issuer))

	findCampaignsByIssuerOutput := s.Tester.Inspect(findCampaignsByIssuerInput)
	s.Len(findCampaignsByIssuerOutput.Reports, 1)

	expectedFindCampaignsByIssuerOutput := fmt.Sprintf(`[{"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}]`, baseTime, closesAt, maturityAt)
	s.Equal(expectedFindCampaignsByIssuerOutput, string(findCampaignsByIssuerOutput.Reports[0].Payload))
}

func (s *DCMSystemSuite) TestFindCampaignsByInvestor() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	anyone := common.HexToAddress("0x0000000000000000000000000000000000000001")
	issuer := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	investor01 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	investor02 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	investor03 := common.HexToAddress("0x0000000000000000000000000000000000000003")
	investor04 := common.HexToAddress("0x0000000000000000000000000000000000000004")
	investor05 := common.HexToAddress("0x0000000000000000000000000000000000000005")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// create issuer user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"issuer"}}`, issuer))
	createUserOutput := s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput := fmt.Sprintf(`user created - {"id":2,"role":"issuer","address":"%s","created_at":%d}`, issuer, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create investors users
	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":3,"role":"investor","address":"%s","created_at":%d}`, investor01, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":4,"role":"investor","address":"%s","created_at":%d}`, investor02, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor03))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":5,"role":"investor","address":"%s","created_at":%d}`, investor03, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor04))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":6,"role":"investor","address":"%s","created_at":%d}`, investor04, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor05))
	createUserOutput = s.Tester.Advance(admin, createUserInput)
	s.Len(createUserOutput.Notices, 1)

	expectedCreateUserOutput = fmt.Sprintf(`user created - {"id":7,"role":"investor","address":"%s","created_at":%d}`, investor05, baseTime)
	s.Equal(expectedCreateUserOutput, string(createUserOutput.Notices[0].Payload))

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createCampaignOutput := s.Tester.DepositERC20(collateral, issuer, big.NewInt(10000), createCampaignInput)
	s.Len(createCampaignOutput.Notices, 1)

	expectedCreateCampaignOutput := fmt.Sprintf(`campaign created - {"id":1,"token":"0x0000000000000000000000000000000000000009","issuer":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateCampaignOutput, string(createCampaignOutput.Notices[0].Payload))

	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	createOrderOutput := s.Tester.DepositERC20(token, investor01, big.NewInt(60000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"8"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor02, big.NewInt(28000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"4"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor03, big.NewInt(2000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"6"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor04, big.NewInt(5000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"4"}}`)
	createOrderOutput = s.Tester.DepositERC20(token, investor05, big.NewInt(5500), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	time.Sleep(5 * time.Second)

	closeCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/close", "data":{"issuer":"%s"}}`, issuer))
	closeCampaignOutput := s.Tester.Advance(anyone, closeCampaignInput)
	s.Len(closeCampaignOutput.Notices, 1)

	expectedCloseCampaignOutput := fmt.Sprintf(`campaign closed - {"id":1,"token":"%s","issuer":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		issuer.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)
	s.Equal(expectedCloseCampaignOutput, string(closeCampaignOutput.Notices[0].Payload))

	// Withdraw raised amount
	withdrawRaisedAmountInput := []byte(fmt.Sprintf(`{"path":"user/erc20-withdraw","data":{"token":"%s","amount":"100000"}}`, token.Hex()))
	withdrawRaisedAmountOutput := s.Tester.Advance(issuer, withdrawRaisedAmountInput)
	s.Len(withdrawRaisedAmountOutput.Notices, 1)

	expectedWithdrawRaisedAmountOutput := fmt.Sprintf(`ERC20 withdrawn - token: %s, amount: 100000, user: %s`, token.Hex(), issuer.Hex())
	s.Equal(expectedWithdrawRaisedAmountOutput, string(withdrawRaisedAmountOutput.Notices[0].Payload))

	expectedFindCampaignByIssuerOutput := fmt.Sprintf(`[{"id":1,"token":"%s","issuer":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}]`,
		token.Hex(),
		issuer.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)

	findCampaignsByIssuerInput := []byte(fmt.Sprintf(`{"path":"campaign/issuer", "data":{"issuer":"%s"}}`, issuer))

	findCampaignsByIssuerOutput := s.Tester.Inspect(findCampaignsByIssuerInput)
	s.Len(findCampaignsByIssuerOutput.Reports, 1)
	s.Equal(expectedFindCampaignByIssuerOutput, string(findCampaignsByIssuerOutput.Reports[0].Payload))
}

func (s *DCMSystemSuite) TestEmergencyERC20Withdraw() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	to := common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955")
	emergencyWithdrawAddress := common.HexToAddress("0x0000000000000000000000000000000000000001")

	// Emergency ERC20 withdraw
	emergencyERC20WithdrawInput := []byte(fmt.Sprintf(`{"path":"user/admin/emergency-erc20-withdraw","data":{"to":"%s","token":"%s","emergency_withdraw_address":"%s"}}`, to.Hex(), token.Hex(), emergencyWithdrawAddress.Hex()))
	emergencyERC20WithdrawOutput := s.Tester.Advance(admin, emergencyERC20WithdrawInput)
	s.Len(emergencyERC20WithdrawOutput.DelegateCallVouchers, 1)
	s.Equal(emergencyWithdrawAddress, emergencyERC20WithdrawOutput.DelegateCallVouchers[0].Destination)

	// Verify the delegate call voucher payload
	abiJSON := `[{
		"type":"function",
		"name":"emergencyERC20Withdraw",
		"inputs":[
            {"type":"address"},
			{"type":"address"},
			{"type":"address"}
		]
	}]`
	emergencyWithdrawABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := emergencyWithdrawABI.Methods["emergencyERC20Withdraw"].Inputs.Unpack(emergencyERC20WithdrawOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(admin, unpacked[0].(common.Address))
	s.Equal(token, unpacked[1].(common.Address))
	s.Equal(to, unpacked[2].(common.Address))
}

func (s *DCMSystemSuite) TestEmergencyEtherWithdraw() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	to := common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955")
	emergencyWithdrawAddress := common.HexToAddress("0x0000000000000000000000000000000000000001")

	// Emergency ETH withdraw
	emergencyEtherWithdrawInput := []byte(fmt.Sprintf(`{"path":"user/admin/emergency-ether-withdraw","data":{"to":"%s","emergency_withdraw_address":"%s"}}`, to.Hex(), emergencyWithdrawAddress.Hex()))
	emergencyEtherWithdrawOutput := s.Tester.Advance(admin, emergencyEtherWithdrawInput)
	s.Len(emergencyEtherWithdrawOutput.DelegateCallVouchers, 1)
	s.Equal(emergencyWithdrawAddress, emergencyEtherWithdrawOutput.DelegateCallVouchers[0].Destination)

	// Verify the delegate call voucher payload
	abiJSON := `[{
		"type":"function",
		"name":"emergencyETHWithdraw",
		"inputs":[
			{"type":"address"},
			{"type":"address"}
		]
	}]`
	emergencyETHWithdrawABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := emergencyETHWithdrawABI.Methods["emergencyETHWithdraw"].Inputs.Unpack(emergencyEtherWithdrawOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(admin, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
}
