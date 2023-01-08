package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)



const (
	ScopeAuthentication = "authentication"
)

// Token is the type for authentication tokens
type Token struct {
	PlainText string    `json:"token"`
	UserID    int64     `json:"-"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}



func GenerateToken(userID int, timeToLive time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: int64(userID),
		Expiry: time.Now().Add(timeToLive),
		Scope: scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	
	hash := sha256.Sum256([]byte(token.PlainText))
	
	token.Hash = hash[:]

	return token, nil
}

func (m *DBModel) InsertToken(t *Token, u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	stmt := `delete from tokens where user_id = $1`
	_, err := m.Db.QueryContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	stmt = `
		insert 
			into tokens 
			(user_id, name, token_hash, expiry, 
			created_at, updated_at) 
		values 
			($1, $2, $3, $4, $5, $6, $7)
		`

	_, err = m.Db.ExecContext(ctx, stmt, u.ID, u.Username, t.Hash, t.Expiry, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) 