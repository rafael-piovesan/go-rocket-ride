package entity

import "errors"

var (
	ErrNotFound                    = errors.New("entity not found")
	ErrIdemKeyParamsMismatch       = errors.New("params mismatch")
	ErrIdemKeyRequestInProgress    = errors.New("request in progress")
	ErrIdemKeyUnknownRecoveryPoint = errors.New("unknown recovery point")
)
