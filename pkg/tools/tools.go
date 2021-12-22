//go:build tools
// +build tools

package tools

import (
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/segmentio/golines"
	_ "github.com/vektra/mockery/cmd/mockery"
)
