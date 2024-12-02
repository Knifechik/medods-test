package app

import (
	"context"
	"github.com/gofrs/uuid"
)

type (
	Repo interface {
		SaveToken(ctx context.Context, refreshToken TokensTable) error
		TokenByUserID(ctx context.Context, userID uuid.UUID) (*TokensTable, error)
	}
	Auth interface {
		GenerateTokens(ip, userID string) (*Token, error)
		DecodeJWT(accessToken string) (uuid.UUID, string, error)
		DecodeRefreshToken(refreshToken string) (string, error)
	}
	Notification interface {
		SendNotify(ctx context.Context, userID uuid.UUID) error
	}
)
