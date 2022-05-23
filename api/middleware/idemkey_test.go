//go:build unit
// +build unit

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestIdemKey(t *testing.T) {
	e := echo.New()
	e.Use(IdempotencyKey())
	e.POST("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	tests := []struct {
		fail   bool
		header string
		value  string
	}{
		{fail: true, header: gofakeit.AnimalType(), value: gofakeit.LetterN(uint(gofakeit.Number(101, 1000)))},
		{fail: true, header: gofakeit.AnimalType(), value: gofakeit.LetterN(uint(gofakeit.Number(0, 100)))},
		{fail: true, header: "idempotency-key", value: gofakeit.LetterN(uint(gofakeit.Number(101, 1000)))},
		{fail: false, header: "idempotency-key", value: gofakeit.LetterN(uint(gofakeit.Number(0, 100)))},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()

		// testing different header values
		req.Header.Set(tc.header, tc.value)

		e.ServeHTTP(rec, req)

		if tc.fail {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		} else {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	}
}
