package auction

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
)

type FindAuctionsByInvestorInputDTO struct {
	Investor Address `json:"investor"`
}

type FindAuctionsByInvestorOutputDTO []*FindAuctionOutputDTO

type FindAuctionsByInvestorUseCase struct {
	AuctionRepository repository.AuctionRepository
}

func NewFindAuctionsByInvestorUseCase(AuctionRepository repository.AuctionRepository) *FindAuctionsByInvestorUseCase {
	return &FindAuctionsByInvestorUseCase{AuctionRepository: AuctionRepository}
}

func (f *FindAuctionsByInvestorUseCase) Execute(ctx context.Context, input *FindAuctionsByInvestorInputDTO) (*FindAuctionsByInvestorOutputDTO, error) {
	res, err := f.AuctionRepository.FindAuctionsByInvestor(ctx, input.Investor)
	if err != nil {
		return nil, err
	}
	output := make(FindAuctionsByInvestorOutputDTO, len(res))
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
			DebtIssued:        Auction.DebtIssued,
			MaxInterestRate:   Auction.MaxInterestRate,
			TotalObligation:   Auction.TotalObligation,
			Orders:            orders,
			State:             string(Auction.State),
			ClosesAt:          Auction.ClosesAt,
			MaturityAt:        Auction.MaturityAt,
			CreatedAt:         Auction.CreatedAt,
			UpdatedAt:         Auction.UpdatedAt,
		}
	}
	return &output, nil
}
