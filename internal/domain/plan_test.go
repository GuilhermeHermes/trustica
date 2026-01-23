package domain

import "testing"

func TestNewPlan(t *testing.T) {
	plan := NewPlan()

	if plan == nil {
		t.Fatal("NewPlan() returned nil")
	}
	if len(plan.Targets) != 0 {
		t.Errorf("NewPlan().Targets = %v, want empty", plan.Targets)
	}
	if len(plan.Skipped) != 0 {
		t.Errorf("NewPlan().Skipped = %v, want empty", plan.Skipped)
	}
	if len(plan.SkipReasons) != 0 {
		t.Errorf("NewPlan().SkipReasons = %v, want empty", plan.SkipReasons)
	}
}

func TestPlan_AddTarget(t *testing.T) {
	plan := NewPlan()
	env := NewEnvironment("git", "Git")

	plan.AddTarget(env)

	if len(plan.Targets) != 1 {
		t.Fatalf("Plan.Targets length = %d, want 1", len(plan.Targets))
	}
	if plan.Targets[0] != env {
		t.Errorf("Plan.Targets[0] = %v, want %v", plan.Targets[0], env)
	}
}

func TestPlan_AddSkipped(t *testing.T) {
	plan := NewPlan()
	env := NewEnvironment("nodejs", "Node.js")
	reason := "not detected"

	plan.AddSkipped(env, reason)

	if len(plan.Skipped) != 1 {
		t.Fatalf("Plan.Skipped length = %d, want 1", len(plan.Skipped))
	}
	if plan.Skipped[0] != env {
		t.Errorf("Plan.Skipped[0] = %v, want %v", plan.Skipped[0], env)
	}
	if plan.SkipReasons[env.ID] != reason {
		t.Errorf("Plan.SkipReasons[%s] = %v, want %v", env.ID, plan.SkipReasons[env.ID], reason)
	}
}

func TestPlan_HasTargets(t *testing.T) {
	plan := NewPlan()

	if plan.HasTargets() {
		t.Error("Empty plan.HasTargets() = true, want false")
	}

	plan.AddTarget(NewEnvironment("git", "Git"))

	if !plan.HasTargets() {
		t.Error("Plan with target.HasTargets() = false, want true")
	}
}

func TestPlan_TargetCount(t *testing.T) {
	plan := NewPlan()

	if got := plan.TargetCount(); got != 0 {
		t.Errorf("Empty plan.TargetCount() = %d, want 0", got)
	}

	plan.AddTarget(NewEnvironment("git", "Git"))
	plan.AddTarget(NewEnvironment("python", "Python"))

	if got := plan.TargetCount(); got != 2 {
		t.Errorf("Plan.TargetCount() = %d, want 2", got)
	}
}

func TestPlan_SkippedCount(t *testing.T) {
	plan := NewPlan()

	if got := plan.SkippedCount(); got != 0 {
		t.Errorf("Empty plan.SkippedCount() = %d, want 0", got)
	}

	plan.AddSkipped(NewEnvironment("nodejs", "Node.js"), "not installed")

	if got := plan.SkippedCount(); got != 1 {
		t.Errorf("Plan.SkippedCount() = %d, want 1", got)
	}
}

func TestPlan_GetSkipReason(t *testing.T) {
	plan := NewPlan()
	env := NewEnvironment("nodejs", "Node.js")
	reason := "not installed"

	plan.AddSkipped(env, reason)

	if got := plan.GetSkipReason(env.ID); got != reason {
		t.Errorf("Plan.GetSkipReason(%s) = %v, want %v", env.ID, got, reason)
	}

	// Non-existent environment
	if got := plan.GetSkipReason("unknown"); got != "" {
		t.Errorf("Plan.GetSkipReason(unknown) = %v, want empty", got)
	}
}

func TestPlan_MultipleOperations(t *testing.T) {
	plan := NewPlan()

	// Add multiple targets and skipped
	git := NewEnvironment("git", "Git")
	python := NewEnvironment("python", "Python")
	nodejs := NewEnvironment("nodejs", "Node.js")
	openssl := NewEnvironment("openssl", "OpenSSL")

	plan.AddTarget(git)
	plan.AddTarget(python)
	plan.AddSkipped(nodejs, "not installed")
	plan.AddSkipped(openssl, "excluded by user")

	if plan.TargetCount() != 2 {
		t.Errorf("Plan.TargetCount() = %d, want 2", plan.TargetCount())
	}
	if plan.SkippedCount() != 2 {
		t.Errorf("Plan.SkippedCount() = %d, want 2", plan.SkippedCount())
	}
	if plan.GetSkipReason("nodejs") != "not installed" {
		t.Errorf("Wrong skip reason for nodejs")
	}
	if plan.GetSkipReason("openssl") != "excluded by user" {
		t.Errorf("Wrong skip reason for openssl")
	}
}
