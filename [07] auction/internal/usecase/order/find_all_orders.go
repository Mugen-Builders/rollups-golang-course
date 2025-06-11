package order

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
)

type FindAllOrdersOutputDTO []*FindOrderOutputDTO

type FindAllOrdersUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewFindAllOrdersUseCase(orderRepository repository.OrderRepository) *FindAllOrdersUseCase {
	return &FindAllOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (f *FindAllOrdersUseCase) Execute(ctx context.Context) (*FindAllOrdersOutputDTO, error) {
	res, err := f.OrderRepository.FindAllOrders(ctx)
	if err != nil {
		return nil, err
	}
	output := make(FindAllOrdersOutputDTO, len(res))
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
