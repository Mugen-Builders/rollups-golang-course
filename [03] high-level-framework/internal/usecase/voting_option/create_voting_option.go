package voting_option

import (
	"context"
	"errors"

	. "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/pkg/custom_type"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type CreateVotingOptionInputDTO struct {
	VotingID int `json:"voting_id"`
}

type CreateVotingOptionOutputDTO struct {
	Id       int `json:"id"`
	VotingID int `json:"voting_id"`
}

type CreateVotingOptionUseCase struct {
	VotingRepository       repository.VotingRepository
	VotingOptionRepository repository.VotingOptionRepository
}

func NewCreateVotingOptionUseCase(votingOptionRepository repository.VotingOptionRepository) *CreateVotingOptionUseCase {
	return &CreateVotingOptionUseCase{VotingOptionRepository: votingOptionRepository}
}

func (uc *CreateVotingOptionUseCase) Execute(ctx context.Context, input *CreateVotingOptionInputDTO) (*CreateVotingOptionOutputDTO, error) {
	voting, err := uc.VotingRepository.FindVotingByID(input.VotingID)
	if err != nil {
		return nil, err
	}
	if voting.Creator != HexToAddress(ctx.Value("msg_sender").(string)) {
		return nil, errors.New("unauthorized")
	}
	option, err := domain.NewVotingOption(input.VotingID)
	if err != nil {
		return nil, err
	}
	err = uc.VotingOptionRepository.CreateOption(option)
	if err != nil {
		return nil, err
	}
	return &CreateVotingOptionOutputDTO{
		Id:       option.ID,
		VotingID: option.VotingID,
	}, nil
}
