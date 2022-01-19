-- name: GetRideByID :one
SELECT * FROM rides
WHERE id = $1 LIMIT 1;

-- name: GetRideByIdempotencyKeyID :one
SELECT * FROM rides
WHERE idempotency_key_id = $1 LIMIT 1;

-- name: CreateRide :one
INSERT INTO rides(
    created_at,
    idempotency_key_id,
    origin_lat,
    origin_lon,
    target_lat,
    target_lon,
    stripe_charge_id,
    user_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateRide :one
UPDATE rides SET
    stripe_charge_id=$2
WHERE id = $1
RETURNING *;