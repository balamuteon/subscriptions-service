//go:build !integration
// +build !integration

package testdb

import (
	"context"
	"errors"
)

// SetupTestDatabase is a stub for non-integration builds.
func SetupTestDatabase(_ context.Context) (string, func(), error) {
	return "", func() {}, errors.New("integration build tag required")
}
