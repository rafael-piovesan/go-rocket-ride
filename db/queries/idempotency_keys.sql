-- name: GetIdempotencyKey :one
SELECT * FROM idempotency_keys
WHERE user_id = $1 AND idempotency_key = $2 LIMIT 1;

-- name: CreateIdempotencyKey :one
INSERT INTO idempotency_keys(
    created_at,
    idempotency_key,
    last_run_at,
    locked_at,
    request_method,
    request_params,
    request_path,
    response_code,
    response_body,
    recovery_point,
    user_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: UpdateIdempotencyKey :one
UPDATE idempotency_keys SET
    last_run_at=$2,
    locked_at=$3,
    response_code=$4,
    response_body=$5,
    recovery_point=$6
WHERE id = $1
RETURNING *;