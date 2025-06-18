package test

import (
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/cartesi-golang-series/auction/cmd/root"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository/factory"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

func TestAuctionSystem(t *testing.T) {
	suite.Run(t, new(AuctionSystemSuite))
}

type AuctionSystemSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *AuctionSystemSuite) SetupTest() {
	repo, err := factory.NewRepositoryFromConnectionString("sqlite://:memory:")
	if err != nil {
		slog.Error("Failed to setup in-memory SQLite database", "error", err)
		os.Exit(1)
	}
	dapp := root.NewAuctionSystem(repo)
	s.tester = rollmelette.NewTester(dapp)
}

func (s *AuctionSystemSuite) TestCreateAuction() {
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	createAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createAuctionResult := s.tester.DepositERC20(collateral, creator, big.NewInt(10000), createAuctionInput)
	s.Len(createAuctionResult.Notices, 1)

	expectedCreateAuctionOutput := fmt.Sprintf(`auction created - {"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateAuctionOutput, string(createAuctionResult.Notices[0].Payload))
}

func (s *AuctionSystemSuite) TestCloseAuction() {
	anyone := common.HexToAddress("0x0000000000000000000000000000000000000001")
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
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

	createAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createAuctionResult := s.tester.DepositERC20(collateral, creator, big.NewInt(10000), createAuctionInput)
	s.Len(createAuctionResult.Notices, 1)

	expectedCreateAuctionOutput := fmt.Sprintf(`auction created - {"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateAuctionOutput, string(createAuctionResult.Notices[0].Payload))

	createOrderInput := []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"9"}}`)
	createOrderResult := s.tester.DepositERC20(token, investor01, big.NewInt(60000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"8"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor02, big.NewInt(28000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"4"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor03, big.NewInt(2000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"6"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor04, big.NewInt(5000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"4"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor05, big.NewInt(5500), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	time.Sleep(5 * time.Second)

	closeAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/close", "data":{"creator":"%s"}}`, creator))
	closeAuctionResult := s.tester.Advance(anyone, closeAuctionInput)
	s.Len(closeAuctionResult.Notices, 1)

	expectedCloseAuctionOutput := fmt.Sprintf(`auction closed - {"id":1,"token":"%s","creator":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"auction_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"auction_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"auction_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"auction_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"auction_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"auction_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)
	s.Equal(expectedCloseAuctionOutput, string(closeAuctionResult.Notices[0].Payload))
	s.Len(closeAuctionResult.DelegateCallVouchers, 1)
}

func (s *AuctionSystemSuite) TestSettleAuction() {
	anyone := common.HexToAddress("0x0000000000000000000000000000000000000001")
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
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

	createAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createAuctionResult := s.tester.DepositERC20(collateral, creator, big.NewInt(10000), createAuctionInput)
	s.Len(createAuctionResult.Notices, 1)

	expectedCreateAuctionOutput := fmt.Sprintf(`auction created - {"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateAuctionOutput, string(createAuctionResult.Notices[0].Payload))

	createOrderInput := []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"9"}}`)
	createOrderResult := s.tester.DepositERC20(token, investor01, big.NewInt(60000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"8"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor02, big.NewInt(28000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"4"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor03, big.NewInt(2000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"6"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor04, big.NewInt(5000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"4"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor05, big.NewInt(5500), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	time.Sleep(5 * time.Second)

	closeAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/close", "data":{"creator":"%s"}}`, creator))
	closeAuctionResult := s.tester.Advance(anyone, closeAuctionInput)
	s.Len(closeAuctionResult.Notices, 1)

	expectedCloseAuctionOutput := fmt.Sprintf(`auction closed - {"id":1,"token":"%s","creator":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"auction_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"auction_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"auction_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"auction_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"auction_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"auction_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)
	s.Equal(expectedCloseAuctionOutput, string(closeAuctionResult.Notices[0].Payload))
	s.Len(closeAuctionResult.DelegateCallVouchers, 1)

	time.Sleep(5 * time.Second)

	settleAuctionInput := []byte(`{"path":"auction/settle", "data":{"auction_id":1}}`)
	settleAuctionresult := s.tester.DepositERC20(token, creator, big.NewInt(108195), settleAuctionInput)
	s.Len(settleAuctionresult.Notices, 1)

	settledAt := baseTime + 10// baseTime

	expectedSettleAuctionOutput := fmt.Sprintf(`auction settled - {"id":1,"token":"%s","creator":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"settled","orders":[`+
		`{"id":1,"auction_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":2,"auction_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":3,"auction_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":4,"auction_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":5,"auction_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":6,"auction_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, settledAt, // Order 1
		investor02.Hex(), baseTime, settledAt, // Order 2
		investor03.Hex(), baseTime, settledAt, // Order 3
		investor04.Hex(), baseTime, settledAt, // Order 4
		investor05.Hex(), baseTime, settledAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, settledAt,
	)
	s.Equal(expectedSettleAuctionOutput, string(settleAuctionresult.Notices[0].Payload))
}

func (s *AuctionSystemSuite) TestExecuteAuctionCollateral() {
	anyone := common.HexToAddress("0x0000000000000000000000000000000000000001")
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
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

	createAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createAuctionResult := s.tester.DepositERC20(collateral, creator, big.NewInt(10000), createAuctionInput)
	s.Len(createAuctionResult.Notices, 1)

	expectedCreateAuctionOutput := fmt.Sprintf(`auction created - {"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateAuctionOutput, string(createAuctionResult.Notices[0].Payload))

	createOrderInput := []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"9"}}`)
	createOrderResult := s.tester.DepositERC20(token, investor01, big.NewInt(60000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"8"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor02, big.NewInt(28000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"4"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor03, big.NewInt(2000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"6"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor04, big.NewInt(5000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"4"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor05, big.NewInt(5500), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	time.Sleep(5 * time.Second)

	closeAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/close", "data":{"creator":"%s"}}`, creator))
	closeAuctionResult := s.tester.Advance(anyone, closeAuctionInput)
	s.Len(closeAuctionResult.Notices, 1)
	s.Len(closeAuctionResult.DelegateCallVouchers, 1)

	expectedCloseAuctionOutput := fmt.Sprintf(`auction closed - {"id":1,"token":"%s","creator":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"auction_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"auction_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"auction_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"auction_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"auction_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"auction_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)
	s.Equal(expectedCloseAuctionOutput, string(closeAuctionResult.Notices[0].Payload))
	s.Len(closeAuctionResult.DelegateCallVouchers, 1)

	findAuctionByIdInput := []byte(fmt.Sprintf(`{"path":"auction/id", "data":{"id":1}}`))

	findAuctionByIdResult := s.tester.Inspect(findAuctionByIdInput)
	s.Len(findAuctionByIdResult.Reports, 1)

	expectedFindAuctionByCreatorOutput := fmt.Sprintf(`[{"id":1,"token":"%s","creator":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"auction_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"auction_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"auction_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"auction_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"auction_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"auction_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}]`,
		token.Hex(),
		creator.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)

	findAuctionsByCreatorInput := []byte(fmt.Sprintf(`{"path":"auction/creator", "data":{"creator":"%s"}}`, creator))

	findAuctionsByCreatorResult := s.tester.Inspect(findAuctionsByCreatorInput)
	s.Len(findAuctionsByCreatorResult.Reports, 1)
	s.Equal(expectedFindAuctionByCreatorOutput, string(findAuctionsByCreatorResult.Reports[0].Payload))

	time.Sleep(6 * time.Second)

	executeAuctionCollateralInput := []byte(fmt.Sprintf(`{"path":"auction/execute-collateral", "data":{"auction_id":1}}`))
	executeAuctionCollateralResult := s.tester.DepositERC20(collateral, creator, big.NewInt(108195), executeAuctionCollateralInput)
	s.Len(executeAuctionCollateralResult.Notices, 1)

	collateralExecutedAt := baseTime + 11// baseTime

	expectedExecuteAuctionCollateralOutput := fmt.Sprintf(`auction collateral executed - {"auction_id":1,"token":"%s","creator":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"collateral_executed","orders":[`+
		`{"id":1,"auction_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":2,"auction_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":3,"auction_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":4,"auction_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":5,"auction_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"settled_by_collateral","created_at":%d,"updated_at":%d},`+
		`{"id":6,"auction_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, collateralExecutedAt, // Order 1
		investor02.Hex(), baseTime, collateralExecutedAt, // Order 2
		investor03.Hex(), baseTime, collateralExecutedAt, // Order 3
		investor04.Hex(), baseTime, collateralExecutedAt, // Order 4
		investor05.Hex(), baseTime, collateralExecutedAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, collateralExecutedAt,
	)
	s.Equal(expectedExecuteAuctionCollateralOutput, string(executeAuctionCollateralResult.Notices[0].Payload))
}

func (s *AuctionSystemSuite) TestFindAllAuctions() {
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	createAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createAuctionResult := s.tester.DepositERC20(collateral, creator, big.NewInt(10000), createAuctionInput)
	s.Len(createAuctionResult.Notices, 1)

	expectedCreateAuctionOutput := fmt.Sprintf(`auction created - {"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateAuctionOutput, string(createAuctionResult.Notices[0].Payload))

	findAllAuctionsInput := []byte(`{"path":"auction"}`)

	findAllAuctionsResult := s.tester.Inspect(findAllAuctionsInput)
	s.Len(findAllAuctionsResult.Reports, 1)

	expectedFindAllAuctionsOutput := fmt.Sprintf(`[{"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}]`, baseTime, closesAt, maturityAt)
	s.Equal(expectedFindAllAuctionsOutput, string(findAllAuctionsResult.Reports[0].Payload))
}

func (s *AuctionSystemSuite) TestFindAuctionById() {
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	createAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createAuctionResult := s.tester.DepositERC20(collateral, creator, big.NewInt(10000), createAuctionInput)
	s.Len(createAuctionResult.Notices, 1)

	expectedCreateAuctionOutput := fmt.Sprintf(`auction created - {"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateAuctionOutput, string(createAuctionResult.Notices[0].Payload))

	findAuctionByIdInput := []byte(fmt.Sprintf(`{"path":"auction/id", "data":{"id":1}}`))

	findAuctionByIdResult := s.tester.Inspect(findAuctionByIdInput)
	s.Len(findAuctionByIdResult.Reports, 1)

	expectedFindAuctionByIdOutput := fmt.Sprintf(`{"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedFindAuctionByIdOutput, string(findAuctionByIdResult.Reports[0].Payload))
}

func (s *AuctionSystemSuite) TestFindAuctionsByCreator() {
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
	collateral := common.HexToAddress("0x0000000000000000000000000000000000000008")
	token := common.HexToAddress("0x0000000000000000000000000000000000000009")

	baseTime := time.Now().Unix()
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	createAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createAuctionResult := s.tester.DepositERC20(collateral, creator, big.NewInt(10000), createAuctionInput)
	s.Len(createAuctionResult.Notices, 1)

	expectedCreateAuctionOutput := fmt.Sprintf(`auction created - {"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateAuctionOutput, string(createAuctionResult.Notices[0].Payload))

	findAuctionsByCreatorInput := []byte(fmt.Sprintf(`{"path":"auction/creator", "data":{"creator":"%s"}}`, creator))

	findAuctionsByCreatorResult := s.tester.Inspect(findAuctionsByCreatorInput)
	s.Len(findAuctionsByCreatorResult.Reports, 1)

	expectedFindAuctionsByCreatorOutput := fmt.Sprintf(`[{"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"0","total_raised":"0","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":0}]`, baseTime, closesAt, maturityAt)
	s.Equal(expectedFindAuctionsByCreatorOutput, string(findAuctionsByCreatorResult.Reports[0].Payload))
}

func (s *AuctionSystemSuite) TestFindAuctionsByInvestor() {
	anyone := common.HexToAddress("0x0000000000000000000000000000000000000001")
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
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

	createAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/create","data":{"token":"%s", "max_interest_rate":"10", "debt_issued":"100000", "closes_at":%d,"maturity_at":%d}}`, token, closesAt, maturityAt))
	createAuctionResult := s.tester.DepositERC20(collateral, creator, big.NewInt(10000), createAuctionInput)
	s.Len(createAuctionResult.Notices, 1)

	expectedCreateAuctionOutput := fmt.Sprintf(`auction created - {"id":1,"token":"0x0000000000000000000000000000000000000009","creator":"0x0000000000000000000000000000000000000007","collateral_address":"0x0000000000000000000000000000000000000008","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`, baseTime, closesAt, maturityAt)
	s.Equal(expectedCreateAuctionOutput, string(createAuctionResult.Notices[0].Payload))

	createOrderInput := []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"9"}}`)
	createOrderResult := s.tester.DepositERC20(token, investor01, big.NewInt(60000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"8"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor02, big.NewInt(28000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"4"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor03, big.NewInt(2000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"6"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor04, big.NewInt(5000), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	createOrderInput = []byte(`{"path": "order/create", "data": {"auction_id":1,"interest_rate":"4"}}`)
	createOrderResult = s.tester.DepositERC20(token, investor05, big.NewInt(5500), createOrderInput)
	s.Len(createOrderResult.Notices, 1)

	time.Sleep(5 * time.Second)

	closeAuctionInput := []byte(fmt.Sprintf(`{"path":"auction/close", "data":{"creator":"%s"}}`, creator))
	closeAuctionResult := s.tester.Advance(anyone, closeAuctionInput)
	s.Len(closeAuctionResult.Notices, 1)

	expectedCloseAuctionOutput := fmt.Sprintf(`auction closed - {"id":1,"token":"%s","creator":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"auction_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"auction_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"auction_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"auction_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"auction_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"auction_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)
	s.Equal(expectedCloseAuctionOutput, string(closeAuctionResult.Notices[0].Payload))
	s.Len(closeAuctionResult.DelegateCallVouchers, 1)

	expectedFindAuctionByCreatorOutput := fmt.Sprintf(`[{"id":1,"token":"%s","creator":"%s","collateral_address":"%s","collateral_amount":"10000","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"auction_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"auction_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"auction_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"auction_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"auction_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"auction_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}]`,
		token.Hex(),
		creator.Hex(),
		collateral.Hex(),
		investor01.Hex(), baseTime, closesAt, // Order 1
		investor02.Hex(), baseTime, closesAt, // Order 2
		investor03.Hex(), baseTime, closesAt, // Order 3
		investor04.Hex(), baseTime, closesAt, // Order 4
		investor05.Hex(), baseTime, closesAt, // Order 5
		investor01.Hex(), baseTime, closesAt, // Order 6 (rejected portion)
		baseTime, closesAt, maturityAt, closesAt,
	)

	findAuctionsByCreatorInput := []byte(fmt.Sprintf(`{"path":"auction/creator", "data":{"creator":"%s"}}`, creator))

	findAuctionsByCreatorResult := s.tester.Inspect(findAuctionsByCreatorInput)
	s.Len(findAuctionsByCreatorResult.Reports, 1)
	s.Equal(expectedFindAuctionByCreatorOutput, string(findAuctionsByCreatorResult.Reports[0].Payload))
}
