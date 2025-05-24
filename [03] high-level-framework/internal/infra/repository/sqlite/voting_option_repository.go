package sqlite

import (
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateOption(option *domain.VotingOption) error {
	return r.db.Create(option).Error
}

func (r *SQLiteRepository) GetOptionByID(id int) (*domain.VotingOption, error) {
	var option domain.VotingOption
	err := r.db.First(&option, id).Error
	if err != nil {
		return nil, err
	}
	return &option, nil
}

func (r *SQLiteRepository) GetOptionsByVotingID(votingID int) ([]*domain.VotingOption, error) {
	var options []*domain.VotingOption
	err := r.db.Where("voting_id = ?", votingID).Find(&options).Error
	if err != nil {
		return nil, err
	}
	return options, nil
}

func (r *SQLiteRepository) UpdateOption(option *domain.VotingOption) error {
	return r.db.Save(option).Error
}

func (r *SQLiteRepository) DeleteOption(id int) error {
	return r.db.Delete(&domain.VotingOption{}, id).Error
}

func (r *SQLiteRepository) IncrementVoteCount(id int) error {
	return r.db.Model(&domain.VotingOption{}).Where("id = ?", id).UpdateColumn("vote_count", gorm.Expr("vote_count + ?", 1)).Error
}
