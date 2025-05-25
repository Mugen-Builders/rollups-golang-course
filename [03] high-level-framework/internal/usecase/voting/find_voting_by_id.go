package voting

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type FindVotingByIDInputDTO struct {
	Id int `json:"id"`
}

type FindVotingByIDOutputDTO struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type FindVotingByIDUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewFindVotingByIDUseCase(votingRepository repository.VotingRepository) *FindVotingByIDUseCase {
	return &FindVotingByIDUseCase{VotingRepository: votingRepository}
}

func (uc *FindVotingByIDUseCase) Execute(ctx context.Context, input *FindVotingByIDInputDTO) (*FindVotingByIDOutputDTO, error) {
	voting, err := uc.VotingRepository.FindVotingByID(input.Id)
	if err != nil {
		return nil, err
	}
	return &FindVotingByIDOutputDTO{
		Id:          voting.ID,
		Title:       voting.Title,
		Description: voting.Description,
	}, nil
}
