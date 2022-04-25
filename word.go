package conlangdev

import (
	"context"
	"time"
)

type Word struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Headword   string    `json:"headword" validate:"required"`
	Slug       string    `json:"slug" validate:"required"`
	Definition string    `json:"definition" validate:"required"`
	Etymology  string    `json:"etymology"`
	Notes      string    `json:"notes"`
	LanguageID uint      `json:"language_id" validate:"required"`
}

type WordUpdate struct {
	Headword   string `json:"headword"`
	Slug       string `json:"slug"`
	Definition string `json:"definition"`
	Etymology  string `json:"etymology"`
	Notes      string `json:"notes"`
}

type WordService interface {
	GetWordByID(ctx context.Context, int uint) (*Word, error)
	GetWordByLanguageAndSlug(ctx context.Context, language *Language, slug string) (*Word, error)
	CreateWord(ctx context.Context, word *Word) error
	UpdateWord(ctx context.Context, word *Word, update WordUpdate) error
	DeleteWord(ctx context.Context, word *Word) error
}
