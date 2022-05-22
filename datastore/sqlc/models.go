// Code generated by sqlc. DO NOT EDIT.

package sqlc

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/tabbed/pqtype"
)

type AuditRecord struct {
	ID           int64
	Action       string
	CreatedAt    time.Time
	Data         json.RawMessage
	OriginIp     string
	ResourceID   int64
	ResourceType string
	UserID       int64
}

type IdempotencyKey struct {
	ID             int64
	CreatedAt      time.Time
	IdempotencyKey string
	LastRunAt      time.Time
	LockedAt       sql.NullTime
	RequestMethod  string
	RequestParams  json.RawMessage
	RequestPath    string
	ResponseCode   sql.NullInt32
	ResponseBody   pqtype.NullRawMessage
	RecoveryPoint  string
	UserID         int64
}

type Ride struct {
	ID               int64
	CreatedAt        time.Time
	IdempotencyKeyID sql.NullInt64
	OriginLat        float64
	OriginLon        float64
	TargetLat        float64
	TargetLon        float64
	StripeChargeID   sql.NullString
	UserID           int64
}

type StagedJob struct {
	ID      int64
	JobName string
	JobArgs json.RawMessage
}

type User struct {
	ID               int64
	Email            string
	StripeCustomerID string
}