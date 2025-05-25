package voting_option

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type UpdateVotingOptionInputDTO struct {
	Id          int    `json:"id"`
	VotingID    int    `json:"voting_id"`
	Description string `json:"description"`
}

type UpdateVotingOptionOutputDTO struct {
	Id          int    `json:"id"`
	VotingID    int    `json:"voting_id"`
	Description string `json:"description"`
}

type UpdateVotingOptionUseCase struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewUpdateVotingOptionUseCase(votingOptionRepository repository.VotingOptionRepository) *UpdateVotingOptionUseCase {
	return &UpdateVotingOptionUseCase{VotingOptionRepository: votingOptionRepository}
}

func (uc *UpdateVotingOptionUseCase) Execute(ctx context.Context, input *UpdateVotingOptionInputDTO) (*UpdateVotingOptionOutputDTO, error) {
	option := &domain.VotingOption{
		ID:          input.Id,
		VotingID:    input.VotingID,
		Description: input.Description,
	}
	err := uc.VotingOptionRepository.UpdateOption(option)
	if err != nil {
		return nil, err
	}
	return &UpdateVotingOptionOutputDTO{
		Id:          option.ID,
		VotingID:    option.VotingID,
		Description: option.Description,
	}, nil
}
