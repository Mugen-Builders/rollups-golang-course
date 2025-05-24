package sqlite

import (
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
)

func (r *SQLiteRepository) CreateVoting(voting *domain.Voting) error {
	return r.db.Create(voting).Error
}

func (r *SQLiteRepository) GetVotingByID(id int) (*domain.Voting, error) {
	var voting domain.Voting
	err := r.db.First(&voting, id).Error
	if err != nil {
		return nil, err
	}
	return &voting, nil
}

func (r *SQLiteRepository) ListVotings() ([]*domain.Voting, error) {
	var votings []*domain.Voting
	err := r.db.Find(&votings).Error
	if err != nil {
		return nil, err
	}
	return votings, nil
}

func (r *SQLiteRepository) UpdateVoting(voting *domain.Voting) error {
	return r.db.Save(voting).Error
}

func (r *SQLiteRepository) DeleteVoting(id int) error {
	return r.db.Delete(&domain.Voting{}, id).Error
}

func (r *SQLiteRepository) ListActiveVotings() ([]*domain.Voting, error) {
	var votings []*domain.Voting
	err := r.db.Where("status = ?", domain.VotingStatusOpen).Find(&votings).Error
	if err != nil {
		return nil, err
	}
	return votings, nil
}
