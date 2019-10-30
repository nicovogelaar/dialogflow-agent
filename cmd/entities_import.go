package cmd

import (
	"log"

	"github.com/nicovogelaar/dialogflow-agent/dialogflow"
	"github.com/spf13/cobra"
)

var (
	entitiesImportFilename string
	entitiesImportURL      string

	entitiesImportCmd = &cobra.Command{
		Use: "import",
		Run: func(_ *cobra.Command, _ []string) {
			entityTypesClient, err := dialogflow.NewEntityTypesClient(projectID, credentialsFile)
			if err != nil {
				log.Fatalf("failed to create entity types client: %v", err)
			}
			defer func() {
				if err = entityTypesClient.Close(); err != nil {
					log.Printf("failed to close entity types client: %v", err)
				}
			}()

			var source dialogflow.Source
			if intentsImportURL != "" {
				source = dialogflow.NewURLSource(entitiesImportURL)
			} else {
				source = dialogflow.NewFileSource(entitiesImportFilename)
			}

			importer := dialogflow.NewEntityTypesImporter(entityTypesClient, source)
			if err = importer.ImportEntityTypes(); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	entitiesImportCmd.Flags().StringVarP(&entitiesImportFilename, "filename", "f", "entities.yaml", "entities filename")
	entitiesImportCmd.Flags().StringVarP(&entitiesImportURL, "url", "u", "", "entities url")
}
