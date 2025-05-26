package root

import (
	"log/slog"
	"os"

	advance_handler "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/cartesi/handler/advance"
	inspect_handler "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/cartesi/handler/inspect"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository/factory"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/pkg/router"
)

func NewVotingSystem() *router.Router {
	repo, err := factory.NewRepositoryFromConnectionString("sqlite://voting.db")
	if err != nil {
		slog.Error("Failed to initialize repository", "error", err)
		os.Exit(1)
	}
	defer repo.Close()

	votingAdvanceHandlers := advance_handler.NewVotingAdvanceHandlers(repo)
	votingInspectHandlers := inspect_handler.NewVotingInspectHandlers(repo, repo)

	voterAdvanceHandlers := advance_handler.NewVoterAdvanceHandlers(repo)
	voterInspectHandlers := inspect_handler.NewVoterInspectHandlers(repo)

	votingOptionAdvanceHandlers := advance_handler.NewVotingOptionAdvanceHandlers(repo)
	votingOptionInspectHandlers := inspect_handler.NewVotingOptionInspectHandlers(repo)

	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)
	r.Use(router.ValidationMiddleware)
	r.Use(router.ErrorHandlingMiddleware)

	votingGroup := r.Group("voting")
	{
		votingGroup.HandleAdvance("create", votingAdvanceHandlers.CreateVoting)
		votingGroup.HandleAdvance("delete", votingAdvanceHandlers.DeleteVoting)

		votingGroup.HandleInspect("", votingInspectHandlers.FindAllVotings)
		votingGroup.HandleInspect("id", votingInspectHandlers.FindVotingByID)
		votingGroup.HandleInspect("active", votingInspectHandlers.FindAllActiveVotings)
		votingGroup.HandleInspect("results", votingInspectHandlers.GetVotingResults)
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
