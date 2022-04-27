package sql

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/conlangdev/conlangdev"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db        *DB
	validate  *validator.Validate
	jwtSecret []byte
}

type customClaim struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func NewUserService(db *DB, validate *validator.Validate, jwtSecret string) *UserService {
	return &UserService{db, validate, []byte(jwtSecret)}
}

// Generates a string hash based on the given raw password.
func (s *UserService) generatePasswordHash(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Retrieves a single user from the database according to their unique
// integer ID.
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*conlangdev.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var user conlangdev.User
	if err := tx.QueryRowContext(ctx,
		`SELECT
			id, created_at, updated_at, username,
			email, display_name, password_hash
		FROM users WHERE id = ? LIMIT 1`,
		id,
	).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username,
		&user.Email, &user.DisplayName, &user.PasswordHash,
	); err == sql.ErrNoRows {
		return nil, &conlangdev.Error{
			Code:       conlangdev.ENOTFOUND,
			Message:    "could not find that user",
			StatusCode: http.StatusNotFound,
		}
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

// Retrieves a single user from the database according to their unique username.
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*conlangdev.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var user conlangdev.User
	if err := tx.QueryRowContext(ctx,
		`SELECT
			id, created_at, updated_at, username,
			email, display_name, password_hash
		FROM users WHERE username = ? LIMIT 1`,
		username,
	).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username,
		&user.Email, &user.DisplayName, &user.PasswordHash,
	); err == sql.ErrNoRows {
		return nil, &conlangdev.Error{
			Code:       conlangdev.ENOTFOUND,
			Message:    "could not find that user",
			StatusCode: http.StatusNotFound,
		}
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUserByJWT(ctx context.Context, tokenString string) (*conlangdev.User, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&customClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return s.jwtSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*customClaim)
	if !ok {
		return nil, errors.New("could not parse JWT claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("JWT token has expired")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var user conlangdev.User
	if err := tx.QueryRowContext(ctx,
		`SELECT
			id, created_at, updated_at, username,
			email, display_name, password_hash
		FROM users WHERE id = ? AND username = ?
		LIMIT 1`,
		claims.UserID, claims.Username,
	).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username,
		&user.Email, &user.DisplayName, &user.PasswordHash,
	); err != nil {
		return nil, errors.New("JWT could not be associated with a user")
	}

	return &user, nil
}

func (s *UserService) GenerateJWTForUser(ctx context.Context, user *conlangdev.User) (string, error) {
	claims := &customClaim{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "conlangdev",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// Creates a new user in the database according to the given parameters in the
// UserCreate DTO.
func (s *UserService) CreateUser(ctx context.Context, create conlangdev.UserCreate) (*conlangdev.User, error) {
	// Validate fields on create user DTO
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
	// Hash password
	hash, err := s.generatePasswordHash(ctx, create.Password)
	if err != nil {
		return nil, err
	}
	// Start a transaction, rolling back if anything goes wrong
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Insert the user into the database, scanning the inserted object back
	// into a new user object
	user := &conlangdev.User{}
	if err := tx.QueryRowContext(
		ctx,
		`INSERT INTO users (
			created_at, updated_at, username,
			email, display_name, password_hash
		) VALUES (NOW(), NOW(), ?, ?, ?, ?) RETURNING
			id, created_at, updated_at, username,
			email, display_name, password_hash`,
		create.Username, create.Email, create.DisplayName, hash,
	).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username,
		&user.Email, &user.DisplayName, &user.PasswordHash,
	); err != nil {
		if sql_err, ok := err.(*mysql.MySQLError); ok {
			if sql_err.Number == 1062 {
				return nil, &conlangdev.Error{
					Code:       conlangdev.ECONFLICT,
					Message:    "user already exists with that username or email",
					StatusCode: http.StatusConflict,
				}
			}
		}
		return nil, err
	}
	// Commit transaction!
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
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
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

func (s *UserService) GetViewForUser(ctx context.Context, user *conlangdev.User) (*conlangdev.UserView, error) {
	return &conlangdev.UserView{
		Username:    user.Username,
		DisplayName: user.DisplayName,
	}, nil
}
