package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nicovogelaar/dialogflow-agent/dialogflow"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	entitiesImportFilename string

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

			if err = importEntityTypes(entityTypesClient, entitiesImportFilename); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	entitiesImportCmd.Flags().StringVarP(&entitiesImportFilename, "filename", "f", "entities.yaml", "entities filename")
}

func importEntityTypes(entityTypesClient *dialogflow.EntityTypesClient, filename string) error {
	entityTypes, err := readEntityTypesFromFile(filename)
	if err != nil {
		return fmt.Errorf("read intents from file: %v", err)
	}

	for _, entityType := range entityTypes {
		if err = createEntityType(entityTypesClient, entityType); err != nil {
			return err
		}
	}

	return nil
}

func createEntityType(entityTypesClient *dialogflow.EntityTypesClient, entityType dialogflow.EntityType) error {
	_, err := entityTypesClient.CreateEntityType(entityType)
	if err != nil {
		return fmt.Errorf("create entity type: %v", err)
	}

	return nil
}

func readEntityTypesFromFile(filename string) ([]dialogflow.EntityType, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	var data struct {
		EntityTypes []struct {
			EntityType string   `yaml:"type"`
			Values     []string `yaml:"values"`
		} `yaml:"entities"`
	}

	err = yaml.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data: %v", err)
	}

	var entityTypes []dialogflow.EntityType
	for _, entityType := range data.EntityTypes {
		var entities []dialogflow.Entity
		for _, val := range entityType.Values {
			entities = append(entities, dialogflow.Entity{
				Value: val,
			})
		}
		entityTypes = append(entityTypes, dialogflow.EntityType{
			DisplayName: entityType.EntityType,
			Entities:    entities,
		})
	}

	return entityTypes, nil
}
