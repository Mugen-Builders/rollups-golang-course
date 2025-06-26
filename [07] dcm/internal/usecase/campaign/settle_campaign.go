package campaign

import (
	"context"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type SettleCampaignInputDTO struct {
	CampaignId uint `json:"campaign_id" validate:"required"`
}

type SettleCampaignOutputDTO struct {
	Id                uint            `json:"id"`
	Token             Address         `json:"token"`
	Issuer            Address         `json:"issuer"`
	CollateralAddress Address         `json:"collateral_address"`
	CollateralAmount  *uint256.Int    `json:"collateral_amount"`
	DebtIssued        *uint256.Int    `json:"debt_issued"`
	MaxInterestRate   *uint256.Int    `json:"max_interest_rate"`
	TotalObligation   *uint256.Int    `json:"total_obligation"`
	TotalRaised       *uint256.Int    `json:"total_raised"`
	State             string          `json:"state"`
	Orders            []*entity.Order `json:"orders"`
	CreatedAt         int64           `json:"created_at"`
	ClosesAt          int64           `json:"closes_at"`
	MaturityAt        int64           `json:"maturity_at"`
	UpdatedAt         int64           `json:"updated_at"`
}

type SettleCampaignUseCase struct {
	CampaignRepository repository.CampaignRepository
	OrderRepository    repository.OrderRepository
}

func NewSettleCampaignUseCase(
	CampaignRepository repository.CampaignRepository,
	orderRepository repository.OrderRepository,
) *SettleCampaignUseCase {
	return &SettleCampaignUseCase{
		CampaignRepository: CampaignRepository,
		OrderRepository:    orderRepository,
	}
}

func (uc *SettleCampaignUseCase) Execute(
	ctx context.Context,
	input *SettleCampaignInputDTO,
	deposit rollmelette.Deposit,
	metadata rollmelette.Metadata,
) (*SettleCampaignOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type: %T", deposit)
	}

	campaign, err := uc.CampaignRepository.FindCampaignById(ctx, input.CampaignId)
	if err != nil {
		return nil, fmt.Errorf("error finding campaign: %w", err)
	}

	if err := uc.Validate(campaign, erc20Deposit, metadata); err != nil {
		return nil, err
	}

	var ordersToUpdate []*entity.Order
	for _, order := range campaign.Orders {
		if order.State == entity.OrderStateAccepted || order.State == entity.OrderStatePartiallyAccepted {
			order.State = entity.OrderStateSettled
			order.UpdatedAt = metadata.BlockTimestamp
			ordersToUpdate = append(ordersToUpdate, order)
		}
	}
	for _, order := range ordersToUpdate {
		if _, err := uc.OrderRepository.UpdateOrder(ctx, order); err != nil {
			return nil, fmt.Errorf("error updating order: %w", err)
		}
	}

	campaign.State = entity.CampaignStateSettled
	campaign.UpdatedAt = metadata.BlockTimestamp
	res, err := uc.CampaignRepository.UpdateCampaign(ctx, campaign)
	if err != nil {
		return nil, fmt.Errorf("error updating campaign: %w", err)
	}

	return &SettleCampaignOutputDTO{
		Id:                res.Id,
		Token:             res.Token,
		Issuer:            res.Issuer,
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

func (uc *SettleCampaignUseCase) Validate(
	Campaign *entity.Campaign,
	deposit *rollmelette.ERC20Deposit,
	metadata rollmelette.Metadata,
) error {
	if metadata.BlockTimestamp > Campaign.MaturityAt {
		return fmt.Errorf("the maturity date of the campaign campaign has passed")
	}

	if Campaign.State == entity.CampaignStateSettled {
		return fmt.Errorf("campaign campaign already settled")
	}

	if Campaign.State != entity.CampaignStateClosed {
		return fmt.Errorf("campaign campaign not closed")
	}

	if deposit.Value.Cmp(Campaign.TotalObligation.ToBig()) < 0 {
		return fmt.Errorf("deposit amount is lower than the total obligation")
	}

	if Campaign.Issuer != Address(deposit.Sender) {
		return fmt.Errorf("only the campaign issuer can settle the campaign")
	}
	return nil
}
