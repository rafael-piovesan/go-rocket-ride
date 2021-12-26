package idempotency

import (
	"encoding/json"
	"net/http"
)

type RecoveryPoint string

const (
	RecoveryPointStarted  RecoveryPoint = "STARTED"
	RecoveryPointCreated  RecoveryPoint = "CREATED"
	RecoveryPointCharged  RecoveryPoint = "CHARGED"
	RecoveryPointFinished RecoveryPoint = "FINISHED"
)

func (r RecoveryPoint) String() string {
	return string(r)
}

type ResponseCode int32

const (
	ResponseCodeOK                ResponseCode = http.StatusOK
	ResponseCodeConflict          ResponseCode = http.StatusConflict
	ResponseCodeErrPayment        ResponseCode = http.StatusPaymentRequired
	ResponseCodeErrPaymentGeneric ResponseCode = http.StatusServiceUnavailable
)

type ResponseBody struct {
	Message string `json:"message"`
}

func (r *ResponseBody) Marshal() (json.RawMessage, error) {
	return json.Marshal(r)
}
