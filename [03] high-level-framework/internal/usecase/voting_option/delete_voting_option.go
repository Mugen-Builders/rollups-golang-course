package voting_option

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type DeleteVotingOptionInputDTO struct {
	Id int `json:"id"`
}

type DeleteVotingOptionOutputDTO struct {
	Success bool `json:"success"`
}

type DeleteVotingOptionUseCase struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewDeleteVotingOptionUseCase(votingOptionRepository repository.VotingOptionRepository) *DeleteVotingOptionUseCase {
	return &DeleteVotingOptionUseCase{VotingOptionRepository: votingOptionRepository}
}

func (uc *DeleteVotingOptionUseCase) Execute(ctx context.Context, input *DeleteVotingOptionInputDTO) (*DeleteVotingOptionOutputDTO, error) {
	err := uc.VotingOptionRepository.DeleteOption(input.Id)
	if err != nil {
		return &DeleteVotingOptionOutputDTO{Success: false}, err
	}
	return &DeleteVotingOptionOutputDTO{Success: true}, nil
}
