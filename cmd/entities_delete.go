package cmd

import (
	"fmt"
	"log"

	"github.com/nicovogelaar/dialogflow-agent/dialogflow"
	"github.com/spf13/cobra"
)

var (
	entitiesDeleteAll          bool
	entitiesDeleteEntityTypeID string

	entitiesDeleteCmd = &cobra.Command{
		Use: "delete",
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

			if entitiesDeleteAll {
				if err := deleteAllEntityTypes(entityTypesClient); err != nil {
					log.Fatal(err)
				}
				return
			}

			if err = deleteEntityType(entityTypesClient, entitiesDeleteEntityTypeID); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	entitiesDeleteCmd.Flags().BoolVarP(&entitiesDeleteAll, "all", "a", false, "delete all entities")
	entitiesDeleteCmd.Flags().StringVarP(&entitiesDeleteEntityTypeID, "id", "i", "", "delete entities for the given entity type id")
}

func deleteEntityType(entityTypesClient *dialogflow.EntityTypesClient, entityTypeID string) error {
	if err := entityTypesClient.DeleteEntityType(entityTypeID); err != nil {
		return fmt.Errorf("delete entity type: %v", err)
	}
	return nil
}

func deleteAllEntityTypes(entityTypesClient *dialogflow.EntityTypesClient) error {
	entityTypes, err := entityTypesClient.ListEntityTypes()
	if err != nil {
		return fmt.Errorf("list entity types: %v", err)
	}
	err = entityTypesClient.DeleteEntityTypes(entityTypes)
	if err != nil {
		return fmt.Errorf("delete entity types: %v", err)
	}
	return nil
}
