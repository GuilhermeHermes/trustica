package main

import (
	"os"

	"github.com/GuilhermeHermes/trustica/internal/adapter/inbound/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
