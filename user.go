package conlangdev

import (
	"context"
	"time"
)

type User struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Username     string    `json:"username" validate:"min=3,max=32"`
	Email        string    `json:"email" validate:"email"`
	DisplayName  string    `json:"display_name"`
	PasswordHash string    `json:"-" validate:"required"`
}

type UserUpdate struct {
	DisplayName string `json:"display_name"`
}

type UserService interface {
	GetUserByID(ctx context.Context, id uint) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByJWT(ctx context.Context, jwt string) (*User, error)
	GenerateJWTForUser(ctx context.Context, user *User) (string, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User, update UserUpdate) error
	DeleteUser(ctx context.Context, user *User) error
	UpdateUserPassword(ctx context.Context, user *User, password string) error
	CheckUserPassword(ctx context.Context, user *User, password string) error
}
