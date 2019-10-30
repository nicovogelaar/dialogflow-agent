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

			var source dialogflow.Source
			if intentsImportURL != "" {
				source = dialogflow.NewURLSource(intentsImportURL)
			} else {
				source = dialogflow.NewFileSource(intentsImportFilename)
			}

			importer := dialogflow.NewIntentsImporter(intentsClient, source)
			if err = importer.ImportIntents(); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	intentsImportCmd.Flags().StringVarP(&intentsImportFilename, "filename", "f", "intents.yaml", "intents filename")
	intentsImportCmd.Flags().StringVarP(&intentsImportURL, "url", "u", "", "intents url")
}
