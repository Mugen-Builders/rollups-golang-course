package voting

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type DeleteVotingInputDTO struct {
	Id int `json:"id"`
}

type DeleteVotingOutputDTO struct {
	Success bool `json:"success"`
}

type DeleteVotingUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewDeleteVotingUseCase(votingRepository repository.VotingRepository) *DeleteVotingUseCase {
	return &DeleteVotingUseCase{VotingRepository: votingRepository}
}

func (uc *DeleteVotingUseCase) Execute(ctx context.Context, input *DeleteVotingInputDTO) (*DeleteVotingOutputDTO, error) {
	err := uc.VotingRepository.DeleteVoting(input.Id)
	if err != nil {
		return &DeleteVotingOutputDTO{Success: false}, err
	}
	return &DeleteVotingOutputDTO{Success: true}, nil
}
