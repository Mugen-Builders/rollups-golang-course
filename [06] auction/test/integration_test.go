package test

import (
	"log/slog"
	"os"
	"testing"
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
}

func (s *AuctionSystemSuite) TestUpdateAuction() {
}

func (s *AuctionSystemSuite) TestSettleAuction() {
}

func (s *AuctionSystemSuite) TestExecuteAuctionCollateral() {
}

func (s *AuctionSystemSuite) TestCloseAuction() {
}

func (s *AuctionSystemSuite) TestFindAllAuctions() {
}

func (s *AuctionSystemSuite) TestFindAuctionById() {
}

func (s *AuctionSystemSuite) TestFindAuctionsByCreator() {
}

func (s *AuctionSystemSuite) TestFindAuctionsByInvestor() {
}