package cmd

import (
	"github.com/spf13/cobra"
)

var intentsCmd = &cobra.Command{
	Use: "intents",
}

func init() {
	intentsCmd.AddCommand(intentsDeleteCmd)
	intentsCmd.AddCommand(intentsImportCmd)
}
