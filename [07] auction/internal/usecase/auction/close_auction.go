package auction

import (
	"context"
	"fmt"
	"sort"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CloseAuctionInputDTO struct {
	Creator Address `json:"creator" validate:"required"`
}

type CloseAuctionOutputDTO struct {
	Id                uint            `json:"id"`
	Token             Address         `json:"token,omitempty"`
	Creator           Address         `json:"creator,omitempty"`
	CollateralAddress Address         `json:"collateral_address,omitempty"`
	CollateralAmount  *uint256.Int    `json:"collateral_amount,omitempty"`
	DebtIssued        *uint256.Int    `json:"debt_issued,omitempty"`
	MaxInterestRate   *uint256.Int    `json:"max_interest_rate,omitempty"`
	TotalObligation   *uint256.Int    `json:"total_obligation,omitempty"`
	TotalRaised       *uint256.Int    `json:"total_raised,omitempty"`
	State             string          `json:"state,omitempty"`
	Orders            []*domain.Order `json:"orders,omitempty"`
	CreatedAt         int64           `json:"created_at,omitempty"`
	ClosesAt          int64           `json:"closes_at,omitempty"`
	MaturityAt        int64           `json:"maturity_at,omitempty"`
	UpdatedAt         int64           `json:"updated_at,omitempty"`
}

type CloseAuctionUseCase struct {
	OrderRepository   repository.OrderRepository
	AuctionRepository repository.AuctionRepository
}

func NewCloseAuctionUseCase(AuctionRepository repository.AuctionRepository, orderRepository repository.OrderRepository) *CloseAuctionUseCase {
	return &CloseAuctionUseCase{
		OrderRepository:   orderRepository,
		AuctionRepository: AuctionRepository,
	}
}

func (u *CloseAuctionUseCase) Execute(ctx context.Context, input *CloseAuctionInputDTO, metadata rollmelette.Metadata) (*CloseAuctionOutputDTO, error) {
	// -------------------------------------------------------------------------
	// 1. Find ongoing auction for the creator
	// -------------------------------------------------------------------------
	auctions, err := u.AuctionRepository.FindAuctionsByCreator(ctx, input.Creator)
	if err != nil {
		return nil, err
	}
	var ongoingAuction *domain.Auction
	for _, auction := range auctions {
		if auction.State == domain.AuctionStateOngoing {
			ongoingAuction = auction
			break
		}
	}
	if ongoingAuction == nil {
		return nil, fmt.Errorf("no ongoing auction found, cannot close it")
	}

	// -------------------------------------------------------------------------
	// 2. Validate auction expiration
	// -------------------------------------------------------------------------
	if metadata.BlockTimestamp < ongoingAuction.ClosesAt {
		return nil, fmt.Errorf("auction not expired yet, cannot close it")
	}

	// -------------------------------------------------------------------------
	// 3. Fetch and sort auction orders
	// -------------------------------------------------------------------------
	orders, err := u.OrderRepository.FindOrdersByAuctionId(ctx, ongoingAuction.Id)
	if err != nil {
		return nil, err
	}
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].InterestRate.Cmp(orders[j].InterestRate) == 0 {
			return orders[i].Amount.Cmp(orders[j].Amount) > 0
		}
		return orders[i].InterestRate.Cmp(orders[j].InterestRate) < 0
	})

	// -------------------------------------------------------------------------
	// 4. Select winning orders and calculate obligations
	// -------------------------------------------------------------------------
	debtRemaining := new(uint256.Int).Set(ongoingAuction.DebtIssued)
	totalCollected := uint256.NewInt(0)
	totalObligation := uint256.NewInt(0)

	for _, order := range orders {
		if debtRemaining.IsZero() {
			// Reject surplus orders
			order.State = domain.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			if _, err := u.OrderRepository.UpdateOrder(ctx, order); err != nil {
				return nil, err
			}
			continue
		}

		// Accept full or partial order
		acceptAmount := order.Amount
		if debtRemaining.Lt(order.Amount) {
			acceptAmount = new(uint256.Int).Set(debtRemaining)
		}
		interest := new(uint256.Int).Mul(acceptAmount, order.InterestRate)
		interest.Div(interest, uint256.NewInt(100))

		orderObligation := new(uint256.Int).Add(acceptAmount, interest)
		totalCollected.Add(totalCollected, acceptAmount)
		totalObligation.Add(totalObligation, orderObligation)

		if debtRemaining.Cmp(order.Amount) >= 0 {
			order.State = domain.OrderStateAccepted
			debtRemaining.Sub(debtRemaining, order.Amount)
		} else {
			order.State = domain.OrderStatePartiallyAccepted
			// Create rejected order for the surplus
			rejectedAmount := new(uint256.Int).Sub(order.Amount, acceptAmount)
			_, err := u.OrderRepository.CreateOrder(ctx, &domain.Order{
				AuctionId:    order.AuctionId,
				Investor:     order.Investor,
				Amount:       rejectedAmount,
				InterestRate: order.InterestRate,
				State:        domain.OrderStateRejected,
				CreatedAt:    order.CreatedAt,
				UpdatedAt:    metadata.BlockTimestamp,
			})
			if err != nil {
				return nil, err
			}
			debtRemaining.Clear()
		}
		order.Amount = acceptAmount
		order.UpdatedAt = metadata.BlockTimestamp
		if _, err := u.OrderRepository.UpdateOrder(ctx, order); err != nil {
			return nil, err
		}
	}

	// -------------------------------------------------------------------------
	// 5. Check if minimum funding (2/3) was reached
	// -------------------------------------------------------------------------
	twoThirds := new(uint256.Int).Mul(ongoingAuction.DebtIssued, uint256.NewInt(2))
	twoThirds.Div(twoThirds, uint256.NewInt(3))
	if totalCollected.Lt(twoThirds) {
		// Cancel auction and reject all orders
		for _, order := range orders {
			order.State = domain.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			if _, err := u.OrderRepository.UpdateOrder(ctx, order); err != nil {
				return nil, err
			}
		}
		ongoingAuction.State = domain.AuctionStateCanceled
		ongoingAuction.UpdatedAt = metadata.BlockTimestamp
		if _, err := u.AuctionRepository.UpdateAuction(ctx, ongoingAuction); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("auction canceled due to insufficient funds collected, expected at least 2/3 of the debt issued: %s, got: %s", twoThirds.String(), totalCollected.String())
	}

	// -------------------------------------------------------------------------
	// 6. Close auction and return result
	// -------------------------------------------------------------------------
	ongoingAuction.State = domain.AuctionStateClosed
	ongoingAuction.TotalObligation = totalObligation
	ongoingAuction.TotalRaised = totalCollected
	ongoingAuction.UpdatedAt = metadata.BlockTimestamp
	res, err := u.AuctionRepository.UpdateAuction(ctx, ongoingAuction)
	if err != nil {
		return nil, err
	}

	return &CloseAuctionOutputDTO{
		Id:                res.Id,
		Token:             res.Token,
		Creator:           res.Creator,
		CollateralAddress: res.CollateralAddress,
		CollateralAmount:  res.CollateralAmount,
		DebtIssued:        res.DebtIssued,
		MaxInterestRate:   res.MaxInterestRate,
		TotalObligation:   res.TotalObligation,
		TotalRaised:       res.TotalRaised,
		Orders:            res.Orders,
		State:             string(res.State),
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		CreatedAt:         res.CreatedAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
