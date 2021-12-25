package entity

import "time"

type Ride struct {
	ID               int64
	CreatedAt        time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	IdempotencyKeyID *int64
	OriginLat        float64
	OriginLon        float64
	TargetLat        float64
	TargetLon        float64
	StripeChargeID   *string
	UserID           int64
}
