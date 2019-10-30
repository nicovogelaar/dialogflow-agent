package dialogflow

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type EntityTypesImporter interface {
	ImportEntityTypes() error
}

type entityTypesImporter struct {
	entityTypesClient *EntityTypesClient
	source            Source
}

func NewEntityTypesImporter(entityTypesClient *EntityTypesClient, source Source) EntityTypesImporter {
	return &entityTypesImporter{
		entityTypesClient: entityTypesClient,
		source:            source,
	}
}

func (importer *entityTypesImporter) ImportEntityTypes() error {
	data, err := ioutil.ReadAll(importer.source)
	if err != nil {
		return fmt.Errorf("failed to read data: %v", err)
	}

	entityTypes, err := readEntityTypes(data)
	if err != nil {
		return fmt.Errorf("read entity types: %v", err)
	}

	for _, entityType := range entityTypes {
		if _, err := importer.entityTypesClient.CreateEntityType(entityType); err != nil {
			return fmt.Errorf("create entity type: %v", err)
		}
	}

	return nil
}

func readEntityTypes(dat []byte) ([]EntityType, error) {
	var data struct {
		EntityTypes []struct {
			EntityType string   `json:"type"`
			Values     []string `json:"values"`
		} `json:"entities"`
	}

	if err := yaml.Unmarshal(dat, &data); err != nil {
		return nil, fmt.Errorf("unmarshal data: %v", err)
	}

	var entityTypes []EntityType
	for _, entityType := range data.EntityTypes {
		var entities []Entity
		for _, val := range entityType.Values {
			entities = append(entities, Entity{
				Value: val,
			})
		}
		entityTypes = append(entityTypes, EntityType{
			DisplayName: entityType.EntityType,
			Entities:    entities,
		})
	}

	return entityTypes, nil
}
