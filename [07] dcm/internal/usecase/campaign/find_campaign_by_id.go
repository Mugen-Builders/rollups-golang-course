package campaign

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
)

type FindCampaignByIdInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type FindCampaignByIdUseCase struct {
	CampaignRepository repository.CampaignRepository
}

func NewFindCampaignByIdUseCase(CampaignRepository repository.CampaignRepository) *FindCampaignByIdUseCase {
	return &FindCampaignByIdUseCase{CampaignRepository: CampaignRepository}
}

func (f *FindCampaignByIdUseCase) Execute(ctx context.Context, input *FindCampaignByIdInputDTO) (*FindCampaignOutputDTO, error) {
	res, err := f.CampaignRepository.FindCampaignById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	orders := make([]*entity.Order, len(res.Orders))
	for i, order := range res.Orders {
		orders[i] = &entity.Order{
			Id:           order.Id,
			CampaignId:   order.CampaignId,
			Investor:     order.Investor,
			Amount:       order.Amount,
			InterestRate: order.InterestRate,
			State:        order.State,
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
		}
	}
	return &FindCampaignOutputDTO{
		Id:                res.Id,
		Token:             res.Token,
		Debtor:            res.Debtor,
		CollateralAddress: res.CollateralAddress,
		CollateralAmount:  res.CollateralAmount,
		DebtIssued:        res.DebtIssued,
		MaxInterestRate:   res.MaxInterestRate,
		TotalObligation:   res.TotalObligation,
		TotalRaised:       res.TotalRaised,
		State:             string(res.State),
		Orders:            orders,
		CreatedAt:         res.CreatedAt,
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
