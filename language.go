package conlangdev

import (
	"context"
	"time"
)

type Language struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name" validate:"required"`
	Slug      string    `json:"slug" validate:"required"`
	Endonym   string    `json:"endonym"`
	UserID    uint      `json:"user_id" validate:"required:"`
}

type LanguageUpdate struct {
	Name    string `json:"name"`
	Endonym string `json:"endonym"`
}

type LanguageService interface {
	GetLanguageByID(ctx context.Context, id uint) (*Language, error)
	CreateLanguage(ctx context.Context, language *Language) error
	UpdateLanguage(ctx context.Context, language *Language, update LanguageUpdate) error
	DeleteLanguage(ctx context.Context, language *Language) error
}
