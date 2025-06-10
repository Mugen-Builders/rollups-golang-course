package order

import (
	"context"
	"errors"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CancelOrderInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type CancelOrderOutputDTO struct {
	Id           uint
	AuctionId    uint
	Token        Address
	Investor     Address
	Amount       *uint256.Int
	InterestRate *uint256.Int
	State        string
	CreatedAt    int64
	UpdatedAt    int64
}

type CancelOrderUseCase struct {
	OrderRepository   repository.OrderRepository
	AuctionRepository repository.AuctionRepository
}

func NewCancelOrderUseCase(orderRepository repository.OrderRepository, auctionRepository repository.AuctionRepository) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		OrderRepository:   orderRepository,
		AuctionRepository: auctionRepository,
	}
}

func (c *CancelOrderUseCase) Execute(ctx context.Context, input *CancelOrderInputDTO, metadata rollmelette.Metadata) (*CancelOrderOutputDTO, error) {
	order, err := c.OrderRepository.FindOrderById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	if order.Investor != Address(metadata.MsgSender) {
		return nil, errors.New("only the investor can cancel the order")
	}
	auction, err := c.AuctionRepository.FindAuctionById(ctx, order.AuctionId)
	if err != nil {
		return nil, err
	}
	if auction.ClosesAt < metadata.BlockTimestamp {
		return nil, errors.New("cannot cancel order after Auction closes")
	}
	err = c.OrderRepository.DeleteOrder(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	return &CancelOrderOutputDTO{
		Id:           order.Id,
		AuctionId:    order.AuctionId,
		Token:        auction.Token,
		Investor:     order.Investor,
		Amount:       order.Amount,
		InterestRate: order.InterestRate,
		State:        string(order.State),
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}, nil
}
