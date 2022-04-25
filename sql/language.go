package sql

import (
	"context"
	"errors"

	"github.com/conlangdev/conlangdev"
)

type LanguageService struct {
	db *DB
}

func NewLanguageService(db *DB) *LanguageService {
	return &LanguageService{db}
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
