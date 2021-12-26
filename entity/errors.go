package entity

import "errors"

var (
	ErrPermissionDenied            = errors.New("permission denied")
	ErrNotFound                    = errors.New("entity not found")
	ErrIdemKeyParamsMismatch       = errors.New("params mismatch")
	ErrIdemKeyRequestInProgress    = errors.New("request in progress")
	ErrIdemKeyUnknownRecoveryPoint = errors.New("unknown recovery point")
	ErrPaymentProvider             = errors.New("card error from payment processor")
	ErrPaymentProviderGeneric      = errors.New("generic error from payment processor")
	ErrInternalError               = errors.New("internal error")
)
