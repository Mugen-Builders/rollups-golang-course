package campaign

import (
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
	"github.com/holiman/uint256"
)

type FindCampaignOutputDTO struct {
	Id                uint            `json:"id"`
	Token             Address         `json:"token"`
	Debtor            Address         `json:"debtor"`
	CollateralAddress Address         `json:"collateral_address"`
	CollateralAmount  *uint256.Int    `json:"collateral_amount"`
	DebtIssued        *uint256.Int    `json:"debt_issued"`
	MaxInterestRate   *uint256.Int    `json:"max_interest_rate"`
	TotalObligation   *uint256.Int    `json:"total_obligation"`
	TotalRaised       *uint256.Int    `json:"total_raised"`
	State             string          `json:"state"`
	Orders            []*entity.Order `json:"orders"`
	CreatedAt         int64           `json:"created_at"`
	ClosesAt          int64           `json:"closes_at"`
	MaturityAt        int64           `json:"maturity_at"`
	UpdatedAt         int64           `json:"updated_at"`
}
