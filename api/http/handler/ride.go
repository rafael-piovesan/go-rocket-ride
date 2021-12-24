package handler

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	rocketride "github.com/rafael-piovesan/go-rocket-ride"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
)

type createRequest struct {
	IdemKey string  `json:"idem_key" header:"idempotency-key" validate:"required,max=100"`
	OrigLat float64 `json:"origin_lat" validate:"min=-90,max=90"`
	OrigLon float64 `json:"origin_lon" validate:"min=-180,max=180"`
	TgtLat  float64 `json:"target_lat" validate:"min=-90,max=90"`
	TgtLon  float64 `json:"target_lon" validate:"min=-180,max=180"`
}

func newCreateRequest() createRequest {
	return createRequest{
		OrigLat: math.Inf(-1),
		OrigLon: math.Inf(-1),
		TgtLat:  math.Inf(-1),
		TgtLon:  math.Inf(-1),
	}
}

type RideHanlder struct {
	*Handler
	uc rocketride.RideUseCase
}

func NewRideHandler(uc rocketride.RideUseCase) *RideHanlder {
	return &RideHanlder{
		Handler: New(),
		uc:      uc,
	}
}

func (r *RideHanlder) Create(c echo.Context) error {
	userID, err := r.GetUserID(c)
	if err != nil {
		return err
	}

	cr := newCreateRequest()
	if err := r.BindAndValidate(c, &cr); err != nil {
		return err
	}

	rd := &entity.Ride{
		OriginLat: cr.OrigLat,
		OriginLon: cr.OrigLon,
		TargetLat: cr.TgtLat,
		TargetLon: cr.TgtLon,
	}

	rp, _ := json.Marshal(rd)

	ik := &entity.IdempotencyKey{
		CreatedAt:      time.Now().UTC(),
		IdempotencyKey: cr.IdemKey,
		UserID:         userID,
		RequestMethod:  c.Request().Method,
		RequestPath:    c.Request().RequestURI,
		RequestParams:  rp,
	}

	ik, err = r.uc.Create(c.Request().Context(), ik, rd)

	if errors.Is(err, entity.ErrIdemKeyParamsMismatch) || errors.Is(err, entity.ErrIdemKeyRequestInProgress) {
		return c.String(http.StatusConflict, err.Error())
	}

	if err != nil || ik.ResponseCode == nil || ik.ResponseBody == nil {
		return c.String(http.StatusInternalServerError, "internal error")
	}

	rCode := int(*ik.ResponseCode)
	rBody, err := json.Marshal(*ik.ResponseBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal error")
	}

	return c.JSONBlob(rCode, rBody)
}
