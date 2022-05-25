package entity

import "time"

type Ride struct {
	ID               int64
	CreatedAt        time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	IdempotencyKeyID *int64
	OriginLat        float64 `json:"origin_lat"`
	OriginLon        float64 `json:"origin_lon"`
	TargetLat        float64 `json:"target_lat"`
	TargetLon        float64 `json:"target_lon"`
	StripeChargeID   *string
	UserID           int64
}
