package voting

import (
	"context"
	"time"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type CreateVotingInputDTO struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type CreateVotingOutputDTO struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type CreateVotingUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewCreateVotingUseCase(votingRepository repository.VotingRepository) *CreateVotingUseCase {
	return &CreateVotingUseCase{VotingRepository: votingRepository}
}

func (uc *CreateVotingUseCase) Execute(ctx context.Context, input *CreateVotingInputDTO) (*CreateVotingOutputDTO, error) {
	voting := &domain.Voting{
		Title:       input.Title,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	}
	err := uc.VotingRepository.CreateVoting(voting)
	if err != nil {
		return nil, err
	}
	return &CreateVotingOutputDTO{
		Id:          voting.ID,
		Title:       voting.Title,
		Description: voting.Description,
		StartDate:   voting.StartDate,
		EndDate:     voting.EndDate,
	}, nil
}
