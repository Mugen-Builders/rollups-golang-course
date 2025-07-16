package sqlite

import (
	"context"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateUser(ctx context.Context, input *entity.User) (*entity.User, error) {
	if err := r.Db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindUserByAddress(ctx context.Context, address Address) (*entity.User, error) {
	var user entity.User
	if err := r.Db.WithContext(ctx).Where("address = ?", address).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by address: %w", err)
	}
	return &user, nil
}

func (r *SQLiteRepository) FindUsersByRole(ctx context.Context, role string) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.Db.WithContext(ctx).Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}
	return users, nil
}

func (r *SQLiteRepository) FindAllUsers(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.Db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}
	return users, nil
}

func (r *SQLiteRepository) DeleteUser(ctx context.Context, address Address) error {
	res := r.Db.WithContext(ctx).Where("address = ?", address).Delete(&entity.User{})
	if res.Error != nil {
		return fmt.Errorf("failed to delete user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}
	return nil
}
