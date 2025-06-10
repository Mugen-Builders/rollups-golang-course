package domain

import (
	"errors"
	"fmt"

	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
	"github.com/holiman/uint256"
)

var (
	ErrAuctionNotFound = errors.New("auction not found")
	ErrInvalidAuction  = errors.New("invalid Auction")
)

type AuctionState string

const (
	AuctionStateClosed             AuctionState = "closed"
	AuctionStateOngoing            AuctionState = "ongoing"
	AuctionStateCanceled           AuctionState = "canceled"
	AuctionStateSettled            AuctionState = "settled"
	AuctionStateCollateralExecuted AuctionState = "collateral_executed"
)

type Auction struct {
	Id                uint         `json:"id" gorm:"primaryKey"`
	Token             Address      `json:"token,omitempty" gorm:"custom_type:text;not null"`
	Creator           Address      `json:"creator,omitempty" gorm:"custom_type:text;not null"`
	CollateralAddress Address      `json:"collateral_address,omitempty" gorm:"custom_type:text;not null"`
	CollateralAmount  *uint256.Int `json:"collateral_amount,omitempty" gorm:"custom_type:text;not null"`
	DebtIssued        *uint256.Int `json:"debt_issued,omitempty" gorm:"custom_type:text;not null"`
	MaxInterestRate   *uint256.Int `json:"max_interest_rate,omitempty" gorm:"custom_type:text;not null"`
	TotalObligation   *uint256.Int `json:"total_obligation,omitempty" gorm:"custom_type:text;not null;default:0"`
	State             AuctionState `json:"state,omitempty" gorm:"custom_type:text;not null"`
	Orders            []*Order     `json:"orders,omitempty" gorm:"foreignKey:AuctionId;constraint:OnDelete:CASCADE"`
	ClosesAt          int64        `json:"closes_at,omitempty" gorm:"not null"`
	MaturityAt        int64        `json:"maturity_at,omitempty" gorm:"not null"`
	CreatedAt         int64        `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt         int64        `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewAuction(token Address, creator Address, collateral_address Address, collateral_amount *uint256.Int, debt_issued *uint256.Int, maxInterestRate *uint256.Int, closesAt int64, maturityAt int64, createdAt int64) (*Auction, error) {
	Auction := &Auction{
		Token:             token,
		Creator:           creator,
		CollateralAddress: collateral_address,
		CollateralAmount:  collateral_amount,
		DebtIssued:        debt_issued,
		MaxInterestRate:   maxInterestRate,
		State:             AuctionStateOngoing,
		Orders:            []*Order{},
		ClosesAt:          closesAt,
		MaturityAt:        maturityAt,
		CreatedAt:         createdAt,
	}
	if err := Auction.validate(); err != nil {
		return nil, err
	}
	return Auction, nil
}

func (a *Auction) validate() error {
	if a.Token == (Address{}) {
		return fmt.Errorf("%w: invalid token address", ErrInvalidAuction)
	}
	if a.Creator == (Address{}) {
		return fmt.Errorf("%w: invalid creator address", ErrInvalidAuction)
	}
	if a.CollateralAddress == (Address{}) {
		return fmt.Errorf("%w: invalid collateral_address address", ErrInvalidAuction)
	}
	if a.CollateralAmount.Sign() == 0 {
		return fmt.Errorf("%w: collateral amount cannot be zero", ErrInvalidAuction)
	}
	if a.DebtIssued.Sign() == 0 {
		return fmt.Errorf("%w: debt issued cannot be zero", ErrInvalidAuction)
	}
	if a.MaxInterestRate.Sign() == 0 {
		return fmt.Errorf("%w: max interest rate cannot be zero", ErrInvalidAuction)
	}
	if a.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidAuction)
	}
	if a.ClosesAt == 0 {
		return fmt.Errorf("%w: close date is missing", ErrInvalidAuction)
	}
	if a.MaturityAt == 0 {
		return fmt.Errorf("%w: maturity date is missing", ErrInvalidAuction)
	}
	return nil
}
