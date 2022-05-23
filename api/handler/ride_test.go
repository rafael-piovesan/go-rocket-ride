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
	"github.com/rafael-piovesan/go-rocket-ride/v2/api/context"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity/idempotency"
	mocks "github.com/rafael-piovesan/go-rocket-ride/v2/mocks/usecase"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/httpserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateValidate(t *testing.T) {
	e := httpserver.New()
	uc := &mocks.Ride{}

	handler := NewRide(uc)

	callArgs := []interface{}{
		mock.Anything,
		mock.AnythingOfType("*entity.IdempotencyKey"),
		mock.AnythingOfType("*entity.Ride"),
	}

	tpl := `{"origin_lat": %v, "origin_lon": %v, "target_lat": %v, "target_lon": %v}`
	tests := []struct {
		fail    bool
		payload string
	}{
		{fail: true, payload: fmt.Sprintf(tpl, -90.0000000001, 0.0, 0.0, 0.0)},
		{fail: true, payload: fmt.Sprintf(tpl, 90.0000000001, 0.0, 0.0, 0.0)},
		{fail: true, payload: fmt.Sprintf(tpl, 0.0, -180.0000000001, 0.0, 0.0)},
		{fail: true, payload: fmt.Sprintf(tpl, 0.0, 180.0000000001, 0.0, 0.0)},
		{fail: true, payload: fmt.Sprintf(tpl, 0.0, 0.0, -90.0000000001, 0.0)},
		{fail: true, payload: fmt.Sprintf(tpl, 0.0, 0.0, 90.0000000001, 0.0)},
		{fail: true, payload: fmt.Sprintf(tpl, 0.0, 0.0, 0.0, -180.0000000001)},
		{fail: true, payload: fmt.Sprintf(tpl, 0.0, 0.0, 0.0, 180.0000000001)},

		{fail: false, payload: fmt.Sprintf(tpl, -90.0, 0.0, 0.0, 0.0)},
		{fail: false, payload: fmt.Sprintf(tpl, 90.0, 0.0, 0.0, 0.0)},
		{fail: false, payload: fmt.Sprintf(tpl, 0.0, -180.0, 0.0, 0.0)},
		{fail: false, payload: fmt.Sprintf(tpl, 0.0, 180.0, 0.0, 0.0)},
		{fail: false, payload: fmt.Sprintf(tpl, 0.0, 0.0, -90.0, 0.0)},
		{fail: false, payload: fmt.Sprintf(tpl, 0.0, 0.0, 90.0, 0.0)},
		{fail: false, payload: fmt.Sprintf(tpl, 0.0, 0.0, 0.0, -180.0)},
		{fail: false, payload: fmt.Sprintf(tpl, 0.0, 0.0, 0.0, 180.0)},
	}
	for _, tc := range tests {
		if !tc.fail {
			rCode := idempotency.ResponseCodeOK
			rBody := idempotency.ResponseBody{Message: "ok"}

			uc.On("Create", callArgs...).Once().Return(nil).Run(func(args mock.Arguments) {
				arg, ok := args.Get(1).(*entity.IdempotencyKey)
				assert.True(t, ok)
				arg.ResponseCode = &rCode
				arg.ResponseBody = &rBody
			})
		}

		// testing different payload values
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.payload))
		rec := httptest.NewRecorder()

		req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

		c := e.NewContext(req, rec)
		context.AddUser(c, entity.User{})
		context.AddIdemKey(c, entity.IdempotencyKey{})

		err := handler.Create(c)

		if tc.fail {
			assert.Error(t, err)

			var he *echo.HTTPError
			if errors.As(err, &he) {
				assert.Equal(t, http.StatusBadRequest, he.Code)
			}
		} else {
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "{\"message\":\"ok\"}\n", rec.Body.String())
		}
	}
}

func TestCreate(t *testing.T) {
	e := echo.New()
	uc := &mocks.Ride{}

	handler := NewRide(uc)

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

		assert.Error(t, err)

		var he *echo.HTTPError
		if errors.As(err, &he) {
			assert.Equal(t, http.StatusUnauthorized, he.Code)
			assert.Equal(t, entity.ErrPermissionDenied.Error(), he.Message)
		}
	})

	t.Run("Error on create ride", func(t *testing.T) {
		payload := `{"origin_lat": 0.0, "origin_lon": 0.0, "target_lat": 0.0, "target_lon": 0.0}`
		rCode := idempotency.ResponseCodeOK
		rBody := idempotency.ResponseBody{Message: "filled"}
		tests := []struct {
			desc     string
			retError error
			retCode  int
			retMsg   string
			retFunc  func(mock.Arguments)
		}{
			{
				desc:     "error on idempotency key params mismatch",
				retError: entity.ErrIdemKeyParamsMismatch,
				retCode:  http.StatusConflict,
				retMsg:   entity.ErrIdemKeyParamsMismatch.Error(),
				retFunc:  func(mock.Arguments) {},
			},
			{
				desc:     "error on idempotency key request in progress",
				retError: entity.ErrIdemKeyRequestInProgress,
				retCode:  http.StatusConflict,
				retMsg:   entity.ErrIdemKeyRequestInProgress.Error(),
				retFunc:  func(mock.Arguments) {},
			},
			{
				desc:     "generic error on create ride uc",
				retError: errors.New("error CreateRide"),
				retCode:  http.StatusInternalServerError,
				retMsg:   "internal error",
				retFunc:  func(mock.Arguments) {},
			},
			{
				desc:     "empty idempotency key response code",
				retError: errors.New("create ride: invalid response"),
				retCode:  http.StatusInternalServerError,
				retMsg:   "internal error",
				retFunc: func(args mock.Arguments) {
					arg, ok := args.Get(1).(*entity.IdempotencyKey)
					assert.True(t, ok)
					arg.ResponseBody = &rBody
				},
			},
			{
				desc:     "empty idempotency key response body",
				retError: errors.New("create ride: invalid response"),
				retCode:  http.StatusInternalServerError,
				retMsg:   "internal error",
				retFunc: func(args mock.Arguments) {
					arg, ok := args.Get(1).(*entity.IdempotencyKey)
					assert.True(t, ok)
					arg.ResponseCode = &rCode
				},
			},
		}

		for _, tc := range tests {
			uc.On("Create", callArgs...).Once().Return(tc.retError).Run(tc.retFunc)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
			rec := httptest.NewRecorder()

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			c := e.NewContext(req, rec)
			context.AddUser(c, entity.User{})
			context.AddIdemKey(c, entity.IdempotencyKey{})

			err := handler.Create(c)

			assert.Equal(t, tc.retError, err)
		}
	})

	t.Run("Success on create ride", func(t *testing.T) {
		payload := `{"origin_lat": 0.0, "origin_lon": 0.0, "target_lat": 0.0, "target_lon": 0.0}`
		rCode := idempotency.ResponseCodeOK
		rBody := idempotency.ResponseBody{Message: gofakeit.UUID()}
		body, err := json.Marshal(rBody)
		require.NoError(t, err)

		uc.On("Create", callArgs...).Once().Return(nil).Run(func(args mock.Arguments) {
			arg, ok := args.Get(1).(*entity.IdempotencyKey)
			assert.True(t, ok)
			arg.ResponseCode = &rCode
			arg.ResponseBody = &rBody
		})

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
		rec := httptest.NewRecorder()

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		c := e.NewContext(req, rec)
		context.AddUser(c, entity.User{})
		context.AddIdemKey(c, entity.IdempotencyKey{})

		err = handler.Create(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), string(body))
		}
	})
}
