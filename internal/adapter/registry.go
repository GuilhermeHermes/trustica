package adapter

import (
	"github.com/GuilhermeHermes/trustica/internal/domain"
	"github.com/GuilhermeHermes/trustica/internal/port"
)

// InMemoryRegistry is a simple in-memory implementation of domain.Registry.
// Adapters are stored in registration order.
type InMemoryRegistry struct {
	adapters []domain.EnvironmentAdapter
}

// NewRegistry creates an empty registry.
func NewRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		adapters: make([]domain.EnvironmentAdapter, 0),
	}
}

// Register adds an adapter to the registry.
// Adapters are stored in the order they are registered.
func (r *InMemoryRegistry) Register(adapter domain.EnvironmentAdapter) {
	r.adapters = append(r.adapters, adapter)
}

// All returns all registered adapters in registration order.
func (r *InMemoryRegistry) All() []domain.EnvironmentAdapter {
	return r.adapters
}

// DefaultRegistry creates a registry with all built-in adapters.
// This is the standard registry used by the CLI.
func DefaultRegistry(sys port.System) *InMemoryRegistry {
	r := NewRegistry()

	// TODO: Register built-in adapters here as we implement them
	// r.Register(git.NewAdapter(sys))
	// r.Register(python.NewAdapter(sys))
	// r.Register(nodejs.NewAdapter(sys))
	// r.Register(openssl.NewAdapter(sys))

	_ = sys // Will be used when we add adapters

	return r
}
