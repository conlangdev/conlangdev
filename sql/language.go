package sql

import (
	"context"
	"database/sql"
	"net/http"

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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var language conlangdev.Language
	if err := tx.QueryRowContext(ctx,
		`SELECT
			id, created_at, updated_at, name,
			slug, endonym, user_id
		FROM languages WHERE id = ? LIMIT 1`,
		id,
	).Scan(
		&language.ID, &language.CreatedAt, &language.UpdatedAt,
		&language.Name, &language.Slug, &language.Endonym, &language.UserID,
	); err == sql.ErrNoRows {
		return nil, &conlangdev.Error{
			Code:       conlangdev.ENOTFOUND,
			Message:    "could not find that language",
			StatusCode: http.StatusNotFound,
		}
	} else if err != nil {
		return nil, err
	}

	return &language, nil
}

func (s *LanguageService) GetLanguageByUserAndSlug(ctx context.Context, user *conlangdev.User, slug string) (*conlangdev.Language, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var language conlangdev.Language
	if err := tx.QueryRowContext(ctx,
		`SELECT
			id, created_at, updated_at, name,
			slug, endonym, user_id
		FROM languages WHERE
			slug = ? AND user_id = ?
		LIMIT 1`,
		slug, user.ID,
	).Scan(
		&language.ID, &language.CreatedAt, &language.UpdatedAt,
		&language.Name, &language.Slug, &language.Endonym, &language.UserID,
	); err == sql.ErrNoRows {
		return nil, &conlangdev.Error{
			Code:       conlangdev.ENOTFOUND,
			Message:    "could not find that language",
			StatusCode: http.StatusNotFound,
		}
	} else if err != nil {
		return nil, err
	}

	return &language, nil
}

func (s *LanguageService) FindLanguagesForUser(ctx context.Context, user *conlangdev.User) ([]*conlangdev.Language, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx,
		`SELECT
			id, created_at, updated_at, name,
			slug, endonym, user_id
		FROM languages WHERE user_id = ?`,
		user.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	languages := make([]*conlangdev.Language, 0)
	for rows.Next() {
		var language conlangdev.Language
		if err := rows.Scan(
			&language.ID, &language.CreatedAt, &language.UpdatedAt,
			&language.Name, &language.Slug, &language.Endonym, &language.UserID,
		); err != nil {
			return nil, err
		}
		languages = append(languages, &language)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return languages, nil
}

func (s *LanguageService) CreateLanguageForUser(ctx context.Context, user *conlangdev.User, create conlangdev.LanguageCreate) (*conlangdev.Language, error) {
	return nil, &conlangdev.Error{
		Code:       conlangdev.ENOTIMPLEMENTED,
		Message:    "not implemented",
		StatusCode: http.StatusInternalServerError,
	}
}

func (s *LanguageService) UpdateLanguage(ctx context.Context, language *conlangdev.Language, update conlangdev.LanguageUpdate) error {
	return &conlangdev.Error{
		Code:       conlangdev.ENOTIMPLEMENTED,
		Message:    "not implemented",
		StatusCode: http.StatusInternalServerError,
	}
}

func (s *LanguageService) DeleteLanguage(ctx context.Context, language *conlangdev.Language) error {
	return &conlangdev.Error{
		Code:       conlangdev.ENOTIMPLEMENTED,
		Message:    "not implemented",
		StatusCode: http.StatusInternalServerError,
	}
}
