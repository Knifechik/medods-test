package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"medods/cmd/login/internal/app"
	"strings"
	"time"
)

var _ app.Auth = &Auth{}

type (
	Auth struct {
		AccessKey []byte
	}
	RefreshCrypt struct {
		IP  string
		Key string
	}
)

// New creates and returns new instance auth.
func New(accessSecretKey string) *Auth {
	return &Auth{
		AccessKey: []byte(accessSecretKey),
	}
}

func (a *Auth) GenerateTokens(ip, userID string) (*app.Token, error) {
	id := idForPair()
	accessToken, err := a.GenerateAccessToken(id, ip, userID)
	if err != nil {
		return nil, fmt.Errorf("a.GenerateAccessToken: %w", err)
	}

	refreshToken, err := a.GenerateRefreshToken(id, ip)
	if err != nil {
		return nil, fmt.Errorf("a.GenerateRefreshToken: %w", err)
	}

	return &app.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (a *Auth) GenerateAccessToken(id, ip, userID string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":  userID,
		"ip":   ip,
		"pair": id,
		"exp":  time.Now().Add(app.AccessExpire).Unix(),
	})

	token, err := claims.SignedString(a.AccessKey)
	if err != nil {
		return "", fmt.Errorf("claims.SignedString: %w", err)
	}

	return token, nil
}

func (a *Auth) GenerateRefreshToken(id, ip string) (string, error) {
	crypto := RefreshCrypt{
		IP:  ip,
		Key: id,
	}

	data, err := json.Marshal(crypto)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}

	refreshToken := base64.StdEncoding.EncodeToString(data)

	return refreshToken, nil
}

func (a *Auth) DecodeJWT(accessToken string) (uuid.UUID, string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return a.AccessKey, nil
	})
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("jwt.Parse: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, "", fmt.Errorf("jwt.MapClaims is not a MapClaims")
	}

	userIDStr := claims["sub"].(string)

	accessPair := claims["pair"].(string)

	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("uuid.FromString: %w", err)
	}

	return userID, accessPair, nil
}

func (a *Auth) DecodeRefreshToken(refreshToken string) (string, error) {
	decodedRefresh, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return "", fmt.Errorf("base64.DecodeString: %w", err)
	}

	var unmarshalRefresh RefreshCrypt
	err = json.Unmarshal(decodedRefresh, &unmarshalRefresh)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	return unmarshalRefresh.Key, nil

}

var chars = []rune("abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789" + "_")

func idForPair() string {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))

	length := 8

	var b strings.Builder

	for i := 0; i < length; i++ {
		b.WriteRune(chars[rd.Intn(len(chars))])
	}

	str := b.String()

	return str
}
