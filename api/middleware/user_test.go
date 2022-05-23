//go:build unit
// +build unit

package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	mocks "github.com/rafael-piovesan/go-rocket-ride/v2/mocks/datastore"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/httpserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUser(t *testing.T) {
	mockUser := &mocks.User{}

	e := httpserver.New()
	e.Use(User(mockUser))

	e.POST("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	tests := []struct {
		mock   bool
		fail   bool
		header string
		value  string
		ret    int
	}{
		{mock: false, fail: true, header: "", value: "", ret: http.StatusBadRequest},
		{mock: false, fail: true, header: gofakeit.AnimalType(), value: gofakeit.Email(), ret: http.StatusBadRequest},
		{mock: true, fail: true, header: "authorization", value: gofakeit.Email(), ret: http.StatusUnauthorized},
		{mock: true, fail: false, header: "authorization", value: gofakeit.Email(), ret: http.StatusOK},
	}

	for _, tc := range tests {
		if tc.mock {
			if tc.fail {
				mockUser.On("FindOne", mock.Anything, mock.Anything).
					Once().
					Return(entity.User{}, errors.New("it failed"))
			} else {
				mockUser.On("FindOne", mock.Anything, mock.Anything).Once().Return(entity.User{}, nil)
			}
		}

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()

		// testing different header values
		if tc.header != "" {
			req.Header.Set(tc.header, tc.value)
		}

		e.ServeHTTP(rec, req)
		assert.Equal(t, tc.ret, rec.Code)
	}
}
