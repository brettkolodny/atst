package main

import (
	"github.com/brettkolodny/atst/internal/cli"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			cli.Simple(args)
		},
	}

	cmd.Execute()
}
