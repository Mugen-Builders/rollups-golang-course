package user

import (
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
	"github.com/holiman/uint256"
)

type BalanceOfInputDTO struct {
	Token   Address `json:"token"`
	Address Address `json:"address" validate:"required"`
}

type FindUserOutputDTO struct {
	Id              uint         `json:"id"`
	Role            string       `json:"role"`
	Address         Address      `json:"address"`
	InvestmentLimit *uint256.Int `json:"investment_limit,omitempty" gorm:"type:bigint"`
	CreatedAt       int64        `json:"created_at"`
	UpdatedAt       int64        `json:"updated_at"`
}
