package voting_option

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type FindVotingOptionByIDInputDTO struct {
	Id int `json:"id"`
}

type FindVotingOptionByIDOutputDTO struct {
	Id          int    `json:"id"`
	VotingID    int    `json:"voting_id"`
	Description string `json:"description"`
}

type FindVotingOptionByIDUseCase struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewFindVotingOptionByIDUseCase(votingOptionRepository repository.VotingOptionRepository) *FindVotingOptionByIDUseCase {
	return &FindVotingOptionByIDUseCase{VotingOptionRepository: votingOptionRepository}
}

func (uc *FindVotingOptionByIDUseCase) Execute(ctx context.Context, input *FindVotingOptionByIDInputDTO) (*FindVotingOptionByIDOutputDTO, error) {
	option, err := uc.VotingOptionRepository.FindOptionByID(input.Id)
	if err != nil {
		return nil, err
	}
	return &FindVotingOptionByIDOutputDTO{
		Id:          option.ID,
		VotingID:    option.VotingID,
		Description: option.Description,
	}, nil
}
