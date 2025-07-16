package domain

import (
	"errors"
	"fmt"
)

var (
	ErrOptionNotFound      = errors.New("voting option not found")
	ErrInvalidVotingOption = errors.New("invalid voting option")
	ErrInvalidOption       = errors.New("voting option does not belong to the voting")
)

type VotingOption struct {
	ID        int     `gorm:"primaryKey;autoIncrement"`
	VotingID  int     `gorm:"not null;index"`
	VoterID   int     `gorm:"not null;index"`
	VoteCount int     `gorm:"not null;default:0"`
	Voting    *Voting `gorm:"foreignKey:VotingID"`
}

func NewVotingOption(votingID int) (*VotingOption, error) {
	option := &VotingOption{
		VotingID:  votingID,
		VoteCount: 0,
	}
	if err := option.validate(); err != nil {
		return nil, err
	}
	return option, nil
}

func (v *VotingOption) validate() error {
	if v.VotingID <= 0 {
		return fmt.Errorf("%w: voting ID must be greater than zero", ErrInvalidVotingOption)
	}
	return nil
}
