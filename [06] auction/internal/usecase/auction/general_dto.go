package auction

import (
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
	"github.com/holiman/uint256"
)

type FindAuctionOutputDTO struct {
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
