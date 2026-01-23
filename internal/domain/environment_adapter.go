package domain

import "context"

// EnvironmentAdapter defines the contract for environment-specific trust configuration.
//
// Each supported environment (OpenSSL, Python, Node.js, Git, etc.) implements this interface.
// The Orchestrator treats all environments uniformly through this contract.
//
// Lifecycle:
//  1. Info() - Identify the environment
//  2. Detect() - Check if present on the system
//  3. Apply() - Install the certificate
//  4. Verify() - Confirm trust works
//
// Implementations must be:
//   - Idempotent: running Apply multiple times must be safe
//   - Isolated: failures don't affect other environments
//   - Self-contained: no dependencies on other adapters
type EnvironmentAdapter interface {
	// Info returns the identity of this environment.
	// This is a pure function with no side effects.
	Info() Environment

	// Detect checks whether this environment is present and configurable.
	//
	// Detection must be:
	//   - Non-destructive: no system changes
	//   - Fast: avoid expensive operations
	//   - Honest: don't guess, verify
	//
	// Returns:
	//   - true, nil: environment detected and ready
	//   - false, nil: environment not present (not an error)
	//   - false, error: detection failed due to an error
	Detect(ctx context.Context) (bool, error)

	// Apply installs the certificate into this environment.
	//
	// Apply must be:
	//   - Idempotent: safe to run multiple times
	//   - Atomic: don't leave partial state on failure
	//   - Reversible: (future) support rollback
	//
	// Returns:
	//   - StateApplied: certificate was installed
	//   - StateAlreadyTrusted: certificate was already present
	//   - StateFailed: installation failed (check error)
	Apply(ctx context.Context, cert Certificate) (State, error)

	// Verify confirms that TLS trust is working correctly.
	//
	// Verification should:
	//   - Test real TLS behavior (not just file presence)
	//   - Be quick (use local verification if possible)
	//   - Be specific to this environment
	//
	// Returns nil if verification succeeds, error otherwise.
	Verify(ctx context.Context) error
}
