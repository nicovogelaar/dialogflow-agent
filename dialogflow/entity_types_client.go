package dialogflow

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/dialogflow/apiv2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type EntityTypesClient struct {
	projectID         string
	entityTypesClient *dialogflow.EntityTypesClient
}

func NewEntityTypesClient(projectID, credentialsFile string) (*EntityTypesClient, error) {
	ctx := context.Background()

	entityTypesClient, err := dialogflow.NewEntityTypesClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}

	return &EntityTypesClient{
		projectID:         projectID,
		entityTypesClient: entityTypesClient,
	}, nil
}

func (client *EntityTypesClient) ListEntityTypes() ([]EntityType, error) {
	iter := client.entityTypesClient.ListEntityTypes(
		context.Background(),
		&dialogflowpb.ListEntityTypesRequest{
			Parent: fmt.Sprintf("projects/%s/agent", client.projectID),
		},
	)

	var entityTypes []EntityType
	for {
		intent, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		entityTypes = append(entityTypes, toEntityType(intent))
	}

	return entityTypes, nil
}

func (client *EntityTypesClient) GetEntityType(entityTypeID string) (EntityType, error) {
	entityType, err := client.entityTypesClient.GetEntityType(
		context.Background(),
		&dialogflowpb.GetEntityTypeRequest{
			Name: fmt.Sprintf("projects/%s/agent/entityTypes/%s", client.projectID, entityTypeID),
		},
	)
	if err != nil {
		return EntityType{}, err
	}

	return toEntityType(entityType), nil
}

func (client *EntityTypesClient) CreateEntityType(entityType EntityType) (EntityType, error) {
	if entityType.DisplayName == "" {
		return EntityType{}, errors.New("display name is empty")
	}

	kind := dialogflowpb.EntityType_KIND_LIST
	if val, ok := dialogflowpb.EntityType_Kind_value[entityType.Kind]; ok {
		kind = dialogflowpb.EntityType_Kind(val)
	}
	autoExpansionMode := dialogflowpb.EntityType_AUTO_EXPANSION_MODE_UNSPECIFIED
	if val, ok := dialogflowpb.EntityType_AutoExpansionMode_value[entityType.AutoExpansionMode]; ok {
		autoExpansionMode = dialogflowpb.EntityType_AutoExpansionMode(val)
	}

	var entities []*dialogflowpb.EntityType_Entity
	for _, entity := range entityType.Entities {
		entities = append(entities, &dialogflowpb.EntityType_Entity{
			Value:    entity.Value,
			Synonyms: entity.Synonyms,
		})
	}

	dialogflowEntityType, err := client.entityTypesClient.CreateEntityType(
		context.Background(),
		&dialogflowpb.CreateEntityTypeRequest{
			Parent: fmt.Sprintf("projects/%s/agent", client.projectID),
			EntityType: &dialogflowpb.EntityType{
				DisplayName:           entityType.DisplayName,
				Kind:                  kind,
				AutoExpansionMode:     autoExpansionMode,
				Entities:              entities,
				EnableFuzzyExtraction: entityType.EnableFuzzyExtraction,
			},
		},
	)
	if err != nil {
		return EntityType{}, err
	}

	return toEntityType(dialogflowEntityType), nil
}

func (client *EntityTypesClient) DeleteEntityType(entityTypeID string) error {
	if entityTypeID == "" {
		return errors.New("missing entity type id")
	}
	err := client.entityTypesClient.DeleteEntityType(context.Background(), &dialogflowpb.DeleteEntityTypeRequest{
		Name: fmt.Sprintf("projects/%s/agent/entityTypes/%s", client.projectID, entityTypeID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (client *EntityTypesClient) DeleteEntityTypes(entityTypes []EntityType) error {
	var entityTypeNames []string
	for _, entityType := range entityTypes {
		entityTypeNames = append(entityTypeNames, entityType.Name)
	}
	_, err := client.entityTypesClient.BatchDeleteEntityTypes(
		context.Background(),
		&dialogflowpb.BatchDeleteEntityTypesRequest{
			Parent:          fmt.Sprintf("projects/%s/agent", client.projectID),
			EntityTypeNames: entityTypeNames,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (client *EntityTypesClient) Close() error {
	return client.entityTypesClient.Close()
}

func toEntityType(dialogflowEntityType *dialogflowpb.EntityType) EntityType {
	var entities []Entity
	for _, entity := range dialogflowEntityType.Entities {
		entities = append(entities, Entity{
			Value:    entity.Value,
			Synonyms: entity.Synonyms,
		})
	}
	return EntityType{
		Name:                  dialogflowEntityType.Name,
		DisplayName:           dialogflowEntityType.DisplayName,
		Kind:                  dialogflowEntityType.Kind.String(),
		AutoExpansionMode:     dialogflowEntityType.AutoExpansionMode.String(),
		Entities:              entities,
		EnableFuzzyExtraction: dialogflowEntityType.EnableFuzzyExtraction,
	}
}
