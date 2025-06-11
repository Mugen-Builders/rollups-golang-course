package auction

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
)

type FindAuctionsByCreatorInputDTO struct {
	Creator Address `json:"creator" validate:"required"`
}

type FindAuctionsByCreatorOutputDTO []*FindAuctionOutputDTO

type FindAuctionsByCreatorUseCase struct {
	AuctionRepository repository.AuctionRepository
}

func NewFindAuctionsByCreatorUseCase(AuctionRepository repository.AuctionRepository) *FindAuctionsByCreatorUseCase {
	return &FindAuctionsByCreatorUseCase{AuctionRepository: AuctionRepository}
}

func (f *FindAuctionsByCreatorUseCase) Execute(ctx context.Context, input *FindAuctionsByCreatorInputDTO) (*FindAuctionsByCreatorOutputDTO, error) {
	res, err := f.AuctionRepository.FindAuctionsByCreator(ctx, input.Creator)
	if err != nil {
		return nil, err
	}
	output := make(FindAuctionsByCreatorOutputDTO, len(res))
	for i, Auction := range res {
		orders := make([]*domain.Order, len(Auction.Orders))
		for j, order := range Auction.Orders {
			orders[j] = &domain.Order{
				Id:           order.Id,
				AuctionId:    order.AuctionId,
				Investor:     order.Investor,
				Amount:       order.Amount,
				InterestRate: order.InterestRate,
				State:        order.State,
				CreatedAt:    order.CreatedAt,
				UpdatedAt:    order.UpdatedAt,
			}
		}
		output[i] = &FindAuctionOutputDTO{
			Id:                Auction.Id,
			Token:             Auction.Token,
			Creator:           Auction.Creator,
			CollateralAddress: Auction.CollateralAddress,
			CollateralAmount:  Auction.CollateralAmount,
			DebtIssued:        Auction.DebtIssued,
			MaxInterestRate:   Auction.MaxInterestRate,
			TotalObligation:   Auction.TotalObligation,
			TotalRaised:       Auction.TotalRaised,
			State:             string(Auction.State),
			Orders:            orders,
			CreatedAt:         Auction.CreatedAt,
			ClosesAt:          Auction.ClosesAt,
			MaturityAt:        Auction.MaturityAt,
			UpdatedAt:         Auction.UpdatedAt,
		}
	}
	return &output, nil
}
