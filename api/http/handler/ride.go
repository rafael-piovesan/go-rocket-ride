package handler

import (
	"encoding/json"
	"errors"
	"math"

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

type RideHandler struct {
	*Handler
	uc rocketride.RideUseCase
}

func NewRideHandler(uc rocketride.RideUseCase) *RideHandler {
	return &RideHandler{
		Handler: New(),
		uc:      uc,
	}
}

func (r *RideHandler) Create(c echo.Context) error {
	user, err := r.GetUserFromCtx(c)
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
		IdempotencyKey: cr.IdemKey,
		RequestMethod:  c.Request().Method,
		RequestPath:    c.Request().RequestURI,
		RequestParams:  rp,
		UserID:         user.ID,
		User:           user,
	}

	ik, err = r.uc.Create(c.Request().Context(), ik, rd)
	if err != nil {
		return r.HandleError(err)
	}

	rCode, rBody, err := r.handleResponse(ik)
	if err != nil {
		return r.HandleError(err)
	}

	return c.JSONBlob(rCode, rBody)
}

func (r *RideHandler) handleResponse(ik *entity.IdempotencyKey) (rCode int, rBody []byte, err error) {
	if ik.ResponseCode == nil || ik.ResponseBody == nil {
		err = errors.New("invalid response")
		return
	}

	rCode = int(*ik.ResponseCode)
	rBody, err = json.Marshal(*ik.ResponseBody)
	return
}
