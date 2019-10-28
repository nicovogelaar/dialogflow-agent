package cmd

import (
	"log"

	"github.com/nicovogelaar/dialogflow-agent/dialogflow"
	"github.com/spf13/cobra"
)

var (
	intentsImportFilename string
	intentsImportURL      string

	intentsImportCmd = &cobra.Command{
		Use: "import",
		Run: func(_ *cobra.Command, _ []string) {
			intentsClient, err := dialogflow.NewIntentsClient(projectID, credentialsFile)
			if err != nil {
				log.Fatalf("failed to create intents client: %v", err)
			}
			defer func() {
				if err = intentsClient.Close(); err != nil {
					log.Printf("failed to close intents client: %v", err)
				}
			}()

			var intentsImporter dialogflow.IntentsImporter
			if intentsImportURL != "" {
				intentsImporter = dialogflow.NewURLIntentsImporter(intentsClient, intentsImportURL)
			} else {
				intentsImporter = dialogflow.NewFileIntentsImporter(intentsClient, intentsImportFilename)
			}

			if err = intentsImporter.ImportIntents(); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	intentsImportCmd.Flags().StringVarP(&intentsImportFilename, "filename", "f", "intents.yaml", "intents filename")
	intentsImportCmd.Flags().StringVarP(&intentsImportURL, "url", "u", "", "intents url")
}
