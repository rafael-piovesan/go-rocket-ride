package entity

import "time"

type Ride struct {
	ID               int64
	CreatedAt        time.Time
	IdempotencyKeyID *int64
	OriginLat        string
	OriginLon        string
	TargetLat        string
	TargetLon        string
	StripeChargeID   *string
	UserID           int64
}
