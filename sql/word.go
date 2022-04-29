package sql

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/conlangdev/conlangdev"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

type WordService struct {
	db       *DB
	validate *validator.Validate
}

func NewWordService(db *DB, validate *validator.Validate) *WordService {
	return &WordService{db, validate}
}

func (s *WordService) GetWordByID(ctx context.Context, id uint) (*conlangdev.Word, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var word conlangdev.Word
	if err := tx.QueryRowContext(ctx,
		`SELECT
			id, uid, created_at, updated_at, headword, part_of_speech,
			definition, pronunciation, grammar_class, gender,
			etymology, notes, language_id
		FROM words WHERE id = ? LIMIT 1`,
		id,
	).Scan(
		&word.ID, &word.UID, &word.CreatedAt, &word.UpdatedAt, &word.Headword,
		&word.PartOfSpeech, &word.Definition, &word.Pronunciation, &word.GrammarClass,
		&word.Gender, &word.Etymology, &word.Notes, &word.LanguageID,
	); err == sql.ErrNoRows {
		return nil, &conlangdev.Error{
			Code:       conlangdev.ENOTFOUND,
			Message:    "could not find that language",
			StatusCode: http.StatusNotFound,
		}
	} else if err != nil {
		return nil, err
	}

	return &word, nil
}

func (s *WordService) GetWordByLanguageAndUID(ctx context.Context, language *conlangdev.Language, uid uint64) (*conlangdev.Word, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var word conlangdev.Word
	if err := tx.QueryRowContext(ctx,
		`SELECT
			id, uid, created_at, updated_at, headword, part_of_speech,
			definition, pronunciation, grammar_class, gender,
			etymology, notes, language_id
		FROM words WHERE uid = ? AND language_id = ? LIMIT 1`,
		uid, language.ID,
	).Scan(
		&word.ID, &word.UID, &word.CreatedAt, &word.UpdatedAt, &word.Headword,
		&word.PartOfSpeech, &word.Definition, &word.Pronunciation, &word.GrammarClass,
		&word.Gender, &word.Etymology, &word.Notes, &word.LanguageID,
	); err == sql.ErrNoRows {
		return nil, &conlangdev.Error{
			Code:       conlangdev.ENOTFOUND,
			Message:    "could not find that word",
			StatusCode: http.StatusNotFound,
		}
	} else if err != nil {
		return nil, err
	}

	return &word, nil
}

func (s *WordService) FindWordsForLanguage(ctx context.Context, language *conlangdev.Language) ([]*conlangdev.WordIndex, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx,
		`SELECT
			id, uid, headword, definition
		FROM words WHERE language_id = ?`,
		language.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	words := make([]*conlangdev.WordIndex, 0)
	for rows.Next() {
		var word conlangdev.WordIndex
		if err := rows.Scan(
			&word.ID, &word.UID, &word.Headword, &word.Definition,
		); err != nil {
			return nil, err
		}
		words = append(words, &word)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func (s *WordService) CreateWordForLanguage(ctx context.Context, language *conlangdev.Language, create conlangdev.WordCreate) (*conlangdev.Word, error) {
	if err := s.validate.Struct(&create); err != nil {
		if val_err, ok := err.(validator.ValidationErrors); ok {
			var fields []string
			for _, field := range val_err {
				fields = append(fields, field.Field())
			}
			return nil, &conlangdev.FieldsError{
				Code:       conlangdev.EVALIDFAIL,
				Message:    "validation failed",
				StatusCode: http.StatusBadRequest,
				Fields:     fields,
			}
		}
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	word := &conlangdev.Word{}
	if err := tx.QueryRowContext(ctx,
		`INSERT INTO words (
			created_at, updated_at, headword, part_of_speech, definition,
			pronunciation, grammar_class, gender, etymology, notes, language_id
		) VALUES (
			NOW(), NOW(), ?, ?, ?, ?, ?, ?, ?, ?, ?
		) RETURNING
			id, uid, created_at, updated_at, headword, part_of_speech,
			definition, pronunciation, grammar_class, gender, etymology,
			notes, language_id`,
		create.Headword, create.PartOfSpeech, create.Definition, create.Pronunciation,
		create.GrammarClass, create.Gender, create.Etymology, create.Notes,
		language.ID,
	).Scan(
		&word.ID, &word.UID, &word.CreatedAt, &word.UpdatedAt, &word.Headword,
		&word.PartOfSpeech, &word.Definition, &word.Pronunciation, &word.GrammarClass,
		&word.Gender, &word.Etymology, &word.Notes, &word.LanguageID,
	); err != nil {
		if sql_err, ok := err.(*mysql.MySQLError); ok {
			if sql_err.Number == 1062 {
				return nil, &conlangdev.Error{
					Code:       conlangdev.ECONFLICT,
					Message:    "something went wrong adding the word - try again",
					StatusCode: http.StatusConflict,
				}
			} else if sql_err.Number == 1452 {
				return nil, &conlangdev.Error{
					Code:       conlangdev.ENOTFOUND,
					Message:    "language with that ID does not exist",
					StatusCode: http.StatusNotFound,
				}
			}
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return word, nil
}

func (s *WordService) UpdateWord(ctx context.Context, word *conlangdev.Word, update conlangdev.WordUpdate) error {
	return errors.New("not implemented")
}

func (s *WordService) DeleteWord(ctx context.Context, word *conlangdev.Word) error {
	return errors.New("not implemented")
}
