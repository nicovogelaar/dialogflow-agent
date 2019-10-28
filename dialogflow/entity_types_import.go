package dialogflow

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ghodss/yaml"
)

type EntityTypesImporter interface {
	ImportEntityTypes() error
}

type fileEntityTypesImporter struct {
	entityTypesClient *EntityTypesClient
	filename          string
}

func NewFileEntityTypesImporter(entityTypesClient *EntityTypesClient, filename string) EntityTypesImporter {
	return &fileEntityTypesImporter{
		entityTypesClient: entityTypesClient,
		filename:          filename,
	}
}

func (importer *fileEntityTypesImporter) ImportEntityTypes() error {
	entityTypes, err := readEntityTypesFromFile(importer.filename)
	if err != nil {
		return fmt.Errorf("read entity types from file: %v", err)
	}

	for _, entityType := range entityTypes {
		if err = createEntityType(importer.entityTypesClient, entityType); err != nil {
			return err
		}
	}

	return nil
}

type urlEntityTypesImporter struct {
	entityTypesClient *EntityTypesClient
	url               string
}

func NewURLEntityTypesImporter(entityTypesClient *EntityTypesClient, url string) EntityTypesImporter {
	return &urlEntityTypesImporter{
		entityTypesClient: entityTypesClient,
		url:               url,
	}
}

func (importer *urlEntityTypesImporter) ImportEntityTypes() error {
	entityTypes, err := readEntityTypesFromURL(importer.url)
	if err != nil {
		return fmt.Errorf("read entity types from url: %v", err)
	}

	for _, entityType := range entityTypes {
		if err = createEntityType(importer.entityTypesClient, entityType); err != nil {
			return err
		}
	}

	return nil
}

func createEntityType(entityTypesClient *EntityTypesClient, entityType EntityType) error {
	_, err := entityTypesClient.CreateEntityType(entityType)
	if err != nil {
		return fmt.Errorf("create entity type: %v", err)
	}

	return nil
}

func readEntityTypesFromFile(filename string) ([]EntityType, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	return readEntityTypes(data)
}

func readEntityTypesFromURL(url string) ([]EntityType, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); err == nil {
			err = closeErr
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %v", err)
	}

	entityTypes, err := readEntityTypes(data)
	if err != nil {
		return nil, err
	}

	return entityTypes, err
}

func readEntityTypes(dat []byte) ([]EntityType, error) {
	var data struct {
		EntityTypes []struct {
			EntityType string   `json:"type"`
			Values     []string `json:"values"`
		} `json:"entities"`
	}

	err := yaml.Unmarshal(dat, &data)
	if err != nil {
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
