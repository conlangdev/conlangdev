package sql

import (
	"context"
	"errors"

	"github.com/conlangdev/conlangdev"
	"github.com/go-playground/validator/v10"
)

type LanguageService struct {
	db       *DB
	validate *validator.Validate
}

func NewLanguageService(db *DB, validate *validator.Validate) *LanguageService {
	return &LanguageService{db, validate}
}

func (s *LanguageService) GetLanguageByID(ctx context.Context, id uint) (*conlangdev.Language, error) {
	return nil, errors.New("not implemented")
}

func (s *LanguageService) GetLanguageByUserAndSlug(ctx context.Context, user *conlangdev.User, slug string) (*conlangdev.Language, error) {
	return nil, errors.New("not implemented")
}

func (s *LanguageService) CreateLanguage(ctx context.Context, language *conlangdev.Language) error {
	return errors.New("not implemented")
}

func (s *LanguageService) UpdateLanguage(ctx context.Context, language *conlangdev.Language, update conlangdev.LanguageUpdate) error {
	return errors.New("not implemented")
}

func (s *LanguageService) DeleteLanguage(ctx context.Context, language *conlangdev.Language) error {
	return errors.New("not implemented")
}
