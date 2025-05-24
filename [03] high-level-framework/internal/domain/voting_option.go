package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidVotingOption = errors.New("invalid voting option")
)

type VotingOption struct {
	ID          int     `gorm:"primaryKey;autoIncrement"`
	VotingID    int     `gorm:"not null;index"`
	Description string  `gorm:"not null"`
	VoteCount   int     `gorm:"not null;default:0"`
	Voting      *Voting `gorm:"foreignKey:VotingID"`
}

func NewVotingOption(votingID int, description string) (*VotingOption, error) {
	option := &VotingOption{
		VotingID:    votingID,
		Description: description,
		VoteCount:   0,
	}
	if err := option.validate(); err != nil {
		return nil, err
	}
	return option, nil
}

func (v *VotingOption) validate() error {
	if v.Description == "" {
		return fmt.Errorf("%w: description cannot be empty", ErrInvalidVotingOption)
	}
	if v.VotingID <= 0 {
		return fmt.Errorf("%w: voting ID must be greater than zero", ErrInvalidVotingOption)
	}
	return nil
}
