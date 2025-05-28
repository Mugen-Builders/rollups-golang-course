package auction

import (
	"context"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CreateAuctionInputDTO struct {
	Token           Address      `json:"token"`
	DebtIssued      *uint256.Int `json:"debt_issued"`
	MaxInterestRate *uint256.Int `json:"max_interest_rate"`
	ClosesAt        int64        `json:"closes_at"`
	MaturityAt      int64        `json:"maturity_at"`
}

type CreateAuctionOutputDTO struct {
	Id                uint            `json:"id"`
	Token             Address         `json:"token,omitempty"`
	Creator           Address         `json:"creator,omitempty"`
	CollateralAddress Address         `json:"collateral_address,omitempty"`
	CollateralAmount  *uint256.Int    `json:"collateral_amount,omitempty"`
	DebtIssued        *uint256.Int    `json:"debt_issued"`
	MaxInterestRate   *uint256.Int    `json:"max_interest_rate"`
	Orders            []*domain.Order `json:"orders"`
	State             string          `json:"state"`
	ClosesAt          int64           `json:"closes_at"`
	MaturityAt        int64           `json:"maturity_at"`
	CreatedAt         int64           `json:"created_at"`
}

type CreateAuctionUseCase struct {
	AuctionRepository repository.AuctionRepository
}

func NewCreateAuctionUseCase(
	AuctionRepository repository.AuctionRepository,
) *CreateAuctionUseCase {
	return &CreateAuctionUseCase{
		AuctionRepository: AuctionRepository,
	}
}

func (c *CreateAuctionUseCase) Execute(ctx context.Context, input *CreateAuctionInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateAuctionOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type: %T", deposit)
	}

	if input.DebtIssued.Cmp(uint256.NewInt(15000000)) > 0 {
		return nil, fmt.Errorf("%w: debt issued exceeds the maximum allowed value", domain.ErrInvalidAuction)
	}
	if input.ClosesAt > metadata.BlockTimestamp+15552000 {
		return nil, fmt.Errorf("%w: close date cannot be greater than 6 months", domain.ErrInvalidAuction)
	}
	if input.ClosesAt > input.MaturityAt {
		return nil, fmt.Errorf("%w: close date cannot be greater than maturity date", domain.ErrInvalidAuction)
	}
	if metadata.BlockTimestamp >= input.ClosesAt {
		return nil, fmt.Errorf("%w: creation date cannot be greater than or equal to close date", domain.ErrInvalidAuction)
	}

	Auctions, err := c.AuctionRepository.FindAuctionsByCreator(ctx, Address(erc20Deposit.Sender))
	if err != nil {
		return nil, fmt.Errorf("error retrieving Auctions: %w", err)
	}
	for _, Auction := range Auctions {
		if Auction.State != domain.AuctionStateSettled && metadata.BlockTimestamp-Auction.CreatedAt < 120*24*60*60 {
			return nil, fmt.Errorf("active Auction exists within the last 120 days")
		}
	}

	Auction, err := domain.NewAuction(
		input.Token,
		Address(erc20Deposit.Sender),
		Address(erc20Deposit.Token),
		uint256.MustFromBig(erc20Deposit.Amount),
		input.DebtIssued,
		input.MaxInterestRate,
		input.ClosesAt,
		input.MaturityAt,
		metadata.BlockTimestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating Auction: %w", err)
	}
	
	createdAuction, err := c.AuctionRepository.CreateAuction(ctx, Auction)
	if err != nil {
		return nil, fmt.Errorf("error creating Auction: %w", err)
	}

	return &CreateAuctionOutputDTO{
		Id:                createdAuction.Id,
		Token:             createdAuction.Token,
		Creator:           createdAuction.Creator,
		CollateralAddress: createdAuction.CollateralAddress,
		CollateralAmount:  createdAuction.CollateralAmount,
		DebtIssued:        createdAuction.DebtIssued,
		MaxInterestRate:   createdAuction.MaxInterestRate,
		Orders:            createdAuction.Orders,
		State:             string(createdAuction.State),
		ClosesAt:          createdAuction.ClosesAt,
		MaturityAt:        createdAuction.MaturityAt,
		CreatedAt:         createdAuction.CreatedAt,
	}, nil
}
