package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
)

func ErrorMapper() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			switch {
			case errors.Is(err, entity.ErrIdemKeyParamsMismatch) || errors.Is(err, entity.ErrIdemKeyRequestInProgress):
				err = echo.NewHTTPError(http.StatusConflict, err.Error())
			case errors.Is(err, entity.ErrPaymentProvider):
				err = echo.NewHTTPError(http.StatusPaymentRequired, err.Error())
			case errors.Is(err, entity.ErrPaymentProviderGeneric):
				err = echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
			}

			return err
		}
	}
}
