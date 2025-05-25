package voting

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type FindAllActiveVotingsOutputDTO struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type FindAllActiveVotingsUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewFindAllActiveVotingsUseCase(votingRepository repository.VotingRepository) *FindAllActiveVotingsUseCase {
	return &FindAllActiveVotingsUseCase{VotingRepository: votingRepository}
}

func (uc *FindAllActiveVotingsUseCase) Execute(ctx context.Context) ([]*FindAllActiveVotingsOutputDTO, error) {
	votings, err := uc.VotingRepository.FindAllActiveVotings()
	if err != nil {
		return nil, err
	}
	var output []*FindAllActiveVotingsOutputDTO
	for _, v := range votings {
		output = append(output, &FindAllActiveVotingsOutputDTO{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
		})
	}
	return output, nil
}
