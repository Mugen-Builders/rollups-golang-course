package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidVoter = errors.New("invalid voter")
)

type Voter struct {
	ID      int    `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"not null"`
	Address string `gorm:"not null;uniqueIndex"`
}

func NewVoter(name, address string) (*Voter, error) {
	voter := &Voter{
		Name:    name,
		Address: address,
	}
	if err := voter.validate(); err != nil {
		return nil, err
	}
	return voter, nil
}

func (v *Voter) validate() error {
	if v.Name == "" {
		return fmt.Errorf("%w: name cannot be empty", ErrInvalidVoter)
	}
	if v.Address == "" {
		return fmt.Errorf("%w: address cannot be empty", ErrInvalidVoter)
	}
	return nil
}
