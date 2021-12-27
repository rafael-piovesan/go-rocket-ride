package entity

type UserCtxKeyType string

const (
	UserCtxKey UserCtxKeyType = "user-ctx-key"
)

type User struct {
	ID               int64
	Email            string
	StripeCustomerID string
}
