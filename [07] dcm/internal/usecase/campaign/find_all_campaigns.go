package campaign

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
)

type FindAllCampaignsOutputDTO []*FindCampaignOutputDTO

type FindAllCampaignsUseCase struct {
	CampaignRepository repository.CampaignRepository
}

func NewFindAllCampaignsUseCase(CampaignRepository repository.CampaignRepository) *FindAllCampaignsUseCase {
	return &FindAllCampaignsUseCase{CampaignRepository: CampaignRepository}
}

func (f *FindAllCampaignsUseCase) Execute(ctx context.Context) (*FindAllCampaignsOutputDTO, error) {
	res, err := f.CampaignRepository.FindAllCampaigns(ctx)
	if err != nil {
		return nil, err
	}
	output := make(FindAllCampaignsOutputDTO, len(res))
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
			Debtor:            Campaign.Debtor,
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
