package sql

import (
	"context"
	"errors"

	"github.com/conlangdev/conlangdev"
	"github.com/go-playground/validator/v10"
)

type WordService struct {
	db       *DB
	validate *validator.Validate
}

func NewWordService(db *DB, validate *validator.Validate) *WordService {
	return &WordService{db, validate}
}

func (s *WordService) GetWordByID(ctx context.Context, int uint) (*conlangdev.Word, error) {
	return nil, errors.New("not implemented")
}

func (s *WordService) GetWordByLanguageAndSlug(ctx context.Context, language *conlangdev.Language, slug string) (*conlangdev.Word, error) {
	return nil, errors.New("not implemented")
}

func (s *WordService) CreateWord(ctx context.Context, word *conlangdev.Word) error {
	return errors.New("not implemented")
}

func (s *WordService) UpdateWord(ctx context.Context, word *conlangdev.Word, update conlangdev.WordUpdate) error {
	return errors.New("not implemented")
}

func (s *WordService) DeleteWord(ctx context.Context, word *conlangdev.Word) error {
	return errors.New("not implemented")
}
