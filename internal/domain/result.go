package domain

import "fmt"

// Result represents the outcome of processing a single environment.
//
// Each environment produces exactly one Result after being processed
// by the Orchestrator. Results are aggregated and returned to the CLI
// for display to the user.
//
// Result is a value object - immutable after creation.
type Result struct {
	// Environment identifies which environment this result is for.
	Environment Environment

	// State is the final state after processing.
	State State

	// Message provides optional human-readable context.
	// Examples: "Certificate added to /etc/ssl/certs/ca-certificates.crt"
	//           "certifi bundle updated"
	//           "git config http.sslCAInfo set"
	Message string

	// Error contains the error if State is StateFailed.
	// Will be nil for successful states.
	Error error
}

// NewResult creates a successful result with the given state and message.
func NewResult(env Environment, state State, message string) Result {
	return Result{
		Environment: env,
		State:       state,
		Message:     message,
		Error:       nil,
	}
}

// NewErrorResult creates a failed result with the given error.
func NewErrorResult(env Environment, err error) Result {
	return Result{
		Environment: env,
		State:       StateFailed,
		Message:     err.Error(),
		Error:       err,
	}
}

// NewSkippedResult creates a skipped result with an optional reason.
func NewSkippedResult(env Environment, reason string) Result {
	return Result{
		Environment: env,
		State:       StateSkipped,
		Message:     reason,
		Error:       nil,
	}
}

// NewNotDetectedResult creates a not-detected result.
func NewNotDetectedResult(env Environment) Result {
	return Result{
		Environment: env,
		State:       StateNotDetected,
		Message:     "Environment not detected",
		Error:       nil,
	}
}

// IsSuccess returns true if this result represents a successful outcome.
func (r Result) IsSuccess() bool {
	return r.State.IsSuccess()
}

// IsFailed returns true if this result represents a failure.
func (r Result) IsFailed() bool {
	return r.State == StateFailed
}

// String returns a human-readable representation of the result.
func (r Result) String() string {
	if r.Message != "" {
		return fmt.Sprintf("%s: %s (%s)", r.Environment.Name, r.State, r.Message)
	}
	return fmt.Sprintf("%s: %s", r.Environment.Name, r.State)
}
