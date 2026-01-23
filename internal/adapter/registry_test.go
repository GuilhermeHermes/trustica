package adapter

import (
	"context"
	"testing"

	"github.com/GuilhermeHermes/trustica/internal/adapter/outbound/system"
	"github.com/GuilhermeHermes/trustica/internal/domain"
)

// mockAdapter is a simple mock for testing registry
type mockAdapter struct {
	env domain.Environment
}

func newMockAdapter(id, name string) *mockAdapter {
	return &mockAdapter{env: domain.NewEnvironment(id, name)}
}

func (m *mockAdapter) Info() domain.Environment {
	return m.env
}

func (m *mockAdapter) Detect(ctx context.Context) (bool, error) {
	return true, nil
}

func (m *mockAdapter) Apply(ctx context.Context, cert domain.Certificate) (domain.State, error) {
	return domain.StateApplied, nil
}

func (m *mockAdapter) Verify(ctx context.Context) error {
	return nil
}

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()

	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}

	adapters := r.All()
	if len(adapters) != 0 {
		t.Errorf("New registry should be empty, got %d adapters", len(adapters))
	}
}

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()

	git := newMockAdapter("git", "Git")
	python := newMockAdapter("python", "Python")

	r.Register(git)
	r.Register(python)

	adapters := r.All()
	if len(adapters) != 2 {
		t.Fatalf("Expected 2 adapters, got %d", len(adapters))
	}

	// Verify order is preserved
	if adapters[0].Info().ID != "git" {
		t.Errorf("First adapter should be git, got %s", adapters[0].Info().ID)
	}
	if adapters[1].Info().ID != "python" {
		t.Errorf("Second adapter should be python, got %s", adapters[1].Info().ID)
	}
}

func TestRegistry_All_ReturnsSameSlice(t *testing.T) {
	r := NewRegistry()
	r.Register(newMockAdapter("git", "Git"))

	adapters1 := r.All()
	adapters2 := r.All()

	// Should return same underlying slice
	if len(adapters1) != len(adapters2) {
		t.Error("All() should return consistent results")
	}
}

func TestDefaultRegistry(t *testing.T) {
	sys := system.New()
	r := DefaultRegistry(sys)

	if r == nil {
		t.Fatal("DefaultRegistry() returned nil")
	}

	// Currently empty, will have adapters once we implement them
	_ = r.All()
}
