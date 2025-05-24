package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidVoting = errors.New("invalid voting")
)

type VotingStatus string

const (
	VotingStatusOpen     VotingStatus = "open"
	VotingStatusClosed   VotingStatus = "closed"
	VotingStatusCanceled VotingStatus = "canceled"
)

type Voting struct {
	ID          int             `gorm:"primaryKey;autoIncrement"`
	Title       string          `gorm:"not null"`
	Description string          `gorm:"not null"`
	StartDate   time.Time       `gorm:"not null;index"`
	EndDate     time.Time       `gorm:"not null;index"`
	Status      VotingStatus    `gorm:"not null;type:string;default:'open'"`
	Options     []*VotingOption `gorm:"foreignKey:VotingID"`
}

func NewVoting(title, description string, startDate, endDate time.Time) (*Voting, error) {
	voting := &Voting{
		Title:       title,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      VotingStatusOpen,
		Options:     make([]*VotingOption, 0),
	}
	if err := voting.validate(); err != nil {
		return nil, err
	}
	return voting, nil
}

func (v *Voting) validate() error {
	if v.Title == "" {
		return fmt.Errorf("%w: title cannot be empty", ErrInvalidVoting)
	}
	if v.StartDate.After(v.EndDate) {
		return fmt.Errorf("%w: start date must be before end date", ErrInvalidVoting)
	}
	if v.StartDate.Before(time.Now()) {
		return fmt.Errorf("%w: start date must be in the future", ErrInvalidVoting)
	}
	if v.Status != VotingStatusOpen && v.Status != VotingStatusClosed && v.Status != VotingStatusCanceled {
		return fmt.Errorf("%w: invalid status", ErrInvalidVoting)
	}
	return nil
}
