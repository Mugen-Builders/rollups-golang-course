//go:build wireinject
// +build wireinject

package root

import (
	"github.com/google/wire"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/cartesi/handler/advance"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/cartesi/handler/inspect"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
)

func NewHandlers(repo repository.Repository) (*Handlers, error) {
	wire.Build(
		// Bind repository interfaces
		wire.Bind(new(repository.UserRepository), new(repository.Repository)),
		wire.Bind(new(repository.OrderRepository), new(repository.Repository)),
		wire.Bind(new(repository.CampaignRepository), new(repository.Repository)),
		// Advance handlers
		advance.NewOrderAdvanceHandlers,
		advance.NewUserAdvanceHandlers,
		advance.NewCampaignAdvanceHandlers,
		// Inspect handlers
		inspect.NewOrderInspectHandlers,
		inspect.NewUserInspectHandlers,
		inspect.NewCampaignInspectHandlers,
		wire.Struct(new(Handlers), "*"),
	)
	return &Handlers{}, nil
}

// Handlers contains all handler dependencies
type Handlers struct {
	// Advance handlers
	OrderAdvanceHandlers    *advance.OrderAdvanceHandlers
	UserAdvanceHandlers     *advance.UserAdvanceHandlers
	CampaignAdvanceHandlers *advance.CampaignAdvanceHandlers

	// Inspect handlers
	OrderInspectHandlers    *inspect.OrderInspectHandlers
	UserInspectHandlers     *inspect.UserInspectHandlers
	CampaignInspectHandlers *inspect.CampaignInspectHandlers
}
