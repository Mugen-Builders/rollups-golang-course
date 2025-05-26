package voter

import (
	"context"
	"errors"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/pkg/custom_type"
)

type CreateVoterOutputDTO struct {
	Id      int     `json:"id"`
	Address Address `json:"address"`
}

type CreateVoterUseCase struct {
	VoterRepository repository.VoterRepository
}

func NewCreateVoterUseCase(voterRepository repository.VoterRepository) *CreateVoterUseCase {
	return &CreateVoterUseCase{VoterRepository: voterRepository}
}

func (uc *CreateVoterUseCase) Execute(ctx context.Context) (*CreateVoterOutputDTO, error) {
	msgSender, ok := ctx.Value("msg_sender").(string)
	if !ok {
		return nil, errors.New("error getting msg_sender")
	}
	voter, err := domain.NewVoter(HexToAddress(msgSender))
	if err != nil {
		return nil, err
	}
	err = uc.VoterRepository.CreateVoter(voter)
	if err != nil {
		return nil, err
	}
	return &CreateVoterOutputDTO{
		Id:      voter.ID,
		Address: voter.Address,
	}, nil
}
