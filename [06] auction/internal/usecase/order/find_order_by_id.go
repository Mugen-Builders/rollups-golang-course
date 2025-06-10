package order

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
)

type FindOrderByIdInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type FindOrderByIdUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewFindOrderByIdUseCase(orderRepository repository.OrderRepository) *FindOrderByIdUseCase {
	return &FindOrderByIdUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *FindOrderByIdUseCase) Execute(ctx context.Context, input *FindOrderByIdInputDTO) (*FindOrderOutputDTO, error) {
	res, err := c.OrderRepository.FindOrderById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	return &FindOrderOutputDTO{
		Id:           res.Id,
		AuctionId:    res.AuctionId,
		Investor:     res.Investor,
		Amount:       res.Amount,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}
