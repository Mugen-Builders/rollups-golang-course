package campaign

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
)

type FindCampaignsByInvestorInputDTO struct {
	Investor Address `json:"investor" validate:"required"`
}

type FindCampaignsByInvestorOutputDTO []*FindCampaignOutputDTO

type FindCampaignsByInvestorUseCase struct {
	CampaignRepository repository.CampaignRepository
}

func NewFindCampaignsByInvestorUseCase(CampaignRepository repository.CampaignRepository) *FindCampaignsByInvestorUseCase {
	return &FindCampaignsByInvestorUseCase{CampaignRepository: CampaignRepository}
}

func (f *FindCampaignsByInvestorUseCase) Execute(ctx context.Context, input *FindCampaignsByInvestorInputDTO) (*FindCampaignsByInvestorOutputDTO, error) {
	res, err := f.CampaignRepository.FindCampaignsByInvestor(ctx, input.Investor)
	if err != nil {
		return nil, err
	}
	output := make(FindCampaignsByInvestorOutputDTO, len(res))
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
