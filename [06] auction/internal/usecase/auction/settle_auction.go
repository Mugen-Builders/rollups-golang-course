package auction

import (
	"context"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type SettleAuctionInputDTO struct {
	AuctionId uint `json:"Auction_id"`
}

type SettleAuctionOutputDTO struct {
	Id                uint            `json:"id"`
	Token             Address         `json:"token"`
	Creator           Address         `json:"creator"`
	CollateralAddress Address         `json:"collateral_address"`
	DebtIssued        *uint256.Int    `json:"debt_issued"`
	MaxInterestRate   *uint256.Int    `json:"max_interest_rate"`
	TotalObligation   *uint256.Int    `json:"total_obligation"`
	Orders            []*domain.Order `json:"orders"`
	State             string          `json:"state"`
	ClosesAt          int64           `json:"closes_at"`
	MaturityAt        int64           `json:"maturity_at"`
	CreatedAt         int64           `json:"created_at"`
	UpdatedAt         int64           `json:"updated_at"`
}

type SettleAuctionUseCase struct {
	AuctionRepository repository.AuctionRepository
	OrderRepository   repository.OrderRepository
}

func NewSettleAuctionUseCase(
	AuctionRepository repository.AuctionRepository,
	orderRepository repository.OrderRepository,
) *SettleAuctionUseCase {
	return &SettleAuctionUseCase{
		AuctionRepository: AuctionRepository,
		OrderRepository:   orderRepository,
	}
}

func (uc *SettleAuctionUseCase) Execute(
	ctx context.Context,
	input *SettleAuctionInputDTO,
	deposit rollmelette.Deposit,
	metadata rollmelette.Metadata,
) (*SettleAuctionOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type: %T", deposit)
	}

	auction, err := uc.AuctionRepository.FindAuctionById(ctx, input.AuctionId)
	if err != nil {
		return nil, fmt.Errorf("error finding auction: %w", err)
	}

	if err := uc.Validate(auction, erc20Deposit, metadata); err != nil {
		return nil, err
	}

	var ordersToUpdate []*domain.Order
	for _, order := range auction.Orders {
		if order.State == domain.OrderStateAccepted || order.State == domain.OrderStatePartiallyAccepted {
			order.State = domain.OrderStateSettled
			order.UpdatedAt = metadata.BlockTimestamp
			ordersToUpdate = append(ordersToUpdate, order)
		}
	}
	for _, order := range ordersToUpdate {
		if _, err := uc.OrderRepository.UpdateOrder(ctx, order); err != nil {
			return nil, fmt.Errorf("error updating order: %w", err)
		}
	}

	auction.State = domain.AuctionStateSettled
	auction.UpdatedAt = metadata.BlockTimestamp
	res, err := uc.AuctionRepository.UpdateAuction(ctx, auction)
	if err != nil {
		return nil, fmt.Errorf("error updating auction: %w", err)
	}

	return &SettleAuctionOutputDTO{
		Id:                res.Id,
		Token:             res.Token,
		Creator:           res.Creator,
		CollateralAddress: res.CollateralAddress,
		DebtIssued:        res.DebtIssued,
		MaxInterestRate:   res.MaxInterestRate,
		TotalObligation:   res.TotalObligation,
		Orders:            res.Orders,
		State:             string(res.State),
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		CreatedAt:         res.CreatedAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}

func (uc *SettleAuctionUseCase) Validate(
	Auction *domain.Auction,
	deposit *rollmelette.ERC20Deposit,
	metadata rollmelette.Metadata,
) error {
	if metadata.BlockTimestamp > Auction.MaturityAt {
		return fmt.Errorf("the maturity date of the auction campaign has passed")
	}

	if Auction.State == domain.AuctionStateSettled {
		return fmt.Errorf("auction campaign already settled")
	}

	if Auction.State != domain.AuctionStateClosed {
		return fmt.Errorf("auction campaign not closed")
	}

	if deposit.Amount.Cmp(Auction.TotalObligation.ToBig()) < 0 {
		return fmt.Errorf("deposit amount is lower than the total obligation")
	}

	return nil
}
