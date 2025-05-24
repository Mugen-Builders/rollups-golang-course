package repository

import "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"

type VotingRepository interface {
	CreateVoting(voting *domain.Voting) error
	GetVotingByID(id int) (*domain.Voting, error)
	ListVotings() ([]*domain.Voting, error)
	UpdateVoting(voting *domain.Voting) error
	DeleteVoting(id int) error
	ListActiveVotings() ([]*domain.Voting, error)
}

type VotingOptionRepository interface {
	CreateOption(option *domain.VotingOption) error
	GetOptionByID(id int) (*domain.VotingOption, error)
	GetOptionsByVotingID(votingID int) ([]*domain.VotingOption, error)
	UpdateOption(option *domain.VotingOption) error
	DeleteOption(id int) error
	IncrementVoteCount(id int) error
}

type VoterRepository interface {
	CreateVoter(voter *domain.Voter) error
	GetVoterByID(id int) (*domain.Voter, error)
	GetVoterByAddress(address string) (*domain.Voter, error)
	UpdateVoter(voter *domain.Voter) error
	DeleteVoter(id int) error
	HasVoted(voterID, votingID int) (bool, error)
}

type Repository interface {
	VotingRepository
	VotingOptionRepository
	VoterRepository
	Close() error
}
