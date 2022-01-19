-- name: CreateStagedJob :one
INSERT INTO staged_jobs(
	job_name,
	job_args
)
VALUES ($1, $2)
RETURNING *;