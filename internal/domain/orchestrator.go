package domain

import (
	"context"
)

// Orchestrator coordinates the trust configuration workflow.
//
// The Orchestrator is the central component of the Core. It:
//   - Retrieves adapters from the Registry
//   - Runs Detect on each adapter to build a Plan
//   - Runs Apply on detected environments
//   - Runs Verify on applied environments
//   - Aggregates and returns Results
//
// The Orchestrator does NOT:
//   - Know how any specific environment works
//   - Access the filesystem or OS directly
//   - Execute shell commands
//
// All environment-specific logic is delegated to EnvironmentAdapter implementations.
//
// ┌─────────────────────────────────────────────────────────────────────────────┐
// │                            ORCHESTRATOR FLOW                                │
// │                                                                             │
// │  1. DETECT  → For each adapter: Detect() → Build Plan                       │
// │  2. APPLY   → For each target in Plan: Apply(cert) → Record state           │
// │  3. VERIFY  → For each applied: Verify() → Update state                     │
// │  4. RETURN  → Aggregate all Results                                         │
// └─────────────────────────────────────────────────────────────────────────────┘
type Orchestrator struct {
	registry Registry
}

// NewOrchestrator creates a new Orchestrator with the given Registry.
// The Registry provides access to all environment adapters.
func NewOrchestrator(registry Registry) *Orchestrator {
	return &Orchestrator{
		registry: registry,
	}
}

// Run executes the complete trust configuration workflow.
//
// Flow:
//  1. Detect: Check which environments are present
//  2. Apply: Install certificate in detected environments
//  3. Verify: Confirm trust is working
//  4. Return: Aggregated results for all environments
//
// Failures in one environment do not affect others (isolation principle).
func (o *Orchestrator) Run(ctx context.Context, cert Certificate) []Result {
	// Phase 1: Detect
	plan, detectResults := o.detect(ctx)

	// If no targets, return early with detection results
	if !plan.HasTargets() {
		return detectResults
	}

	// Phase 2: Apply
	applyResults := o.apply(ctx, cert, plan)

	// Phase 3: Verify
	verifyResults := o.verify(ctx, applyResults)

	// Combine skipped/not-detected results with verified results
	results := append(detectResults, verifyResults...)
	return results
}

// detect runs detection on all registered adapters.
// Returns the execution Plan and Results for failed/skipped environments.
func (o *Orchestrator) detect(ctx context.Context) (*Plan, []Result) {
	plan := NewPlan()
	var results []Result

	adapters := o.registry.All()

	for _, adapter := range adapters {
		env := adapter.Info()

		detected, err := adapter.Detect(ctx)
		if err != nil {
			// Detection failed - record as failed result
			results = append(results, NewErrorResult(env, err))
			continue
		}

		if detected {
			plan.AddTarget(env)
		} else {
			plan.AddSkipped(env, "not detected")
			results = append(results, NewNotDetectedResult(env))
		}
	}

	return plan, results
}

// apply installs the certificate in all target environments.
// Returns Results for each environment (Applied, AlreadyTrusted, or Failed).
func (o *Orchestrator) apply(ctx context.Context, cert Certificate, plan *Plan) []Result {
	var results []Result

	adapters := o.registry.All()

	// Build a map of adapters by environment ID for quick lookup
	adapterMap := make(map[string]EnvironmentAdapter)
	for _, adapter := range adapters {
		adapterMap[adapter.Info().ID] = adapter
	}

	for _, env := range plan.Targets {
		adapter, ok := adapterMap[env.ID]
		if !ok {
			// Should not happen, but handle gracefully
			results = append(results, NewErrorResult(env, ErrAdapterNotFound))
			continue
		}

		state, err := adapter.Apply(ctx, cert)
		if err != nil {
			results = append(results, NewErrorResult(env, err))
			continue
		}

		results = append(results, NewResult(env, state, ""))
	}

	return results
}

// verify confirms trust is working for applied environments.
// Updates state to Verified on success, Failed on error.
func (o *Orchestrator) verify(ctx context.Context, applyResults []Result) []Result {
	adapters := o.registry.All()

	// Build a map of adapters by environment ID
	adapterMap := make(map[string]EnvironmentAdapter)
	for _, adapter := range adapters {
		adapterMap[adapter.Info().ID] = adapter
	}

	results := make([]Result, len(applyResults))

	for i, result := range applyResults {
		// Only verify if Apply was successful
		if result.State != StateApplied && result.State != StateAlreadyTrusted {
			results[i] = result
			continue
		}

		adapter, ok := adapterMap[result.Environment.ID]
		if !ok {
			results[i] = NewErrorResult(result.Environment, ErrAdapterNotFound)
			continue
		}

		err := adapter.Verify(ctx)
		if err != nil {
			results[i] = NewErrorResult(result.Environment, err)
			continue
		}

		results[i] = NewResult(result.Environment, StateVerified, result.Message)
	}

	return results
}