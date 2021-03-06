package conlangdev

import (
	"context"
	"time"
)

type User struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	PasswordHash string    `json:"-"`
}

type UserUpdate struct {
	DisplayName string `json:"display_name"`
}

type UserCreate struct {
	Username    string `json:"username" validate:"min=3,max=32"`
	Email       string `json:"email" validate:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password" validate:"min=8"`
}

type UserView struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type UserService interface {
	GetUserByID(ctx context.Context, id uint) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByJWT(ctx context.Context, jwt string) (*User, error)
	GenerateJWTForUser(ctx context.Context, user *User) (string, error)
	CreateUser(ctx context.Context, create UserCreate) (*User, error)
	UpdateUser(ctx context.Context, user *User, update UserUpdate) error
	DeleteUser(ctx context.Context, user *User) error
	UpdateUserPassword(ctx context.Context, user *User, password string) error
	CheckUserPassword(ctx context.Context, user *User, password string) error
	GetViewForUser(ctx context.Context, user *User) (*UserView, error)
}
