package http_test

import (
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	http_api "medods/cmd/login/internal/api/http"
	"medods/cmd/login/internal/app"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttp_loginHandler(t *testing.T) {
	t.Parallel()

	var (
		accessToken  = "1234"
		refreshToken = "5678"
		ip           = "127.0.0.1"
		userID       = uuid.Must(uuid.NewV4())

		ok          = http.StatusOK
		serverError = http.StatusInternalServerError
		appRes      = &app.Token{
			AccessToken:  "1234",
			RefreshToken: "5678",
		}
	)

	testcases := map[string]struct {
		appRes *app.Token
		appErr error
		code   int
	}{
		"success":     {appRes, nil, ok},
		"a.app.login": {nil, errAny, serverError},
	}
	for name, tc := range testcases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockApp, c, assert := start(t)

			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/login/"+userID.String(), nil)
			response := httptest.NewRecorder()

			req.RemoteAddr = ip

			mockApp.EXPECT().Login(gomock.Any(), userID, ip).Return(tc.appRes, tc.appErr)

			c.ServeHTTP(response, req)
			assert.Equal(tc.code, response.Code)

			if tc.appErr == nil {

				res := response.Result()

				cookies := res.Cookies()
				for _, cookie := range cookies {
					if cookie.Name == http_api.CookieNameAccessToken {
						assert.Equal(accessToken, cookie.Value)
						continue
					}
					if cookie.Name == http_api.CookieNameRefreshToken {
						assert.Equal(refreshToken, cookie.Value)
						continue
					}
					t.Error("not cookies")
				}
			}
		})
	}
}

func TestHttp_refreshHandler(t *testing.T) {
	t.Parallel()

	var (
		accessToken  = "1234"
		refreshToken = "5678"
		ip           = "127.0.0.1"

		ok          = http.StatusOK
		serverError = http.StatusInternalServerError
		appRes      = &app.Token{
			AccessToken:  "1234",
			RefreshToken: "5678",
		}
		accessCookie = &http.Cookie{
			Name:     http_api.CookieNameAccessToken,
			Value:    accessToken,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		refreshCookie = &http.Cookie{
			Name:     http_api.CookieNameRefreshToken,
			Value:    refreshToken,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
	)

	testcases := map[string]struct {
		appRes *app.Token
		appErr error
		code   int
	}{
		"success":     {appRes, nil, ok},
		"a.app.login": {nil, errAny, serverError},
	}
	for name, tc := range testcases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockApp, c, assert := start(t)

			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/refresh", nil)
			response := httptest.NewRecorder()

			req.AddCookie(accessCookie)
			req.AddCookie(refreshCookie)

			req.RemoteAddr = ip

			mockApp.EXPECT().Refresh(gomock.Any(), accessCookie.Value, refreshCookie.Value, ip).Return(tc.appRes, tc.appErr)

			c.ServeHTTP(response, req)
			assert.Equal(tc.code, response.Code)

			if tc.appErr == nil {

				res := response.Result()

				cookies := res.Cookies()
				for _, cookie := range cookies {
					if cookie.Name == http_api.CookieNameAccessToken {
						assert.Equal(accessToken, cookie.Value)
						continue
					}
					if cookie.Name == http_api.CookieNameRefreshToken {
						assert.Equal(refreshToken, cookie.Value)
						continue
					}
					t.Error("not cookies")
				}
			}
		})
	}
}
