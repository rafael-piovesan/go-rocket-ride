package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
)

type Handler struct {
	binder   *echo.DefaultBinder
	validate *validator.Validate
}

func New() *Handler {
	return &Handler{
		binder:   &echo.DefaultBinder{},
		validate: validator.New(),
	}
}

// BindAndValidate binds and validates the data struct pointed to by 'i'.
// It expects a stuct pointer as parameter, since it needs to populate its
// fields, otherwise it'll panic.
func (h *Handler) BindAndValidate(c echo.Context, i interface{}) error {
	if i == nil {
		return nil
	}

	if err := h.binder.BindHeaders(c, i); err != nil {
		return err
	}

	if err := h.binder.BindBody(c, i); err != nil {
		return err
	}

	if err := h.validate.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (h *Handler) GetUserFromCtx(c echo.Context) (uid entity.User, err error) {
	uid, ok := c.Get(string(entity.UserCtxKey)).(entity.User)
	if !ok {
		err = echo.NewHTTPError(http.StatusUnauthorized, entity.ErrPermissionDenied.Error())
	}
	return
}

func (h *Handler) HandleError(err error) error {
	switch {
	case errors.Is(err, entity.ErrIdemKeyParamsMismatch) || errors.Is(err, entity.ErrIdemKeyRequestInProgress):
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	case errors.Is(err, entity.ErrPaymentProvider):
		return echo.NewHTTPError(http.StatusPaymentRequired, err.Error())
	case errors.Is(err, entity.ErrPaymentProviderGeneric):
		return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, entity.ErrInternalError.Error())
	}
}
