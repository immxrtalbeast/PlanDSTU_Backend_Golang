package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"github.com/immxrtalbeast/plandstu/internal/lib"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("user not found")
)

type UserInteractor struct {
	userRepo  domain.UserRepository
	tokenTTL  time.Duration
	appSecret string
}

func NewUserInteractor(userRepo domain.UserRepository, tokenTTL time.Duration, appSecret string) *UserInteractor {
	return &UserInteractor{
		userRepo:  userRepo,
		tokenTTL:  tokenTTL,
		appSecret: appSecret,
	}
}

func (ui *UserInteractor) CreateUser(ctx context.Context, login string, pass string, group string) (uuid.UUID, error) {
	const op = "uc.user.create"
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	user := domain.User{
		Login:    login,
		PassHash: passHash,
		Group:    group,
	}
	id, err := ui.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (ui *UserInteractor) Login(ctx context.Context, login string, passhash string) (string, error) {
	const op = "uc.user.login"
	user, err := ui.userRepo.UserByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(passhash)); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	return lib.NewToken(user, ui.tokenTTL, ui.appSecret)

}

func (ui *UserInteractor) User(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const op = "uc.user.get"
	user, err := ui.userRepo.User(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		} else {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return user, nil
}
