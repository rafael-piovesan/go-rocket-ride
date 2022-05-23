package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity/originip"
)

func OriginIP() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			oip := &originip.OriginIP{
				IP: c.RealIP(),
			}
			ctxIP := originip.NewContext(c.Request().Context(), oip)
			c.SetRequest(c.Request().WithContext(ctxIP))
			return next(c)
		}
	}
}
