package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Login     string    `gorm:"unique;not null"`
	PassHash  []byte    `gorm:"not null"`
	CreatedAt time.Time
	Faculty   string
	Direction string
	Histories []History `gorm:"foreignKey:UserID"`
}

type UserDTO struct {
	Login string `json:"login"`
	Pass  string `json:"password"`
}
type UserInteractor interface {
	CreateUser(ctx context.Context, login string, pass string) (uuid.UUID, error)
	Login(ctx context.Context, login string, passhash string) (string, error)
	User(ctx context.Context, id uuid.UUID) (*User, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (uuid.UUID, error)
	User(ctx context.Context, id uuid.UUID) (*User, error)
	UserByLogin(ctx context.Context, login string) (*User, error)
}
