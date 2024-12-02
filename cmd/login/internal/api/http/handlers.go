package http

import (
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	ErrInvalidUserID = errors.New("invalid userID")
	ErrInvalidIP     = errors.New("invalid IP")
	ErrInvalidCookie = errors.New("invalid cookie")
	ErrServerError   = errors.New("server error")
)

// loginHandler create a new pair of tokens
func (a *api) loginHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.FromString(vars["id"])
	if err != nil {
		log.Println(err)
		errorHandler(w, http.StatusBadRequest, ErrInvalidUserID)
		return
	}
	ip, err := getIP(r)
	if err != nil {
		log.Println(err)
		errorHandler(w, http.StatusBadRequest, ErrInvalidIP)
		return
	}
	token, err := a.app.Login(r.Context(), userID, ip)
	if err != nil {
		log.Println(err)
		errorHandler(w, http.StatusUnauthorized, ErrServerError)
		return
	}

	accessCookie := BuildCookieAccessToken(token.AccessToken)
	http.SetCookie(w, accessCookie)

	refreshTokenCookie := BuildCookieRefreshToken(token.RefreshToken)
	http.SetCookie(w, refreshTokenCookie)

	w.WriteHeader(http.StatusOK)
}

// refreshHandler receives a pair of tokens, checks it and creates new pair
func (a *api) refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie(CookieNameRefreshToken)
	if err != nil {
		log.Println(err)
		errorHandler(w, http.StatusBadRequest, ErrInvalidCookie)
		return
	}
	accessCookie, err := r.Cookie(CookieNameAccessToken)
	if err != nil {
		log.Println(err)
		errorHandler(w, http.StatusBadRequest, ErrInvalidCookie)
		return
	}

	ip, err := getIP(r)
	if err != nil {
		log.Println(err)
		errorHandler(w, http.StatusBadRequest, ErrInvalidIP)
		return
	}

	tokens, err := a.app.Refresh(r.Context(), accessCookie.Value, refreshCookie.Value, ip)
	if err != nil {
		log.Println(err)
		errorHandler(w, http.StatusInternalServerError, ErrServerError)
		return
	}

	newAccessCookie := BuildCookieAccessToken(tokens.AccessToken)
	http.SetCookie(w, newAccessCookie)

	newRefreshTokenCookie := BuildCookieRefreshToken(tokens.RefreshToken)
	http.SetCookie(w, newRefreshTokenCookie)

	w.WriteHeader(http.StatusOK)

}

// errorHandler for error response
func errorHandler(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)

	erR := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if erR != nil {
		log.Println(erR)
	}
}
