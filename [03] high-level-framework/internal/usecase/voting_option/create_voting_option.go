package voting_option

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type CreateVotingOptionInputDTO struct {
	VotingID int `json:"voting_id"`
}

type CreateVotingOptionOutputDTO struct {
	Id       int `json:"id"`
	VotingID int `json:"voting_id"`
}

type CreateVotingOptionUseCase struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewCreateVotingOptionUseCase(votingOptionRepository repository.VotingOptionRepository) *CreateVotingOptionUseCase {
	return &CreateVotingOptionUseCase{VotingOptionRepository: votingOptionRepository}
}

func (uc *CreateVotingOptionUseCase) Execute(ctx context.Context, input *CreateVotingOptionInputDTO) (*CreateVotingOptionOutputDTO, error) {
	option := &domain.VotingOption{
		VotingID: input.VotingID,
	}
	err := uc.VotingOptionRepository.CreateOption(option)
	if err != nil {
		return nil, err
	}
	return &CreateVotingOptionOutputDTO{
		Id:       option.ID,
		VotingID: option.VotingID,
	}, nil
}
