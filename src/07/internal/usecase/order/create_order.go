package order

import (
	"context"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CreateOrderInputDTO struct {
	CampaignId   uint         `json:"campaign_id" validate:"required"`
	InterestRate *uint256.Int `json:"interest_rate" validate:"required"`
}

type CreateOrderOutputDTO struct {
	Id           uint         `json:"id"`
	CampaignId   uint         `json:"campaign_id"`
	Investor     Address      `json:"investor"`
	Amount       *uint256.Int `json:"amount"`
	InterestRate *uint256.Int `json:"interest_rate"`
	State        string       `json:"state"`
	CreatedAt    int64        `json:"created_at"`
}

type CreateOrderUseCase struct {
	OrderRepository    repository.OrderRepository
	CampaignRepository repository.CampaignRepository
}

func NewCreateOrderUseCase(orderRepository repository.OrderRepository, campaignRepository repository.CampaignRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository:    orderRepository,
		CampaignRepository: campaignRepository,
	}
}

func (c *CreateOrderUseCase) Execute(ctx context.Context, input *CreateOrderInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateOrderOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type provided for order creation: %T", deposit)
	}

	campaign, err := c.CampaignRepository.FindCampaignById(ctx, input.CampaignId)
	if err != nil {
		return nil, fmt.Errorf("error finding campaign campaigns: %w", err)
	}

	if campaign.ClosesAt < metadata.BlockTimestamp {
		return nil, fmt.Errorf("campaign campaign closed, order cannot be placed")
	}

	if Address(erc20Deposit.Token) != campaign.Token {
		return nil, fmt.Errorf("invalid contract address provided for order creation: %v", erc20Deposit.Token)
	}

	if input.InterestRate.Gt(campaign.MaxInterestRate) {
		return nil, fmt.Errorf("order interest rate exceeds active Campaign max interest rate")
	}

	order, err := entity.NewOrder(
		campaign.Id,
		Address(erc20Deposit.Sender),
		uint256.MustFromBig(erc20Deposit.Value),
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
		CampaignId:   res.CampaignId,
		Investor:     res.Investor,
		Amount:       res.Amount,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		CreatedAt:    res.CreatedAt,
	}, nil
}
