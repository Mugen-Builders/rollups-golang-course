package campaign

import (
	"context"
	"fmt"
	"sort"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CloseCampaignInputDTO struct {
	Issuer Address `json:"issuer" validate:"required"`
}

type CloseCampaignOutputDTO struct {
	Id                uint            `json:"id"`
	Token             Address         `json:"token,omitempty"`
	Issuer            Address         `json:"issuer,omitempty"`
	CollateralAddress Address         `json:"collateral_address,omitempty"`
	CollateralAmount  *uint256.Int    `json:"collateral_amount,omitempty"`
	DebtIssued        *uint256.Int    `json:"debt_issued,omitempty"`
	MaxInterestRate   *uint256.Int    `json:"max_interest_rate,omitempty"`
	TotalObligation   *uint256.Int    `json:"total_obligation,omitempty"`
	TotalRaised       *uint256.Int    `json:"total_raised,omitempty"`
	State             string          `json:"state,omitempty"`
	Orders            []*entity.Order `json:"orders,omitempty"`
	CreatedAt         int64           `json:"created_at,omitempty"`
	ClosesAt          int64           `json:"closes_at,omitempty"`
	MaturityAt        int64           `json:"maturity_at,omitempty"`
	UpdatedAt         int64           `json:"updated_at,omitempty"`
}

type CloseCampaignUseCase struct {
	OrderRepository    repository.OrderRepository
	CampaignRepository repository.CampaignRepository
}

func NewCloseCampaignUseCase(CampaignRepository repository.CampaignRepository, orderRepository repository.OrderRepository) *CloseCampaignUseCase {
	return &CloseCampaignUseCase{
		OrderRepository:    orderRepository,
		CampaignRepository: CampaignRepository,
	}
}

func (u *CloseCampaignUseCase) Execute(ctx context.Context, input *CloseCampaignInputDTO, metadata rollmelette.Metadata) (*CloseCampaignOutputDTO, error) {
	// -------------------------------------------------------------------------
	// 1. Find ongoing campaign for the issuer
	// -------------------------------------------------------------------------
	campaigns, err := u.CampaignRepository.FindCampaignsByIssuer(ctx, input.Issuer)
	if err != nil {
		return nil, err
	}
	var ongoingCampaign *entity.Campaign
	for _, campaign := range campaigns {
		if campaign.State == entity.CampaignStateOngoing {
			ongoingCampaign = campaign
			break
		}
	}
	if ongoingCampaign == nil {
		return nil, fmt.Errorf("no ongoing campaign found, cannot close it")
	}

	// -------------------------------------------------------------------------
	// 2. Validate campaign expiration
	// -------------------------------------------------------------------------
	if metadata.BlockTimestamp < ongoingCampaign.ClosesAt {
		return nil, fmt.Errorf("campaign not expired yet, cannot close it")
	}

	// -------------------------------------------------------------------------
	// 3. Fetch and sort campaign orders
	// -------------------------------------------------------------------------
	orders, err := u.OrderRepository.FindOrdersByCampaignId(ctx, ongoingCampaign.Id)
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
	debtRemaining := new(uint256.Int).Set(ongoingCampaign.DebtIssued)
	totalCollected := uint256.NewInt(0)
	totalObligation := uint256.NewInt(0)

	for _, order := range orders {
		if debtRemaining.IsZero() {
			// Reject surplus orders
			order.State = entity.OrderStateRejected
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
			order.State = entity.OrderStateAccepted
			debtRemaining.Sub(debtRemaining, order.Amount)
		} else {
			order.State = entity.OrderStatePartiallyAccepted
			// Create rejected order for the surplus
			rejectedAmount := new(uint256.Int).Sub(order.Amount, acceptAmount)
			_, err := u.OrderRepository.CreateOrder(ctx, &entity.Order{
				CampaignId:   order.CampaignId,
				Investor:     order.Investor,
				Amount:       rejectedAmount,
				InterestRate: order.InterestRate,
				State:        entity.OrderStateRejected,
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
	twoThirds := new(uint256.Int).Mul(ongoingCampaign.DebtIssued, uint256.NewInt(2))
	twoThirds.Div(twoThirds, uint256.NewInt(3))
	if totalCollected.Lt(twoThirds) {
		// Cancel campaign and reject all orders
		for _, order := range orders {
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			if _, err := u.OrderRepository.UpdateOrder(ctx, order); err != nil {
				return nil, err
			}
		}
		ongoingCampaign.State = entity.CampaignStateCanceled
		ongoingCampaign.UpdatedAt = metadata.BlockTimestamp
		if _, err := u.CampaignRepository.UpdateCampaign(ctx, ongoingCampaign); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("campaign canceled due to insufficient funds collected, expected at least 2/3 of the debt issued: %s, got: %s", twoThirds.String(), totalCollected.String())
	}

	// -------------------------------------------------------------------------
	// 6. Close campaign and return result
	// -------------------------------------------------------------------------
	ongoingCampaign.State = entity.CampaignStateClosed
	ongoingCampaign.TotalObligation = totalObligation
	ongoingCampaign.TotalRaised = totalCollected
	ongoingCampaign.UpdatedAt = metadata.BlockTimestamp
	res, err := u.CampaignRepository.UpdateCampaign(ctx, ongoingCampaign)
	if err != nil {
		return nil, err
	}

	return &CloseCampaignOutputDTO{
		Id:                res.Id,
		Token:             res.Token,
		Issuer:            res.Issuer,
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
