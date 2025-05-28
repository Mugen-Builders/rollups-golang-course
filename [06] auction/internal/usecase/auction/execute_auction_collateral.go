package auction

import (
	"context"
	"fmt"

	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type ExecuteAuctionCollateralInputDTO struct {
	AuctionId uint `json:"auction_id"`
}

type ExecuteAuctionCollateralOutputDTO struct {
	AuctionId         uint            `json:"auction_id"`
	CollateralAddress Address         `json:"collateral_address"`
	CollateralAmount  *uint256.Int    `json:"collateral_amount"`
	Orders            []*domain.Order `json:"orders"`
}

type ExecuteAuctionCollateralUseCase struct {
	AuctionRepository repository.AuctionRepository
}

func NewExecuteAuctionCollateralUseCase(auctionRepository repository.AuctionRepository) *ExecuteAuctionCollateralUseCase {
	return &ExecuteAuctionCollateralUseCase{AuctionRepository: auctionRepository}
}

func (uc *ExecuteAuctionCollateralUseCase) Execute(ctx context.Context, input *ExecuteAuctionCollateralInputDTO, metadata rollmelette.Metadata) (*ExecuteAuctionCollateralOutputDTO, error) {
	auction, err := uc.AuctionRepository.FindAuctionById(ctx, input.AuctionId)
	if err != nil {
		return nil, err
	}

	if metadata.BlockTimestamp < auction.MaturityAt {
		return nil, fmt.Errorf("the maturity date of the auction campaign has not passed")
	}

	if auction.State != domain.AuctionStateClosed {
		return nil, fmt.Errorf("auction campaign not closed")
	}

	auction.State = domain.AuctionStateCollateralExecuted

	res, err := uc.AuctionRepository.UpdateAuction(ctx, auction)
	if err != nil {
		return nil, err
	}

	return &ExecuteAuctionCollateralOutputDTO{
		AuctionId: res.Id,
		CollateralAddress: res.CollateralAddress,
		CollateralAmount:  res.CollateralAmount,
		Orders:    res.Orders,
	}, nil
}
