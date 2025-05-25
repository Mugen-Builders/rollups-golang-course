package voting

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type UpdateVotingInputDTO struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateVotingOutputDTO struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateVotingUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewUpdateVotingUseCase(votingRepository repository.VotingRepository) *UpdateVotingUseCase {
	return &UpdateVotingUseCase{VotingRepository: votingRepository}
}

func (uc *UpdateVotingUseCase) Execute(ctx context.Context, input *UpdateVotingInputDTO) (*UpdateVotingOutputDTO, error) {
	voting := &domain.Voting{
		ID:          input.Id,
		Title:       input.Title,
		Description: input.Description,
	}
	err := uc.VotingRepository.UpdateVoting(voting)
	if err != nil {
		return nil, err
	}
	return &UpdateVotingOutputDTO{
		Id:          voting.ID,
		Title:       voting.Title,
		Description: voting.Description,
	}, nil
}
