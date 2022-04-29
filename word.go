package conlangdev

import (
	"context"
	"time"
)

type Word struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Headword   string    `json:"headword"`
	Definition string    `json:"definition"`
	Etymology  string    `json:"etymology"`
	Notes      string    `json:"notes"`
	LanguageID uint      `json:"language_id"`
}

type WordUpdate struct {
	Headword   string `json:"headword"`
	Definition string `json:"definition"`
	Etymology  string `json:"etymology"`
	Notes      string `json:"notes"`
}

type WordCreate struct {
	Headword   string `json:"headword"`
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
