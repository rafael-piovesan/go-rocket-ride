-- name: GetRideByID :one
SELECT * FROM rides
WHERE id = $1 LIMIT 1;

-- name: GetRideByIdempotencyKeyID :one
SELECT * FROM rides
WHERE idempotency_key_id = $1 LIMIT 1;

-- name: UpdateRide :one
UPDATE rides SET
    stripe_charge_id=$2
WHERE id = $1
RETURNING *;