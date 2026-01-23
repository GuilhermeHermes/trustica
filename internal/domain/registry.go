package domain

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
// The Orchestrator uses the Registry to get all adapters, then orchestrates them.
//
// This interface is defined in the domain because the Orchestrator depends on it.
// Concrete implementations live in adapter/registry/.
type Registry interface {
	// All returns all registered environment adapters.
	// The order is deterministic and reflects registration order.
	All() []EnvironmentAdapter
}
