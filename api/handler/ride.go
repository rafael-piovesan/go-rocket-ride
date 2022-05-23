package handler

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/usecase"
)

type createRequest struct {
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

type Ride struct {
	Handler
	uc usecase.Ride
}

func NewRide(uc usecase.Ride) Ride {
	return Ride{
		Handler: New(),
		uc:      uc,
	}
}

func (r Ride) Create(c echo.Context) error {
	ik, err := r.IdempotencyKey(c)
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
	ik.RequestParams = rp

	err = r.uc.Create(c.Request().Context(), &ik, rd)
	if err != nil {
		return err
	}

	if ik.ResponseCode == nil || ik.ResponseBody == nil {
		return errors.New("create ride: invalid response")
	}

	return c.JSON(int(*ik.ResponseCode), ik.ResponseBody)
}
