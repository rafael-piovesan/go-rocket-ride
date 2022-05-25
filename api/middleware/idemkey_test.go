//go:build unit
// +build unit

package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/v2/api/context"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/httpserver"
	"github.com/stretchr/testify/assert"
)

func TestIdemKey(t *testing.T) {
	e := httpserver.New()
	e.Use(IdempotencyKey())

	payload := "{\"key\":\"value\"}"

	e.POST("/", func(c echo.Context) error {
		ik, ok := context.GetIdemKey(c)

		assert.True(t, ok)
		assert.Equal(t, http.MethodPost, ik.RequestMethod)
		assert.Equal(t, "/", ik.RequestPath)
		assert.Equal(t, json.RawMessage(payload), ik.RequestParams)

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

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Body = ioutil.NopCloser(strings.NewReader(payload))

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
