package root

import (
	"context"
	"log/slog"
	"os"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/cartesi/handler/advance"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/cartesi/handler/inspect"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository/factory"
	"github.com/henriquemarlon/cartesi-golang-series/auction/pkg/router"
	"github.com/rollmelette/rollmelette"
	"github.com/spf13/cobra"
)

const (
	CMD_NAME = "rollup"
)

var (
	useMemoryDB bool
	Cmd         = &cobra.Command{
		Use:   "auction-" + CMD_NAME,
		Short: "Runs Auction Rollup",
		Long:  `Cartesi Rollup Application for auction`,
		Run:   run,
	}
)

func init() {
	Cmd.PersistentFlags().BoolVar(
		&useMemoryDB,
		"memory-db",
		false,
		"Use in-memory SQLite database instead of persistent",
	)
}

func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	repo, err := factory.NewRepositoryFromConnectionString(
		map[bool]string{
			true:  "sqlite://:memory:",
			false: "sqlite://auction.db",
		}[useMemoryDB],
	)
	if err != nil {
		slog.Error("Failed to setup database", "error", err, "type", map[bool]string{true: "in-memory", false: "persistent"}[useMemoryDB])
		os.Exit(1)
	}
	slog.Info("Database initialized", "type", map[bool]string{true: "in-memory", false: "persistent"}[useMemoryDB])
	defer repo.Close()

	r := NewAuctionSystem(repo)
	opts := rollmelette.NewRunOpts()
	if err := rollmelette.Run(ctx, opts, r); err != nil {
		slog.Error("Failed to run rollmelette", "error", err)
		os.Exit(1)
	}
}

func NewAuctionSystem(repo repository.Repository) *router.Router {
	auctionAdvanceHandlers := advance.NewAuctionAdvanceHandlers(repo, repo)
	auctionInspectHandlers := inspect.NewAuctionInspectHandlers(repo)

	orderAdvanceHandlers := advance.NewOrderAdvanceHandlers(repo, repo)
	orderInspectHandlers := inspect.NewOrderInspectHandlers(repo)

	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)
	r.Use(router.ValidationMiddleware)
	r.Use(router.ErrorHandlingMiddleware)

	orderGroup := r.Group("order")
	{
		orderGroup.HandleAdvance("create", orderAdvanceHandlers.CreateOrder)
		orderGroup.HandleAdvance("cancel", orderAdvanceHandlers.CancelOrder)
		orderGroup.HandleInspect("", orderInspectHandlers.FindAllOrders)
		orderGroup.HandleInspect("id", orderInspectHandlers.FindOrderById)
		orderGroup.HandleInspect("investor", orderInspectHandlers.FindOrdersByInvestor)
		orderGroup.HandleInspect("auction", orderInspectHandlers.FindBidsByAuctionId)
	}

	auctionGroup := r.Group("auction")
	{
		auctionGroup.HandleAdvance("create", auctionAdvanceHandlers.CreateAuction)
		auctionGroup.HandleAdvance("settle", auctionAdvanceHandlers.SettleAuction)
		auctionGroup.HandleAdvance("execute-collateral", auctionAdvanceHandlers.ExecuteAuctionCollateral)
		auctionGroup.HandleAdvance("close", auctionAdvanceHandlers.CloseAuction)
		auctionGroup.HandleInspect("", auctionInspectHandlers.FindAllAuctions)
		auctionGroup.HandleInspect("id", auctionInspectHandlers.FindAuctionById)
		auctionGroup.HandleInspect("creator", auctionInspectHandlers.FindAuctionsByCreator)
		auctionGroup.HandleInspect("investor", auctionInspectHandlers.FindAuctionsByInvestor)
	}
	return r
}
