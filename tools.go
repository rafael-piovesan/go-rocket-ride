//go:build tools
// +build tools

// Package tool helps keeping track of dev dependencies.
package tools

import (
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/kyleconroy/sqlc/cmd/sqlc"
	_ "github.com/lib/pq"
	_ "github.com/segmentio/golines"
)
