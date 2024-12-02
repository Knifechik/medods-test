package repo

import (
	"github.com/gofrs/uuid"
	"time"
)

type RefreshToken struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"hash_refresh_token"`
	IP           string    `db:"ip"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	ExpiresAt    time.Time `db:"expires_at"`
}
