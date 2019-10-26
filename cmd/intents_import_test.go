package cmd

import (
	"reflect"
	"testing"

	"github.com/nicovogelaar/dialogflow-agent/dialogflow"
)

func TestParseTrainingPhraseParts(t *testing.T) {
	trainingPhrase := `Lorem ipsum @entity:dolor sit amet, @anotherentity:'consectetur adipiscing' elit.`

	parts := parseTrainingPhrase(trainingPhrase)

	expected := []dialogflow.TrainingPhrasePart{
		{
			Text: "Lorem ipsum ",
		},
		{
			Text:        "dolor",
			EntityType:  "@entity",
			Alias:       "entity",
			UserDefined: true,
		},
		{
			Text: " sit amet, ",
		},
		{
			Text:        "consectetur adipiscing",
			EntityType:  "@anotherentity",
			Alias:       "anotherentity",
			UserDefined: true,
		},
		{
			Text: " elit.",
		},
	}

	if !reflect.DeepEqual(expected, parts) {
		t.Fail()
	}
}
