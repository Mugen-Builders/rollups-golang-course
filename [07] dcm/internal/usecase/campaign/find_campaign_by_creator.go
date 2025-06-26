package campaign

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
)

type FindCampaignsByIssuerInputDTO struct {
	Issuer Address `json:"issuer" validate:"required"`
}

type FindCampaignsByIssuerOutputDTO []*FindCampaignOutputDTO

type FindCampaignsByIssuerUseCase struct {
	CampaignRepository repository.CampaignRepository
}

func NewFindCampaignsByIssuerUseCase(CampaignRepository repository.CampaignRepository) *FindCampaignsByIssuerUseCase {
	return &FindCampaignsByIssuerUseCase{CampaignRepository: CampaignRepository}
}

func (f *FindCampaignsByIssuerUseCase) Execute(ctx context.Context, input *FindCampaignsByIssuerInputDTO) (*FindCampaignsByIssuerOutputDTO, error) {
	res, err := f.CampaignRepository.FindCampaignsByIssuer(ctx, input.Issuer)
	if err != nil {
		return nil, err
	}
	output := make(FindCampaignsByIssuerOutputDTO, len(res))
	for i, Campaign := range res {
		orders := make([]*entity.Order, len(Campaign.Orders))
		for j, order := range Campaign.Orders {
			orders[j] = &entity.Order{
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
		output[i] = &FindCampaignOutputDTO{
			Id:                Campaign.Id,
			Token:             Campaign.Token,
			Issuer:            Campaign.Issuer,
			CollateralAddress: Campaign.CollateralAddress,
			CollateralAmount:  Campaign.CollateralAmount,
			DebtIssued:        Campaign.DebtIssued,
			MaxInterestRate:   Campaign.MaxInterestRate,
			TotalObligation:   Campaign.TotalObligation,
			TotalRaised:       Campaign.TotalRaised,
			State:             string(Campaign.State),
			Orders:            orders,
			CreatedAt:         Campaign.CreatedAt,
			ClosesAt:          Campaign.ClosesAt,
			MaturityAt:        Campaign.MaturityAt,
			UpdatedAt:         Campaign.UpdatedAt,
		}
	}
	return &output, nil
}
