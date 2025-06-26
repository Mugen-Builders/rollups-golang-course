package entity

import (
	"errors"
	"fmt"

	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
)

var (
	ErrInvalidUser  = errors.New("invalid user")
	ErrUserNotFound = errors.New("user not found")
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleDebtor   UserRole = "debtor"
	UserRoleInvestor UserRole = "investor"
)

type User struct {
	Id        uint     `json:"id" gorm:"primaryKey"`
	Role      UserRole `json:"role,omitempty" gorm:"not null"`
	Address   Address  `json:"address,omitempty" gorm:"custom_type:text;uniqueIndex;not null"`
	CreatedAt int64    `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt int64    `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewUser(role string, address Address, createdAt int64) (*User, error) {
	user := &User{
		Role:      UserRole(role),
		Address:   address,
		CreatedAt: createdAt,
	}
	if err := user.validate(); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) validate() error {
	if u.Role == "" {
		return fmt.Errorf("%w: role cannot be empty", ErrInvalidUser)
	}
	if u.Role != UserRoleAdmin && u.Role != UserRoleDebtor && u.Role != UserRoleInvestor {
		return fmt.Errorf("%w: invalid role", ErrInvalidUser)
	}
	if u.Address == (Address{}) {
		return fmt.Errorf("%w: address cannot be empty", ErrInvalidUser)
	}
	if u.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidUser)
	}
	return nil
}
