package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Login create a new pair of tokens
func (a *App) Login(ctx context.Context, userID uuid.UUID, ip string) (*Token, error) {

	tokens, err := a.auth.GenerateTokens(ip, userID.String())
	if err != nil {
		return nil, fmt.Errorf("a.auth.GenerateTokens: %w", err)
	}

	hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(tokens.RefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	err = a.repo.SaveToken(ctx, TokensTable{
		UserID:    userID,
		Hash:      hashedRefresh,
		IP:        ip,
		ExpiresAt: time.Now().Add(RefreshExpire),
	})
	if err != nil {
		return nil, fmt.Errorf("a.repo.SaveToken: %w", err)
	}

	return tokens, nil
}

// Refresh receives a pair of tokens, checks it and creates new pair
func (a *App) Refresh(ctx context.Context, accessToken string, refreshToken string, ip string) (*Token, error) {
	userID, accessPair, err := a.auth.DecodeJWT(accessToken)
	if err != nil {
		return nil, fmt.Errorf("a.auth.DecodeJWT: %w", err)
	}

	dbToken, err := a.repo.TokenByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("a.repo.TokenByUserID: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(dbToken.Hash, []byte(refreshToken))
	if err != nil {
		return nil, fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
	}

	if time.Now().After(dbToken.ExpiresAt) {
		return nil, errors.New("refresh token is expired")
	}

	if dbToken.IP != ip {
		err = a.notify.SendNotify(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("notify.SendNotify: %w", err)
		}

	}
	refreshPair, err := a.auth.DecodeRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("a.auth.DecodeRefreshToken: %w", err)
	}

	if refreshPair != accessPair {
		return nil, fmt.Errorf("refreshPair != accessPair")
	}

	tokens, err := a.auth.GenerateTokens(ip, userID.String())
	if err != nil {
		return nil, fmt.Errorf("a.auth.GenerateTokens: %w", err)
	}

	hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(tokens.RefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	err = a.repo.SaveToken(ctx, TokensTable{
		UserID:    userID,
		Hash:      hashedRefresh,
		IP:        ip,
		ExpiresAt: time.Now().Add(RefreshExpire),
	})
	if err != nil {
		return nil, fmt.Errorf("a.repo.SaveToken: %w", err)
	}

	return tokens, nil
}
