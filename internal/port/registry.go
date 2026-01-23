package port

// Registry provides access to all known environment adapters.
//
// The Registry is a catalog - it declares which environments Trustica supports.
// It does NOT:
//   - Perform detection
//   - Apply certificates
//   - Execute any verification
//
// The Registry is static and explicit. New environments are added by
// registering their adapters here.
//
// The Core uses the Registry to get all adapters, then orchestrates them.
type Registry interface {
	// All returns all registered environment adapters.
	// The order is deterministic and reflects registration order.
	All() []EnvironmentAdapter

	// Get returns a specific adapter by environment ID.
	// Returns nil if the environment is not registered.
	Get(envID string) EnvironmentAdapter
}
