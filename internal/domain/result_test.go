package domain

import (
	"errors"
	"testing"
)

func TestNewResult(t *testing.T) {
	env := NewEnvironment("git", "Git")
	result := NewResult(env, StateApplied, "Certificate installed")

	if result.Environment != env {
		t.Errorf("Result.Environment = %v, want %v", result.Environment, env)
	}
	if result.State != StateApplied {
		t.Errorf("Result.State = %v, want %v", result.State, StateApplied)
	}
	if result.Message != "Certificate installed" {
		t.Errorf("Result.Message = %v, want %v", result.Message, "Certificate installed")
	}
	if result.Error != nil {
		t.Errorf("Result.Error = %v, want nil", result.Error)
	}
}

func TestNewErrorResult(t *testing.T) {
	env := NewEnvironment("git", "Git")
	testErr := errors.New("permission denied")
	result := NewErrorResult(env, testErr)

	if result.Environment != env {
		t.Errorf("Result.Environment = %v, want %v", result.Environment, env)
	}
	if result.State != StateFailed {
		t.Errorf("Result.State = %v, want %v", result.State, StateFailed)
	}
	if result.Error != testErr {
		t.Errorf("Result.Error = %v, want %v", result.Error, testErr)
	}
	if result.Message != testErr.Error() {
		t.Errorf("Result.Message = %v, want %v", result.Message, testErr.Error())
	}
}

func TestNewSkippedResult(t *testing.T) {
	env := NewEnvironment("git", "Git")
	result := NewSkippedResult(env, "excluded by user")

	if result.State != StateSkipped {
		t.Errorf("Result.State = %v, want %v", result.State, StateSkipped)
	}
	if result.Message != "excluded by user" {
		t.Errorf("Result.Message = %v, want %v", result.Message, "excluded by user")
	}
	if result.Error != nil {
		t.Errorf("Result.Error = %v, want nil", result.Error)
	}
}

func TestNewNotDetectedResult(t *testing.T) {
	env := NewEnvironment("git", "Git")
	result := NewNotDetectedResult(env)

	if result.State != StateNotDetected {
		t.Errorf("Result.State = %v, want %v", result.State, StateNotDetected)
	}
	if result.Message != "Environment not detected" {
		t.Errorf("Result.Message = %v, want %v", result.Message, "Environment not detected")
	}
}

func TestResult_IsSuccess(t *testing.T) {
	env := NewEnvironment("git", "Git")

	tests := []struct {
		name     string
		result   Result
		expected bool
	}{
		{"Verified", NewResult(env, StateVerified, ""), true},
		{"AlreadyTrusted", NewResult(env, StateAlreadyTrusted, ""), true},
		{"Applied", NewResult(env, StateApplied, ""), false},
		{"Failed", NewErrorResult(env, errors.New("error")), false},
		{"Skipped", NewSkippedResult(env, "reason"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.IsSuccess(); got != tt.expected {
				t.Errorf("Result.IsSuccess() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestResult_IsFailed(t *testing.T) {
	env := NewEnvironment("git", "Git")

	tests := []struct {
		name     string
		result   Result
		expected bool
	}{
		{"Verified", NewResult(env, StateVerified, ""), false},
		{"Failed", NewErrorResult(env, errors.New("error")), true},
		{"Skipped", NewSkippedResult(env, "reason"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.IsFailed(); got != tt.expected {
				t.Errorf("Result.IsFailed() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestResult_String(t *testing.T) {
	env := NewEnvironment("git", "Git")

	// With message
	result1 := NewResult(env, StateApplied, "done")
	expected1 := "Git: Applied (done)"
	if got := result1.String(); got != expected1 {
		t.Errorf("Result.String() = %v, want %v", got, expected1)
	}

	// Without message
	result2 := NewResult(env, StateVerified, "")
	expected2 := "Git: Verified"
	if got := result2.String(); got != expected2 {
		t.Errorf("Result.String() = %v, want %v", got, expected2)
	}
}
