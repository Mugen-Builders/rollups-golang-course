package order

import (
	"context"
	"errors"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CancelOrderInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type CancelOrderOutputDTO struct {
	Id           uint
	CampaignId   uint
	Token        Address
	Investor     Address
	Amount       *uint256.Int
	InterestRate *uint256.Int
	State        string
	CreatedAt    int64
	UpdatedAt    int64
}

type CancelOrderUseCase struct {
	OrderRepository    repository.OrderRepository
	CampaignRepository repository.CampaignRepository
}

func NewCancelOrderUseCase(orderRepository repository.OrderRepository, campaignRepository repository.CampaignRepository) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		OrderRepository:    orderRepository,
		CampaignRepository: campaignRepository,
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
	campaign, err := c.CampaignRepository.FindCampaignById(ctx, order.CampaignId)
	if err != nil {
		return nil, err
	}
	if campaign.State == entity.CampaignStateClosed {
		return nil, errors.New("cannot cancel order after Campaign closes")
	}
	err = c.OrderRepository.DeleteOrder(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	return &CancelOrderOutputDTO{
		Id:           order.Id,
		CampaignId:   order.CampaignId,
		Token:        campaign.Token,
		Investor:     order.Investor,
		Amount:       order.Amount,
		InterestRate: order.InterestRate,
		State:        string(order.State),
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}, nil
}
