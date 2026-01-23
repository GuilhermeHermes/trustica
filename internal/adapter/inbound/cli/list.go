package cli

import (
	"fmt"

	"github.com/GuilhermeHermes/trustica/internal/adapter"
	"github.com/GuilhermeHermes/trustica/internal/adapter/outbound/system"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported environments",
	Long:  `Lists all environments that Trustica can configure for CA trust.`,
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	sys := system.New()
	registry := adapter.DefaultRegistry(sys)

	adapters := registry.All()
	if len(adapters) == 0 {
		fmt.Println("No environments registered.")
		fmt.Println("\n(Adapters will be added in future versions)")
		return nil
	}

	fmt.Println("Supported environments:")
	for _, a := range adapters {
		info := a.Info()
		fmt.Printf("  %-15s %s\n", info.ID, info.Name)
	}

	return nil
}
