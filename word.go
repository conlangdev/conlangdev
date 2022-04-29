package conlangdev

import (
	"context"
	"time"
)

type Word struct {
	ID            uint      `json:"id"`
	UID           uint64    `json:"uid"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Headword      string    `json:"headword"`
	PartOfSpeech  string    `json:"part_of_speech"`
	Definition    string    `json:"definition"`
	Pronunciation string    `json:"pronunciation"`
	GrammarClass  string    `json:"grammar_class"`
	Gender        string    `json:"gender"`
	Etymology     string    `json:"etymology"`
	Notes         string    `json:"notes"`
	LanguageID    uint      `json:"language_id"`
}

type WordUpdate struct {
	Headword      string `json:"headword"`
	PartOfSpeech  string `json:"part_of_speech"`
	Definition    string `json:"definition"`
	Pronunciation string `json:"pronunciation"`
	GrammarClass  string `json:"grammar_class"`
	Gender        string `json:"gender"`
	Etymology     string `json:"etymology"`
	Notes         string `json:"notes"`
}

type WordCreate struct {
	Headword      string `json:"headword" validate:"required"`
	PartOfSpeech  string `json:"part_of_speech" validate:"required"`
	Definition    string `json:"definition" validate:"required"`
	Pronunciation string `json:"pronunciation"`
	GrammarClass  string `json:"grammar_class"`
	Gender        string `json:"gender"`
	Etymology     string `json:"etymology"`
	Notes         string `json:"notes"`
}

type WordIndex struct {
	ID         uint   `json:"id"`
	UID        uint64 `json:"uid"`
	Headword   string `json:"headword"`
	Definition string `json:"definition"`
}

type WordService interface {
	GetWordByID(ctx context.Context, id uint) (*Word, error)
	GetWordByLanguageAndUID(ctx context.Context, language *Language, uid uint64) (*Word, error)
	FindWordsForLanguage(ctx context.Context, language *Language) ([]*WordIndex, error)
	CreateWordForLanguage(ctx context.Context, language *Language, create WordCreate) (*Word, error)
	UpdateWord(ctx context.Context, word *Word, update WordUpdate) error
	DeleteWord(ctx context.Context, word *Word) error
}
