//go:build unit
// +build unit

package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/entity/idempotency"
	"github.com/rafael-piovesan/go-rocket-ride/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	e := echo.New()
	uc := &mocks.RideUseCase{}
	user := entity.User{
		ID:               gofakeit.Int64(),
		Email:            gofakeit.Email(),
		StripeCustomerID: gofakeit.UUID(),
	}
	handler := NewRideHandler(uc)

	callArgs := []interface{}{
		mock.Anything,
		mock.AnythingOfType("*entity.IdempotencyKey"),
		mock.AnythingOfType("*entity.Ride"),
	}

	t.Run("User info not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// context is missing the "user-id" key
		c := e.NewContext(req, rec)

		err := handler.Create(c)

		if assert.Error(t, err) && assert.IsType(t, (&echo.HTTPError{}), err) {
			he := (err).(*echo.HTTPError)
			assert.Equal(t, http.StatusUnauthorized, he.Code)
			assert.Equal(t, entity.ErrPermissionDenied.Error(), he.Message)
		}
	})

	t.Run("Header input validation", func(t *testing.T) {
		payload := `{"origin_lat": 0.0, "origin_lon": 0.0, "target_lat": 0.0, "target_lon": 0.0}`
		tests := []struct {
			wantErr bool
			header  string
			value   string
		}{
			{wantErr: true, header: gofakeit.AnimalType(), value: gofakeit.LetterN(uint(gofakeit.Number(101, 1000)))},
			{wantErr: true, header: gofakeit.AnimalType(), value: gofakeit.LetterN(uint(gofakeit.Number(0, 100)))},
			{wantErr: true, header: "idempotency-key", value: gofakeit.LetterN(uint(gofakeit.Number(101, 1000)))},
			{wantErr: false, header: "idempotency-key", value: gofakeit.LetterN(uint(gofakeit.Number(0, 100)))},
		}
		for _, tc := range tests {
			if !tc.wantErr {
				rCode := idempotency.ResponseCodeOK
				rBody := idempotency.ResponseBody{Message: "ok"}

				ik := &entity.IdempotencyKey{
					ResponseCode: &rCode,
					ResponseBody: &rBody,
				}
				uc.On("Create", callArgs...).Once().Return(ik, nil)
			}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			// testing different header values
			req.Header.Set(tc.header, tc.value)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set(string(entity.UserCtxKey), user)

			err := handler.Create(c)

			if tc.wantErr {
				if assert.Error(t, err) && assert.IsType(t, (&echo.HTTPError{}), err) {
					he := (err).(*echo.HTTPError)
					assert.Equal(t, http.StatusBadRequest, he.Code)
				}
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.Equal(t, "{\"message\":\"ok\"}", rec.Body.String())
				}
			}
		}
	})

	t.Run("Body input validation", func(t *testing.T) {
		tpl := `{"origin_lat": %v, "origin_lon": %v, "target_lat": %v, "target_lon": %v}`
		tests := []struct {
			wantErr bool
			payload string
		}{
			{wantErr: true, payload: fmt.Sprintf(tpl, -90.0000000001, 0.0, 0.0, 0.0)},
			{wantErr: true, payload: fmt.Sprintf(tpl, 90.0000000001, 0.0, 0.0, 0.0)},
			{wantErr: true, payload: fmt.Sprintf(tpl, 0.0, -180.0000000001, 0.0, 0.0)},
			{wantErr: true, payload: fmt.Sprintf(tpl, 0.0, 180.0000000001, 0.0, 0.0)},
			{wantErr: true, payload: fmt.Sprintf(tpl, 0.0, 0.0, -90.0000000001, 0.0)},
			{wantErr: true, payload: fmt.Sprintf(tpl, 0.0, 0.0, 90.0000000001, 0.0)},
			{wantErr: true, payload: fmt.Sprintf(tpl, 0.0, 0.0, 0.0, -180.0000000001)},
			{wantErr: true, payload: fmt.Sprintf(tpl, 0.0, 0.0, 0.0, 180.0000000001)},

			{wantErr: false, payload: fmt.Sprintf(tpl, -90.0, 0.0, 0.0, 0.0)},
			{wantErr: false, payload: fmt.Sprintf(tpl, 90.0, 0.0, 0.0, 0.0)},
			{wantErr: false, payload: fmt.Sprintf(tpl, 0.0, -180.0, 0.0, 0.0)},
			{wantErr: false, payload: fmt.Sprintf(tpl, 0.0, 180.0, 0.0, 0.0)},
			{wantErr: false, payload: fmt.Sprintf(tpl, 0.0, 0.0, -90.0, 0.0)},
			{wantErr: false, payload: fmt.Sprintf(tpl, 0.0, 0.0, 90.0, 0.0)},
			{wantErr: false, payload: fmt.Sprintf(tpl, 0.0, 0.0, 0.0, -180.0)},
			{wantErr: false, payload: fmt.Sprintf(tpl, 0.0, 0.0, 0.0, 180.0)},
		}
		for _, tc := range tests {
			if !tc.wantErr {
				rCode := idempotency.ResponseCodeOK
				rBody := idempotency.ResponseBody{Message: "ok"}

				ik := &entity.IdempotencyKey{
					ResponseCode: &rCode,
					ResponseBody: &rBody,
				}
				uc.On("Create", callArgs...).Once().Return(ik, nil)
			}

			// testing different payload values
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.payload))

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("idempotency-key", gofakeit.UUID())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set(string(entity.UserCtxKey), user)

			err := handler.Create(c)

			if tc.wantErr {
				if assert.Error(t, err) && assert.IsType(t, (&echo.HTTPError{}), err) {
					he := (err).(*echo.HTTPError)
					assert.Equal(t, http.StatusBadRequest, he.Code)
				}
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, http.StatusOK, rec.Code)
					assert.Equal(t, "{\"message\":\"ok\"}", rec.Body.String())
				}
			}
		}
	})

	t.Run("Error on create ride", func(t *testing.T) {
		payload := `{"origin_lat": 0.0, "origin_lon": 0.0, "target_lat": 0.0, "target_lon": 0.0}`
		rCode := idempotency.ResponseCodeOK
		rBody := idempotency.ResponseBody{Message: "filled"}
		tests := []struct {
			desc    string
			retArgs []interface{}
			retCode int
			retMsg  string
		}{
			{
				desc:    "error on idempotency key params mismatch",
				retArgs: []interface{}{nil, entity.ErrIdemKeyParamsMismatch},
				retCode: http.StatusConflict,
				retMsg:  entity.ErrIdemKeyParamsMismatch.Error(),
			},
			{
				desc:    "error on idempotency key request in progress",
				retArgs: []interface{}{nil, entity.ErrIdemKeyRequestInProgress},
				retCode: http.StatusConflict,
				retMsg:  entity.ErrIdemKeyRequestInProgress.Error(),
			},
			{
				desc:    "generic error on create ride uc",
				retArgs: []interface{}{nil, errors.New("error CreateRide")},
				retCode: http.StatusInternalServerError,
				retMsg:  "internal error",
			},
			{
				desc:    "empty idempotency key response code",
				retArgs: []interface{}{&entity.IdempotencyKey{ResponseBody: &rBody}, nil},
				retCode: http.StatusInternalServerError,
				retMsg:  "internal error",
			},
			{
				desc:    "empty idempotency key response body",
				retArgs: []interface{}{&entity.IdempotencyKey{ResponseCode: &rCode}, nil},
				retCode: http.StatusInternalServerError,
				retMsg:  "internal error",
			},
		}

		for _, tc := range tests {
			uc.On("Create", callArgs...).Once().Return(tc.retArgs...)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("idempotency-key", gofakeit.UUID())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set(string(entity.UserCtxKey), user)

			err := handler.Create(c)

			if assert.Error(t, err) && assert.IsType(t, (&echo.HTTPError{}), err) {
				he := (err).(*echo.HTTPError)
				assert.Equal(t, tc.retCode, he.Code, tc.desc)
				assert.Equal(t, tc.retMsg, he.Message, tc.desc)
			}
		}
	})

	t.Run("Success on create ride", func(t *testing.T) {
		payload := `{"origin_lat": 0.0, "origin_lon": 0.0, "target_lat": 0.0, "target_lon": 0.0}`
		rCode := idempotency.ResponseCodeOK
		rBody := idempotency.ResponseBody{Message: gofakeit.UUID()}
		body, err := json.Marshal(rBody)
		require.NoError(t, err)

		ik := &entity.IdempotencyKey{
			ResponseCode: &rCode,
			ResponseBody: &rBody,
		}
		uc.On("Create", callArgs...).Once().Return(ik, nil)

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("idempotency-key", gofakeit.UUID())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(string(entity.UserCtxKey), user)

		err = handler.Create(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, body, rec.Body.Bytes())
		}
	})
}
