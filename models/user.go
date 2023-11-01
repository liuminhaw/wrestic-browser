package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Username     string
	Email        sql.NullString
	PasswordHash string
	Role         string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(username, password, role, email string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	hashedPassword := string(hashedBytes)

	var row *sql.Row
	var emailField sql.NullString
	if email == "" {
		row = us.DB.QueryRow(`
			INSERT INTO users (username, password_hash, role_id)
			VALUES ($1, $2, (SELECT id FROM roles WHERE role = $3));
		`, username, hashedPassword, role)
		emailField = sql.NullString{String: email, Valid: false}
	} else {
		row = us.DB.QueryRow(`
			INSERT INTO users (username, password_hash, email, role_id)
			VALUES ($1, $2, $3, (SELECT id FROM roles WHERE role = $4));
		`, username, hashedPassword, email, role)
		emailField = sql.NullString{String: email, Valid: true}
	}

	user := User{
		Username:     username,
		Email:        emailField,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	err = row.Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrUsernameTaken
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (us *UserService) Authenticate(username, password string) (*User, error) {
	username = strings.ToLower(username)
	user := User{
		Username: username,
	}
	row := us.DB.QueryRow(`
		SELECT id, password_hash FROM users WHERE username = $1;
	`, user.Username)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate user: %w", err)
	}
	return &user, nil
}
