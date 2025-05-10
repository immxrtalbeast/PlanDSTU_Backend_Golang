package psql

import (
	"context"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return user.ID, nil
}

func (r *UserRepository) User(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepository) UserByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("login = ?", login).First(&user).Error
	return &user, err
}
