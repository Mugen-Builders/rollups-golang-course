package voting

import (
	"context"
	"errors"
	"time"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/pkg/custom_type"
)

type CreateVotingInputDTO struct {
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type CreateVotingOutputDTO struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Creator   Address   `json:"creator"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type CreateVotingUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewCreateVotingUseCase(votingRepository repository.VotingRepository) *CreateVotingUseCase {
	return &CreateVotingUseCase{VotingRepository: votingRepository}
}

func (uc *CreateVotingUseCase) Execute(ctx context.Context, input *CreateVotingInputDTO) (*CreateVotingOutputDTO, error) {
	msgSender, ok := ctx.Value("msg_sender").(string)
	if !ok {
		return nil, errors.New("error getting msg_sender")
	}
	voting, err := domain.NewVoting(input.Title, HexToAddress(msgSender), input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}
	err = uc.VotingRepository.CreateVoting(voting)
	if err != nil {
		return nil, err
	}
	return &CreateVotingOutputDTO{
		Id:        voting.ID,
		Title:     voting.Title,
		StartDate: voting.StartDate,
		EndDate:   voting.EndDate,
	}, nil
}
