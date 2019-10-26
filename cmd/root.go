package cmd

import (
	"github.com/spf13/cobra"
)

var (
	projectID       string
	credentialsFile string

	rootCmd = &cobra.Command{
		Use:   "dialogflow-agent",
		Short: "A dialogflow agent CLI tool",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&projectID, "project-id", "", "project ID")
	rootCmd.PersistentFlags().StringVar(&credentialsFile, "credentials-file", "credentials.json", "credentials file")
	rootCmd.AddCommand(entitiesCmd)
	rootCmd.AddCommand(intentsCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
