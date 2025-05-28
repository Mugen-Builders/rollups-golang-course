package order

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
)

type FindOrdersByStateInputDTO struct {
	AuctionId uint   `json:"Auction_id"`
	State     string `json:"state"`
}

type FindOrdersByStateOutputDTO []*FindOrderOutputDTO

type FindOrdersByStateUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewFindOrdersByStateUseCase(orderRepository repository.OrderRepository) *FindOrdersByStateUseCase {
	return &FindOrdersByStateUseCase{
		OrderRepository: orderRepository,
	}
}

func (f *FindOrdersByStateUseCase) Execute(ctx context.Context, input *FindOrdersByStateInputDTO) (FindOrdersByStateOutputDTO, error) {
	res, err := f.OrderRepository.FindOrdersByState(ctx, input.AuctionId, input.State)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByStateOutputDTO, len(res))
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
	return output, nil
}
