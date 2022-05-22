// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: user.sql

package sqlc

import (
	"context"
)

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, stripe_customer_id FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.StripeCustomerID)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, stripe_customer_id FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.StripeCustomerID)
	return i, err
}
