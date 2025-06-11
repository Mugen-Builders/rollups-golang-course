package order

import (
	"context"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CreateOrderInputDTO struct {
	AuctionId    uint         `json:"auction_id" validate:"required"`
	InterestRate *uint256.Int `json:"interest_rate" validate:"required"`
}

type CreateOrderOutputDTO struct {
	Id           uint         `json:"id"`
	AuctionId    uint         `json:"auction_id"`
	Investor     Address      `json:"investor"`
	Amount       *uint256.Int `json:"amount"`
	InterestRate *uint256.Int `json:"interest_rate"`
	State        string       `json:"state"`
	CreatedAt    int64        `json:"created_at"`
}

type CreateOrderUseCase struct {
	OrderRepository   repository.OrderRepository
	AuctionRepository repository.AuctionRepository
}

func NewCreateOrderUseCase(orderRepository repository.OrderRepository, auctionRepository repository.AuctionRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository:   orderRepository,
		AuctionRepository: auctionRepository,
	}
}

func (c *CreateOrderUseCase) Execute(ctx context.Context, input *CreateOrderInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateOrderOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type provided for order creation: %T", deposit)
	}

	auction, err := c.AuctionRepository.FindAuctionById(ctx, input.AuctionId)
	if err != nil {
		return nil, fmt.Errorf("error finding auction campaigns: %w", err)
	}

	if auction.ClosesAt < metadata.BlockTimestamp {
		return nil, fmt.Errorf("auction campaign closed, order cannot be placed")
	}

	if Address(erc20Deposit.Token) != auction.Token {
		return nil, fmt.Errorf("invalid contract address provided for order creation: %v", erc20Deposit.Token)
	}

	if input.InterestRate.Gt(auction.MaxInterestRate) {
		return nil, fmt.Errorf("order interest rate exceeds active Auction max interest rate")
	}

	order, err := domain.NewOrder(
		auction.Id,
		Address(erc20Deposit.Sender),
		uint256.MustFromBig(erc20Deposit.Amount),
		input.InterestRate,
		metadata.BlockTimestamp,
	)
	if err != nil {
		return nil, err
	}

	res, err := c.OrderRepository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return &CreateOrderOutputDTO{
		Id:           res.Id,
		AuctionId:    res.AuctionId,
		Investor:     res.Investor,
		Amount:       res.Amount,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		CreatedAt:    res.CreatedAt,
	}, nil
}
