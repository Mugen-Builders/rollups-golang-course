package root

import (
	"context"
	"log/slog"
	"os"

	advance_handler "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/cartesi/handler/advance"
	inspect_handler "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/cartesi/handler/inspect"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository/factory"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/pkg/router"
	"github.com/rollmelette/rollmelette"
	"github.com/spf13/cobra"
)

const (
	CMD_NAME = "rollup"
)

var (
	useMemoryDB bool
	Cmd         = &cobra.Command{
		Use:   "voting-" + CMD_NAME,
		Short: "Runs Voting Rollup",
		Long:  `Cartesi Rollup Application for voting`,
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
			false: "sqlite://voting.db",
		}[useMemoryDB],
	)
	if err != nil {
		slog.Error("Failed to setup database", "error", err, "type", map[bool]string{true: "in-memory", false: "persistent"}[useMemoryDB])
		os.Exit(1)
	}
	slog.Info("Database initialized", "type", map[bool]string{true: "in-memory", false: "persistent"}[useMemoryDB])
	defer repo.Close()

	r := NewVotingSystem(repo)
	opts := rollmelette.NewRunOpts()
	if err := rollmelette.Run(ctx, opts, r); err != nil {
		slog.Error("Failed to run rollmelette", "error", err)
		os.Exit(1)
	}
}

func NewVotingSystem(repo repository.Repository) *router.Router {
	votingAdvanceHandlers := advance_handler.NewVotingAdvanceHandlers(repo)
	votingInspectHandlers := inspect_handler.NewVotingInspectHandlers(repo, repo)

	voterAdvanceHandlers := advance_handler.NewVoterAdvanceHandlers(repo)
	voterInspectHandlers := inspect_handler.NewVoterInspectHandlers(repo)

	votingOptionAdvanceHandlers := advance_handler.NewVotingOptionAdvanceHandlers(repo, repo)
	votingOptionInspectHandlers := inspect_handler.NewVotingOptionInspectHandlers(repo)

	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)
	r.Use(router.ValidationMiddleware)
	r.Use(router.ErrorHandlingMiddleware)

	votingGroup := r.Group("voting")
	{
		votingGroup.HandleAdvance("create", votingAdvanceHandlers.CreateVoting)
		votingGroup.HandleAdvance("delete", votingAdvanceHandlers.DeleteVoting)
		votingGroup.HandleAdvance("vote", votingAdvanceHandlers.Vote)
		votingGroup.HandleAdvance("update-status", votingAdvanceHandlers.UpdateStatus)

		votingGroup.HandleInspect("", votingInspectHandlers.FindAllVotings)
		votingGroup.HandleInspect("id", votingInspectHandlers.FindVotingByID)
		votingGroup.HandleInspect("active", votingInspectHandlers.FindAllActiveVotings)
		votingGroup.HandleInspect("results", votingInspectHandlers.GetResults)
	}

	voterGroup := r.Group("voter")
	{
		voterGroup.HandleAdvance("create", voterAdvanceHandlers.CreateVoter)
		voterGroup.HandleAdvance("delete", voterAdvanceHandlers.DeleteVoter)

		voterGroup.HandleInspect("id", voterInspectHandlers.FindVoterByID)
		voterGroup.HandleInspect("address", voterInspectHandlers.FindVoterByAddress)
	}

	votingOptionGroup := r.Group("voting-option")
	{
		votingOptionGroup.HandleAdvance("create", votingOptionAdvanceHandlers.CreateVotingOption)
		votingOptionGroup.HandleAdvance("delete", votingOptionAdvanceHandlers.DeleteVotingOption)

		votingOptionGroup.HandleInspect("id", votingOptionInspectHandlers.FindVotingOptionByID)
		votingOptionGroup.HandleInspect("voting", votingOptionInspectHandlers.FindAllOptionsByVotingID)
	}
	return r
}
