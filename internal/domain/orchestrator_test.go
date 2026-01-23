package domain

import (
	"context"
	"errors"
	"testing"
)

// =============================================================================
// Mock Implementations
// =============================================================================

// MockAdapter is a configurable mock for EnvironmentAdapter
type MockAdapter struct {
	env          Environment
	detectResult bool
	detectError  error
	applyState   State
	applyError   error
	verifyError  error
	detectCalled bool
	applyCalled  bool
	verifyCalled bool
	appliedCert  Certificate
}

func NewMockAdapter(id, name string) *MockAdapter {
	return &MockAdapter{
		env:          NewEnvironment(id, name),
		detectResult: true,
		applyState:   StateApplied,
	}
}

func (m *MockAdapter) Info() Environment {
	return m.env
}

func (m *MockAdapter) Detect(ctx context.Context) (bool, error) {
	m.detectCalled = true
	return m.detectResult, m.detectError
}

func (m *MockAdapter) Apply(ctx context.Context, cert Certificate) (State, error) {
	m.applyCalled = true
	m.appliedCert = cert
	return m.applyState, m.applyError
}

func (m *MockAdapter) Verify(ctx context.Context) error {
	m.verifyCalled = true
	return m.verifyError
}

// WithDetect configures detection behavior
func (m *MockAdapter) WithDetect(result bool, err error) *MockAdapter {
	m.detectResult = result
	m.detectError = err
	return m
}

// WithApply configures apply behavior
func (m *MockAdapter) WithApply(state State, err error) *MockAdapter {
	m.applyState = state
	m.applyError = err
	return m
}

// WithVerify configures verify behavior
func (m *MockAdapter) WithVerify(err error) *MockAdapter {
	m.verifyError = err
	return m
}

// MockRegistry is a mock for Registry
type MockRegistry struct {
	adapters []EnvironmentAdapter
}

func NewMockRegistry(adapters ...EnvironmentAdapter) *MockRegistry {
	return &MockRegistry{adapters: adapters}
}

func (r *MockRegistry) All() []EnvironmentAdapter {
	return r.adapters
}

// =============================================================================
// Orchestrator Tests
// =============================================================================

func TestNewOrchestrator(t *testing.T) {
	registry := NewMockRegistry()
	orch := NewOrchestrator(registry)

	if orch == nil {
		t.Fatal("NewOrchestrator() returned nil")
	}
	if orch.registry != registry {
		t.Error("Orchestrator.registry not set correctly")
	}
}

func TestOrchestrator_Run_AllDetected(t *testing.T) {
	// Setup: 2 adapters, both detected, both apply successfully
	git := NewMockAdapter("git", "Git")
	python := NewMockAdapter("python", "Python")
	registry := NewMockRegistry(git, python)

	orch := NewOrchestrator(registry)
	cert := Certificate{Path: "/path/to/ca.pem", Content: []byte("test")}

	results := orch.Run(context.Background(), cert)

	// Should have 2 results
	if len(results) != 2 {
		t.Fatalf("Run() returned %d results, want 2", len(results))
	}

	// Both should be verified
	for _, r := range results {
		if r.State != StateVerified {
			t.Errorf("Result for %s: State = %v, want Verified", r.Environment.ID, r.State)
		}
	}

	// All phases should be called
	if !git.detectCalled || !git.applyCalled || !git.verifyCalled {
		t.Error("Git adapter: not all phases called")
	}
	if !python.detectCalled || !python.applyCalled || !python.verifyCalled {
		t.Error("Python adapter: not all phases called")
	}
}

func TestOrchestrator_Run_NoneDetected(t *testing.T) {
	// Setup: 2 adapters, neither detected
	git := NewMockAdapter("git", "Git").WithDetect(false, nil)
	python := NewMockAdapter("python", "Python").WithDetect(false, nil)
	registry := NewMockRegistry(git, python)

	orch := NewOrchestrator(registry)
	cert := Certificate{Path: "/path/to/ca.pem", Content: []byte("test")}

	results := orch.Run(context.Background(), cert)

	// Should have 2 results (both not detected)
	if len(results) != 2 {
		t.Fatalf("Run() returned %d results, want 2", len(results))
	}

	for _, r := range results {
		if r.State != StateNotDetected {
			t.Errorf("Result for %s: State = %v, want NotDetected", r.Environment.ID, r.State)
		}
	}

	// Apply and Verify should NOT be called
	if git.applyCalled || git.verifyCalled {
		t.Error("Git: Apply/Verify called for non-detected adapter")
	}
}

func TestOrchestrator_Run_PartialDetection(t *testing.T) {
	// Setup: git detected, python not detected
	git := NewMockAdapter("git", "Git").WithDetect(true, nil)
	python := NewMockAdapter("python", "Python").WithDetect(false, nil)
	registry := NewMockRegistry(git, python)

	orch := NewOrchestrator(registry)
	cert := Certificate{Path: "/path/to/ca.pem", Content: []byte("test")}

	results := orch.Run(context.Background(), cert)

	// Should have 2 results
	if len(results) != 2 {
		t.Fatalf("Run() returned %d results, want 2", len(results))
	}

	// Find results by environment
	var gitResult, pythonResult Result
	for _, r := range results {
		if r.Environment.ID == "git" {
			gitResult = r
		} else if r.Environment.ID == "python" {
			pythonResult = r
		}
	}

	if gitResult.State != StateVerified {
		t.Errorf("Git: State = %v, want Verified", gitResult.State)
	}
	if pythonResult.State != StateNotDetected {
		t.Errorf("Python: State = %v, want NotDetected", pythonResult.State)
	}
}

func TestOrchestrator_Run_DetectError(t *testing.T) {
	// Setup: git detection fails with error
	detectErr := errors.New("cannot run git")
	git := NewMockAdapter("git", "Git").WithDetect(false, detectErr)
	registry := NewMockRegistry(git)

	orch := NewOrchestrator(registry)
	cert := Certificate{Path: "/path/to/ca.pem", Content: []byte("test")}

	results := orch.Run(context.Background(), cert)

	if len(results) != 1 {
		t.Fatalf("Run() returned %d results, want 1", len(results))
	}

	if results[0].State != StateFailed {
		t.Errorf("State = %v, want Failed", results[0].State)
	}
	if results[0].Error != detectErr {
		t.Errorf("Error = %v, want %v", results[0].Error, detectErr)
	}
}

func TestOrchestrator_Run_ApplyError(t *testing.T) {
	// Setup: detection succeeds, apply fails
	applyErr := errors.New("permission denied")
	git := NewMockAdapter("git", "Git").WithApply(StateFailed, applyErr)
	registry := NewMockRegistry(git)

	orch := NewOrchestrator(registry)
	cert := Certificate{Path: "/path/to/ca.pem", Content: []byte("test")}

	results := orch.Run(context.Background(), cert)

	if len(results) != 1 {
		t.Fatalf("Run() returned %d results, want 1", len(results))
	}

	if results[0].State != StateFailed {
		t.Errorf("State = %v, want Failed", results[0].State)
	}

	// Verify should NOT be called on failed apply
	if git.verifyCalled {
		t.Error("Verify called after failed Apply")
	}
}

func TestOrchestrator_Run_VerifyError(t *testing.T) {
	// Setup: detection and apply succeed, verify fails
	verifyErr := errors.New("TLS handshake failed")
	git := NewMockAdapter("git", "Git").WithVerify(verifyErr)
	registry := NewMockRegistry(git)

	orch := NewOrchestrator(registry)
	cert := Certificate{Path: "/path/to/ca.pem", Content: []byte("test")}

	results := orch.Run(context.Background(), cert)

	if len(results) != 1 {
		t.Fatalf("Run() returned %d results, want 1", len(results))
	}

	if results[0].State != StateFailed {
		t.Errorf("State = %v, want Failed", results[0].State)
	}
	if results[0].Error != verifyErr {
		t.Errorf("Error = %v, want %v", results[0].Error, verifyErr)
	}
}

func TestOrchestrator_Run_AlreadyTrusted(t *testing.T) {
	// Setup: apply returns AlreadyTrusted
	git := NewMockAdapter("git", "Git").WithApply(StateAlreadyTrusted, nil)
	registry := NewMockRegistry(git)

	orch := NewOrchestrator(registry)
	cert := Certificate{Path: "/path/to/ca.pem", Content: []byte("test")}

	results := orch.Run(context.Background(), cert)

	if len(results) != 1 {
		t.Fatalf("Run() returned %d results, want 1", len(results))
	}

	// AlreadyTrusted should still be verified
	if results[0].State != StateVerified {
		t.Errorf("State = %v, want Verified", results[0].State)
	}
	if !git.verifyCalled {
		t.Error("Verify not called for AlreadyTrusted")
	}
}

func TestOrchestrator_Run_Isolation(t *testing.T) {
	// Setup: git fails, python succeeds
	// This tests that one failure doesn't affect others
	gitErr := errors.New("git error")
	git := NewMockAdapter("git", "Git").WithApply(StateFailed, gitErr)
	python := NewMockAdapter("python", "Python")
	registry := NewMockRegistry(git, python)

	orch := NewOrchestrator(registry)
	cert := Certificate{Path: "/path/to/ca.pem", Content: []byte("test")}

	results := orch.Run(context.Background(), cert)

	if len(results) != 2 {
		t.Fatalf("Run() returned %d results, want 2", len(results))
	}

	// Find results by environment
	var gitResult, pythonResult Result
	for _, r := range results {
		if r.Environment.ID == "git" {
			gitResult = r
		} else if r.Environment.ID == "python" {
			pythonResult = r
		}
	}

	// Git should fail, Python should succeed
	if gitResult.State != StateFailed {
		t.Errorf("Git: State = %v, want Failed", gitResult.State)
	}
	if pythonResult.State != StateVerified {
		t.Errorf("Python: State = %v, want Verified", pythonResult.State)
	}
}

func TestOrchestrator_Run_CertificatePassedToApply(t *testing.T) {
	git := NewMockAdapter("git", "Git")
	registry := NewMockRegistry(git)

	orch := NewOrchestrator(registry)
	cert := Certificate{Path: "/path/to/ca.pem", Content: []byte("test content")}

	orch.Run(context.Background(), cert)

	// Verify the certificate was passed correctly
	if git.appliedCert.Path != cert.Path {
		t.Errorf("Applied cert path = %v, want %v", git.appliedCert.Path, cert.Path)
	}
	if string(git.appliedCert.Content) != string(cert.Content) {
		t.Errorf("Applied cert content = %v, want %v", git.appliedCert.Content, cert.Content)
	}
}

func TestOrchestrator_Run_EmptyRegistry(t *testing.T) {
	registry := NewMockRegistry() // No adapters
	orch := NewOrchestrator(registry)
	cert := Certificate{Path: "/path/to/ca.pem", Content: []byte("test")}

	results := orch.Run(context.Background(), cert)

	if len(results) != 0 {
		t.Errorf("Run() with empty registry returned %d results, want 0", len(results))
	}
}
