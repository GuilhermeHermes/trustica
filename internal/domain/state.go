package domain

// State represents the lifecycle state of an environment during trust configuration.
//
// State transitions follow this flow:
//
//	                ┌──────────────┐
//	                │ NotDetected  │ ──────────────────┐
//	                └──────────────┘                   │
//	                       │                           │
//	                 (if detected)                     │
//	                       ▼                           ▼
//	                ┌──────────────┐            ┌──────────┐
//	                │   Detected   │            │  Skipped │
//	                └──────────────┘            └──────────┘
//	                       │
//	        ┌──────────────┼──────────────┐
//	        │              │              │
//	   (apply ok)    (already ok)    (error)
//	        ▼              ▼              ▼
//	 ┌──────────┐  ┌──────────────┐  ┌────────┐
//	 │ Applied  │  │AlreadyTrusted│  │ Failed │
//	 └──────────┘  └──────────────┘  └────────┘
//	        │              │
//	        └──────┬───────┘
//	               ▼
//	          (verify)
//	               │
//	      ┌────────┴────────┐
//	      ▼                 ▼
//	┌──────────┐      ┌────────┐
//	│ Verified │      │ Failed │
//	└──────────┘      └────────┘
//
// State transitions are controlled exclusively by the Orchestrator (Core).
// Adapters report outcomes; the Core decides state changes.
type State int

const (
	// StateNotDetected means the environment is not present on this system.
	// This is the initial state when detection fails to find the environment.
	StateNotDetected State = iota

	// StateDetected means the environment exists and can be configured.
	// The environment is ready for the Apply phase.
	StateDetected

	// StateApplied means the certificate was installed successfully.
	// The environment is ready for verification.
	StateApplied

	// StateAlreadyTrusted means the certificate was already present.
	// No changes were made; idempotency preserved.
	StateAlreadyTrusted

	// StateVerified means TLS trust was confirmed working.
	// This is a successful terminal state.
	StateVerified

	// StateFailed means an error occurred during any phase.
	// Check the Result.Error for details.
	StateFailed

	// StateSkipped means the environment was intentionally ignored.
	// This can happen due to user flags or configuration.
	StateSkipped
)

// String returns a human-readable representation of the state.
func (s State) String() string {
	switch s {
	case StateNotDetected:
		return "NotDetected"
	case StateDetected:
		return "Detected"
	case StateApplied:
		return "Applied"
	case StateAlreadyTrusted:
		return "AlreadyTrusted"
	case StateVerified:
		return "Verified"
	case StateFailed:
		return "Failed"
	case StateSkipped:
		return "Skipped"
	default:
		return "Unknown"
	}
}

// IsTerminal returns true if this state represents a final outcome.
func (s State) IsTerminal() bool {
	switch s {
	case StateVerified, StateFailed, StateSkipped, StateNotDetected, StateAlreadyTrusted:
		return true
	default:
		return false
	}
}

// IsSuccess returns true if this state represents a successful outcome.
func (s State) IsSuccess() bool {
	switch s {
	case StateVerified, StateAlreadyTrusted:
		return true
	default:
		return false
	}
}
