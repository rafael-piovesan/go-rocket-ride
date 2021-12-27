package entity

import (
	"encoding/json"
	"time"

	"github.com/rafael-piovesan/go-rocket-ride/entity/idempotency"
)

type IdempotencyKey struct {
	ID             int64
	CreatedAt      time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	IdempotencyKey string
	LastRunAt      time.Time
	LockedAt       *time.Time
	RequestMethod  string
	RequestParams  json.RawMessage
	RequestPath    string
	ResponseCode   *idempotency.ResponseCode
	ResponseBody   *idempotency.ResponseBody
	RecoveryPoint  idempotency.RecoveryPoint
	UserID         int64
	User           *User `json:"-" bun:"-"`
}
