package domain

import "errors"

// Domain errors.
//
// These errors represent domain-level failures that can occur during
// trust configuration. They are distinct from adapter-specific errors.
var (
	// ErrAdapterNotFound is returned when an adapter cannot be found by environment ID.
	// This typically indicates an internal error in the orchestration logic.
	ErrAdapterNotFound = errors.New("adapter not found for environment")
)
