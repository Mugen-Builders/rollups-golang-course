package order

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
)

type FindOrdersByAuctionIdInputDTO struct {
	AuctionId uint `json:"auction_id" validate:"required"`
}

type FindOrdersByAuctionIdOutputDTO []*FindOrderOutputDTO

type FindOrdersByAuctionIdUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewFindOrdersByAuctionIdUseCase(orderRepository repository.OrderRepository) *FindOrdersByAuctionIdUseCase {
	return &FindOrdersByAuctionIdUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *FindOrdersByAuctionIdUseCase) Execute(ctx context.Context, input *FindOrdersByAuctionIdInputDTO) (*FindOrdersByAuctionIdOutputDTO, error) {
	res, err := c.OrderRepository.FindOrdersByAuctionId(ctx, input.AuctionId)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByAuctionIdOutputDTO, len(res))
	for i, order := range res {
		output[i] = &FindOrderOutputDTO{
			Id:           order.Id,
			AuctionId:    order.AuctionId,
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
