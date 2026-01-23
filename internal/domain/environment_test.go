package domain

import "testing"

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment("git", "Git")

	if env.ID != "git" {
		t.Errorf("Environment.ID = %v, want %v", env.ID, "git")
	}
	if env.Name != "Git" {
		t.Errorf("Environment.Name = %v, want %v", env.Name, "Git")
	}
}

func TestEnvironment_String(t *testing.T) {
	tests := []struct {
		env      Environment
		expected string
	}{
		{NewEnvironment("git", "Git"), "Git"},
		{NewEnvironment("python", "Python (certifi)"), "Python (certifi)"},
		{NewEnvironment("nodejs", "Node.js"), "Node.js"},
		{NewEnvironment("openssl", "OpenSSL"), "OpenSSL"},
	}

	for _, tt := range tests {
		t.Run(tt.env.ID, func(t *testing.T) {
			if got := tt.env.String(); got != tt.expected {
				t.Errorf("Environment.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEnvironment_Equality(t *testing.T) {
	env1 := NewEnvironment("git", "Git")
	env2 := NewEnvironment("git", "Git")
	env3 := NewEnvironment("python", "Python")

	if env1 != env2 {
		t.Errorf("Expected env1 == env2 (same values)")
	}
	if env1 == env3 {
		t.Errorf("Expected env1 != env3 (different values)")
	}
}
