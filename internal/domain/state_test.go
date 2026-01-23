package domain

import "testing"

func TestState_String(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{StateNotDetected, "NotDetected"},
		{StateDetected, "Detected"},
		{StateApplied, "Applied"},
		{StateAlreadyTrusted, "AlreadyTrusted"},
		{StateVerified, "Verified"},
		{StateFailed, "Failed"},
		{StateSkipped, "Skipped"},
		{State(99), "Unknown"}, // Unknown state
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("State.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestState_IsTerminal(t *testing.T) {
	tests := []struct {
		state    State
		expected bool
	}{
		{StateNotDetected, true},
		{StateDetected, false},
		{StateApplied, false},
		{StateAlreadyTrusted, true},
		{StateVerified, true},
		{StateFailed, true},
		{StateSkipped, true},
	}

	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			if got := tt.state.IsTerminal(); got != tt.expected {
				t.Errorf("State.IsTerminal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestState_IsSuccess(t *testing.T) {
	tests := []struct {
		state    State
		expected bool
	}{
		{StateNotDetected, false},
		{StateDetected, false},
		{StateApplied, false},
		{StateAlreadyTrusted, true},
		{StateVerified, true},
		{StateFailed, false},
		{StateSkipped, false},
	}

	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			if got := tt.state.IsSuccess(); got != tt.expected {
				t.Errorf("State.IsSuccess() = %v, want %v", got, tt.expected)
			}
		})
	}
}
