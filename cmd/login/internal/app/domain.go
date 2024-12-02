package app

import (
	"github.com/gofrs/uuid"
	"time"
)

type (
	Token struct {
		AccessToken  string
		RefreshToken string
	}
	TokensTable struct {
		ID        uuid.UUID
		UserID    uuid.UUID
		Hash      []byte
		IP        string
		CreatedAt time.Time
		UpdatedAt time.Time
		ExpiresAt time.Time
	}
)

const (
	AccessExpire  = time.Minute * 15
	RefreshExpire = time.Hour * 24 * 30
)
