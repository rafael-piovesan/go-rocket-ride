package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/v2/api/context"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
)

type idemKeyRequest struct {
	IdemKey string `header:"idempotency-key" validate:"required,max=100"`
}

func IdempotencyKey() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			binder := &echo.DefaultBinder{}
			validate := validator.New()
			ikr := idemKeyRequest{}

			if err := binder.BindHeaders(c, &ikr); err != nil {
				return err
			}

			if err := validate.Struct(&ikr); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			rawBody, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return err
			}

			// Restore the io.ReadCloser to it's original state
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))

			ik := entity.IdempotencyKey{
				IdempotencyKey: ikr.IdemKey,
				RequestMethod:  c.Request().Method,
				RequestPath:    c.Request().RequestURI,
				RequestParams:  rawBody,
			}

			context.AddIdemKey(c, ik)

			return next(c)
		}
	}
}
