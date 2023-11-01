package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/liuminhaw/wrestic-brw/rand"
)

const (
	// The minimum number of bytes to be used for each session token
	MinBytesPerToken = 32
)

// Token is only set when creating a new session
// When look up a session this will be left empty,
// as we only store the hash of a session token in out database
// and we cannot reverse it into a raw token
type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each session token. If this value is not set or is less than the
	// MinBytesePerToken const it will be ignored and MinBytesPerToken will be used.
	BytesPerToken int
}

func (service *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: service.hash(token),
	}
	row := service.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2
		RETURNING id;
	`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (service *SessionService) Delete(token string) error {
	tokenHash := service.hash(token)
	_, err := service.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1
	`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (service *SessionService) User(token string) (*User, error) {
	// 1. Hash the session token
	tokenHash := service.hash(token)
	// 2. Query for the session with that hash
	var user User
	row := service.DB.QueryRow(`
	SELECT users.id, users.username, users.password_hash
	FROM sessions
	JOIN users ON users.id = sessions.user_id
	WHERE sessions.token_hash = $1
	`, tokenHash)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	// 4. return user
	return &user, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
