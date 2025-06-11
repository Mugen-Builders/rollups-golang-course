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

type ExecuteAuctionCollateralInputDTO struct {
	AuctionId uint `json:"auction_id" validate:"required"`
}

type ExecuteAuctionCollateralOutputDTO struct {
	AuctionId         uint            `json:"auction_id"`
	Token             Address         `json:"token"`
	Creator           Address         `json:"creator"`
	CollateralAddress Address         `json:"collateral_address"`
	CollateralAmount  *uint256.Int    `json:"collateral_amount"`
	DebtIssued        *uint256.Int    `json:"debt_issued"`
	MaxInterestRate   *uint256.Int    `json:"max_interest_rate"`
	TotalObligation   *uint256.Int    `json:"total_obligation"`
	TotalRaised       *uint256.Int    `json:"total_raised"`
	State             string          `json:"state"`
	Orders            []*domain.Order `json:"orders"`
	CreatedAt         int64           `json:"created_at"`
	ClosesAt          int64           `json:"closes_at"`
	MaturityAt        int64           `json:"maturity_at"`
	UpdatedAt         int64           `json:"updated_at"`
}

type ExecuteAuctionCollateralUseCase struct {
	AuctionRepository repository.AuctionRepository
	OrderRepository   repository.OrderRepository
}

func NewExecuteAuctionCollateralUseCase(auctionRepository repository.AuctionRepository, orderRepository repository.OrderRepository) *ExecuteAuctionCollateralUseCase {
	return &ExecuteAuctionCollateralUseCase{
		AuctionRepository: auctionRepository,
		OrderRepository:   orderRepository,
	}
}

func (uc *ExecuteAuctionCollateralUseCase) Execute(ctx context.Context, input *ExecuteAuctionCollateralInputDTO, metadata rollmelette.Metadata) (*ExecuteAuctionCollateralOutputDTO, error) {
	auction, err := uc.AuctionRepository.FindAuctionById(ctx, input.AuctionId)
	if err != nil {
		return nil, err
	}

	if err := uc.Validate(auction, metadata); err != nil {
		return nil, err
	}

	var ordersToUpdate []*domain.Order
	for _, order := range auction.Orders {
		if order.State == domain.OrderStateAccepted || order.State == domain.OrderStatePartiallyAccepted {
			order.State = domain.OrderStateSettledByCollateral
			order.UpdatedAt = metadata.BlockTimestamp
			ordersToUpdate = append(ordersToUpdate, order)
		}
	}
	for _, order := range ordersToUpdate {
		if _, err := uc.OrderRepository.UpdateOrder(ctx, order); err != nil {
			return nil, fmt.Errorf("error updating order: %w", err)
		}
	}

	auction.State = domain.AuctionStateCollateralExecuted

	res, err := uc.AuctionRepository.UpdateAuction(ctx, auction)
	if err != nil {
		return nil, err
	}

	return &ExecuteAuctionCollateralOutputDTO{
		AuctionId:         res.Id,
		Token:             res.Token,
		Creator:           res.Creator,
		CollateralAddress: res.CollateralAddress,
		CollateralAmount:  res.CollateralAmount,
		DebtIssued:        res.DebtIssued,
		MaxInterestRate:   res.MaxInterestRate,
		TotalObligation:   res.TotalObligation,
		TotalRaised:       res.TotalRaised,
		State:             string(res.State),
		Orders:            res.Orders,
		CreatedAt:         res.CreatedAt,
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}

func (uc *ExecuteAuctionCollateralUseCase) Validate(auction *domain.Auction, metadata rollmelette.Metadata) error {
	if metadata.BlockTimestamp < auction.MaturityAt {
		return fmt.Errorf("the maturity date of the auction campaign has not passed")
	}
	if auction.State != domain.AuctionStateClosed {
		return fmt.Errorf("auction campaign not closed")
	}
	return nil
}
