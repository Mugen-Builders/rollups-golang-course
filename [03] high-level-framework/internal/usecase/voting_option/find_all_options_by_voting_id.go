package voting_option

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type FindAllOptionsByVotingIDInputDTO struct {
	VotingID int `json:"voting_id"`
}

type FindAllOptionsByVotingIDOutputDTO struct {
	Id       int `json:"id"`
	VotingID int `json:"voting_id"`
}

type FindAllOptionsByVotingIDUseCase struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewFindAllOptionsByVotingIDUseCase(votingOptionRepository repository.VotingOptionRepository) *FindAllOptionsByVotingIDUseCase {
	return &FindAllOptionsByVotingIDUseCase{VotingOptionRepository: votingOptionRepository}
}

func (uc *FindAllOptionsByVotingIDUseCase) Execute(ctx context.Context, input *FindAllOptionsByVotingIDInputDTO) ([]*FindAllOptionsByVotingIDOutputDTO, error) {
	options, err := uc.VotingOptionRepository.FindAllOptionsByVotingID(input.VotingID)
	if err != nil {
		return nil, err
	}
	var output []*FindAllOptionsByVotingIDOutputDTO
	for _, o := range options {
		output = append(output, &FindAllOptionsByVotingIDOutputDTO{
			Id:       o.ID,
			VotingID: o.VotingID,
		})
	}
	return output, nil
}
