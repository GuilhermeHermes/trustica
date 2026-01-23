package cli

import (
	"github.com/spf13/cobra"
)

// Version information (set via ldflags at build time)
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "trustica",
	Short: "Certificate trust orchestration for developer environments",
	Long: `Trustica configures CA certificate trust for development tools and runtimes
without modifying the operating system trust store.

It automatically detects installed environments (Git, Python, Node.js, etc.)
and configures each one to trust your custom CA certificate.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
}
