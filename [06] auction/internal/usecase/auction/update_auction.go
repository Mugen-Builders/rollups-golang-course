package auction

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type UpdateAuctionInputDTO struct {
	Id              uint         `json:"id"`
	DebtIssued      *uint256.Int `json:"debt_issued"`
	MaxInterestRate *uint256.Int `json:"max_interest_rate"`
	TotalObligation *uint256.Int `json:"total_obligation"`
	State           string       `json:"state"`
	ClosesAt        int64        `json:"closes_at"`
	MaturityAt      int64        `json:"maturity_at"`
}

type UpdateAuctionOutputDTO struct {
	Id                uint            `json:"id"`
	Token             Address         `json:"token"`
	Creator           Address         `json:"creator"`
	CollateralAddress Address         `json:"collateral_address"`
	DebtIssued        *uint256.Int    `json:"debt_issued"`
	MaxInterestRate   *uint256.Int    `json:"max_interest_rate"`
	TotalObligation   *uint256.Int    `json:"total_obligation"`
	Orders            []*domain.Order `json:"orders"`
	State             string          `json:"state"`
	ClosesAt          int64           `json:"closes_at"`
	MaturityAt        int64           `json:"maturity_at"`
	CreatedAt         int64           `json:"created_at"`
	UpdatedAt         int64           `json:"updated_at"`
}

type UpdateAuctionUsecase struct {
	AuctionRepository repository.AuctionRepository
}

func NewUpdateAuctionUseCase(AuctionRepository repository.AuctionRepository) *UpdateAuctionUsecase {
	return &UpdateAuctionUsecase{
		AuctionRepository: AuctionRepository,
	}
}

func (uc *UpdateAuctionUsecase) Execute(ctx context.Context, input UpdateAuctionInputDTO, metadata rollmelette.Metadata) (*UpdateAuctionOutputDTO, error) {
	Auction, err := uc.AuctionRepository.UpdateAuction(ctx, &domain.Auction{
		Id:              input.Id,
		DebtIssued:      input.DebtIssued,
		MaxInterestRate: input.MaxInterestRate,
		TotalObligation: input.TotalObligation,
		State:           domain.AuctionState(input.State),
		ClosesAt:        input.ClosesAt,
		MaturityAt:      input.ClosesAt,
		UpdatedAt:       metadata.BlockTimestamp,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateAuctionOutputDTO{
		Id:                Auction.Id,
		Token:             Auction.Token,
		Creator:           Auction.Creator,
		CollateralAddress: Auction.CollateralAddress,
		DebtIssued:        Auction.DebtIssued,
		MaxInterestRate:   Auction.MaxInterestRate,
		TotalObligation:   Auction.TotalObligation,
		Orders:            Auction.Orders,
		State:             string(Auction.State),
		ClosesAt:          Auction.ClosesAt,
		MaturityAt:        Auction.MaturityAt,
		CreatedAt:         Auction.CreatedAt,
		UpdatedAt:         Auction.UpdatedAt,
	}, nil
}
