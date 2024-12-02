package http

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"medods/cmd/login/internal/app"
	"net/http"
)

type application interface {
	Login(ctx context.Context, userID uuid.UUID, ip string) (*app.Token, error)
	Refresh(ctx context.Context, accessToken string, refreshToken string, ip string) (*app.Token, error)
}

type api struct {
	app application
}

func New(application application) http.Handler {
	api := &api{
		app: application,
	}

	r := mux.NewRouter()

	r.HandleFunc("/login/{id}", api.loginHandler).Methods(http.MethodPost)
	r.HandleFunc("/refresh", api.refreshHandler).Methods(http.MethodPost)

	return r
}

const (
	CookieNameRefreshToken = `cookie_refresh`
	CookieNameAccessToken  = `access_token`
)

func BuildCookieAccessToken(value string) *http.Cookie {
	return buildCookie(CookieNameAccessToken, value)
}

func BuildCookieRefreshToken(value string) *http.Cookie {
	return buildCookie(CookieNameRefreshToken, value)
}

func buildCookie(name string, value string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}

	return cookie
}
