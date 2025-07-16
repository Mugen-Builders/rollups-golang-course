package order

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
)

type FindOrdersByCampaignIdInputDTO struct {
	CampaignId uint `json:"campaign_id" validate:"required"`
}

type FindOrdersByCampaignIdOutputDTO []*FindOrderOutputDTO

type FindOrdersByCampaignIdUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewFindOrdersByCampaignIdUseCase(orderRepository repository.OrderRepository) *FindOrdersByCampaignIdUseCase {
	return &FindOrdersByCampaignIdUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *FindOrdersByCampaignIdUseCase) Execute(ctx context.Context, input *FindOrdersByCampaignIdInputDTO) (*FindOrdersByCampaignIdOutputDTO, error) {
	res, err := c.OrderRepository.FindOrdersByCampaignId(ctx, input.CampaignId)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByCampaignIdOutputDTO, len(res))
	for i, order := range res {
		output[i] = &FindOrderOutputDTO{
			Id:           order.Id,
			CampaignId:   order.CampaignId,
			Investor:     order.Investor,
			Amount:       order.Amount,
			InterestRate: order.InterestRate,
			State:        string(order.State),
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
		}
	}
	return &output, nil
}
