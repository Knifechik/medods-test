package http_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	http_api "medods/cmd/login/internal/api/http"
	"net/http"
	"testing"
)

var (
	errAny = errors.New("any err")
)

func start(t *testing.T) (*Mockapplication, http.Handler, *require.Assertions) {
	ctrl := gomock.NewController(t)
	mockApp := NewMockapplication(ctrl)
	assert := require.New(t)

	return mockApp, http_api.New(mockApp), assert
}
