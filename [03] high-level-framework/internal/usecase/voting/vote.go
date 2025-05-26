package voting

import (
	"errors"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/pkg/custom_type"
	"github.com/rollmelette/rollmelette"
)

var (
	ErrVotingClosed   = errors.New("voting is closed")
	ErrVoterNotFound  = errors.New("voter not found")
	ErrOptionNotFound = errors.New("voting option not found")
	ErrAlreadyVoted   = errors.New("voter has already voted in this voting")
	ErrInvalidOption  = errors.New("voting option does not belong to the voting")
)

type VoteInputDTO struct {
	VotingID int `json:"voting_id"`
	OptionID int `json:"option_id"`
}

type VoteOutputDTO struct {
	VotingID  int     `json:"voting_id"`
	OptionID  int     `json:"option_id"`
	Voter     Address `json:"voter"`
	VoteCount int     `json:"vote_count"`
}

type VoteUseCase struct {
	VotingRepository       repository.VotingRepository
	VotingOptionRepository repository.VotingOptionRepository
	VoterRepository        repository.VoterRepository
}

func NewVoteUseCase(
	votingRepository repository.VotingRepository,
	votingOptionRepository repository.VotingOptionRepository,
	voterRepository repository.VoterRepository,
) *VoteUseCase {
	return &VoteUseCase{
		VotingRepository:       votingRepository,
		VotingOptionRepository: votingOptionRepository,
		VoterRepository:        voterRepository,
	}
}

func (u *VoteUseCase) Execute(input VoteInputDTO, metadata *rollmelette.Metadata) (*VoteOutputDTO, error) {
	voting, err := u.VotingRepository.FindVotingByID(input.VotingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find voting: %w", err)
	}

	if voting.Status != domain.VotingStatusOpen {
		return nil, ErrVotingClosed
	}

	voter, err := u.VoterRepository.FindVoterByAddress(Address(metadata.MsgSender))
	if err != nil {
		return nil, ErrVoterNotFound
	}

	hasVoted, err := u.VoterRepository.HasVoted(voter.ID, voting.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if voter has voted: %w", err)
	}
	if hasVoted {
		return nil, ErrAlreadyVoted
	}

	option, err := u.VotingOptionRepository.FindOptionByID(input.OptionID)
	if err != nil {
		return nil, ErrOptionNotFound
	}

	if option.VotingID != voting.ID {
		return nil, ErrInvalidOption
	}

	err = u.VotingOptionRepository.IncrementVoteCount(option.ID, voter.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to increment vote count: %w", err)
	}

	return &VoteOutputDTO{
		VotingID:  voting.ID,
		OptionID:  option.ID,
		Voter:     voter.Address,
		VoteCount: option.VoteCount + 1,
	}, nil
}
