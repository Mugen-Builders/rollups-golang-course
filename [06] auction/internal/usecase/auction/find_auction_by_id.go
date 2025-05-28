package auction

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
)

type FindAuctionByIdInputDTO struct {
	Id uint `json:"id"`
}

type FindAuctionByIdUseCase struct {
	AuctionRepository repository.AuctionRepository
}

func NewFindAuctionByIdUseCase(AuctionRepository repository.AuctionRepository) *FindAuctionByIdUseCase {
	return &FindAuctionByIdUseCase{AuctionRepository: AuctionRepository}
}

func (f *FindAuctionByIdUseCase) Execute(ctx context.Context, input *FindAuctionByIdInputDTO) (*FindAuctionOutputDTO, error) {
	res, err := f.AuctionRepository.FindAuctionById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	var orders []*domain.Order
	for _, order := range res.Orders {
		orders = append(orders, &domain.Order{
			Id:           order.Id,
			AuctionId:    order.AuctionId,
			Investor:     order.Investor,
			Amount:       order.Amount,
			InterestRate: order.InterestRate,
			State:        order.State,
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
		})
	}
	return &FindAuctionOutputDTO{
		Id:                res.Id,
		Token:             res.Token,
		Creator:           res.Creator,
		CollateralAddress: res.CollateralAddress,
		DebtIssued:        res.DebtIssued,
		MaxInterestRate:   res.MaxInterestRate,
		TotalObligation:   res.TotalObligation,
		Orders:            orders,
		State:             string(res.State),
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		CreatedAt:         res.CreatedAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
