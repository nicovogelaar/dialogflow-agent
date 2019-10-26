package cmd

import (
	"github.com/spf13/cobra"
)

var entitiesCmd = &cobra.Command{
	Use: "entities",
}

func init() {
	entitiesCmd.AddCommand(entitiesDeleteCmd)
	entitiesCmd.AddCommand(entitiesImportCmd)
}
