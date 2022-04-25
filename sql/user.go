package sql

import (
	"context"
	"errors"

	"github.com/conlangdev/conlangdev"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db}
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*conlangdev.User, error) {
	return nil, errors.New("not implemented")
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*conlangdev.User, error) {
	return nil, errors.New("not implemented")
}

func (s *UserService) GetUserByJWT(ctx context.Context, jwt string) (*conlangdev.User, error) {
	return nil, errors.New("not implemented")
}

func (s *UserService) GenerateJWTForUser(ctx context.Context, user *conlangdev.User) (string, error) {
	return "", errors.New("not implemented")
}

func (s *UserService) CreateUser(ctx context.Context, user *conlangdev.User) error {
	return errors.New("not implemented")
}

func (s *UserService) UpdateUser(ctx context.Context, user *conlangdev.User, update conlangdev.UserUpdate) error {
	return errors.New("not implemented")
}

func (s *UserService) DeleteUser(ctx context.Context, user *conlangdev.User) error {
	return errors.New("not implemented")
}

func (s *UserService) UpdateUserPassword(ctx context.Context, user *conlangdev.User, password string) error {
	return errors.New("not implemented")
}

func (s *UserService) CheckUserPassword(ctx context.Context, user *conlangdev.User, password string) error {
	return errors.New("not implemented")
}
