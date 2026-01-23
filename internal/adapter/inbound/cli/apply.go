package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/GuilhermeHermes/trustica/internal/adapter"
	"github.com/GuilhermeHermes/trustica/internal/adapter/outbound/system"
	"github.com/GuilhermeHermes/trustica/internal/domain"
	"github.com/spf13/cobra"
)

var (
	dryRun bool
)

var applyCmd = &cobra.Command{
	Use:   "apply <certificate.pem>",
	Short: "Apply CA certificate trust to detected environments",
	Long: `Detects installed development environments and configures each one
to trust the specified CA certificate.

The certificate must be in PEM format.

Examples:
  trustica apply ./my-ca.pem
  trustica apply /path/to/corporate-ca.pem --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: runApply,
}

func init() {
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")
}

func runApply(cmd *cobra.Command, args []string) error {
	certPath := args[0]

	// Load and validate certificate
	cert, err := domain.LoadCertificate(certPath)
	if err != nil {
		return fmt.Errorf("failed to load certificate: %w", err)
	}

	fmt.Printf("Certificate: %s\n", cert.Path)

	// Create context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nInterrupted, cleaning up...")
		cancel()
	}()

	// Create infrastructure
	sys := system.New()
	registry := adapter.DefaultRegistry(sys)

	if dryRun {
		return runDryRun(registry)
	}

	// Run orchestration
	orchestrator := domain.NewOrchestrator(registry)
	results := orchestrator.Run(ctx, cert)

	// Print results
	printResults(results)

	// Determine exit code
	if hasFailures(results) {
		return fmt.Errorf("some environments failed")
	}

	return nil
}

func runDryRun(registry *adapter.InMemoryRegistry) error {
	fmt.Println("\n[Dry Run] Would process these environments:")

	adapters := registry.All()
	if len(adapters) == 0 {
		fmt.Println("  (no adapters registered yet)")
		return nil
	}

	for _, a := range adapters {
		info := a.Info()
		fmt.Printf("  - %s (%s)\n", info.Name, info.ID)
	}

	return nil
}

func printResults(results []domain.Result) {
	if len(results) == 0 {
		fmt.Println("\nNo environments processed.")
		return
	}

	fmt.Println("\nResults:")
	for _, r := range results {
		status := resultStatus(r)
		fmt.Printf("  %-20s %s\n", r.Environment.Name, status)
		if r.Error != nil {
			fmt.Printf("    Error: %v\n", r.Error)
		}
	}

	// Summary
	success, failed, skipped := countResults(results)
	fmt.Printf("\nSummary: %d succeeded, %d failed, %d skipped\n", success, failed, skipped)
}

func resultStatus(r domain.Result) string {
	switch r.State {
	case domain.StateVerified:
		return "[OK] Verified"
	case domain.StateAlreadyTrusted:
		return "[OK] Already trusted"
	case domain.StateApplied:
		return "[OK] Applied"
	case domain.StateFailed:
		return "[FAIL]"
	case domain.StateSkipped:
		return "[SKIP]"
	case domain.StateNotDetected:
		return "[--] Not detected"
	default:
		return "[??] " + r.State.String()
	}
}

func countResults(results []domain.Result) (success, failed, skipped int) {
	for _, r := range results {
		switch {
		case r.State.IsSuccess():
			success++
		case r.State == domain.StateFailed:
			failed++
		default:
			skipped++
		}
	}
	return
}

func hasFailures(results []domain.Result) bool {
	for _, r := range results {
		if r.State == domain.StateFailed {
			return true
		}
	}
	return false
}
