package dialogflow

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/ghodss/yaml"
)

type IntentsImporter interface {
	ImportIntents() error
}

type intentsImporter struct {
	intentsClient *IntentsClient
	source        Source
}

func NewIntentsImporter(intentsClient *IntentsClient, source Source) IntentsImporter {
	return &intentsImporter{
		intentsClient: intentsClient,
		source:        source,
	}
}

func (importer *intentsImporter) ImportIntents() error {
	data, err := ioutil.ReadAll(importer.source)
	if err != nil {
		return fmt.Errorf("read data: %v", err)
	}

	intents, err := readIntents(data)
	if err != nil {
		return fmt.Errorf("read intents: %v", err)
	}

	for _, intent := range intents {
		if err = importer.createIntent(intent); err != nil {
			return err
		}
	}

	return nil
}

func (importer *intentsImporter) createIntent(intent Intent) error {
	newIntent, err := importer.intentsClient.CreateIntent(intent)
	if err != nil {
		return fmt.Errorf("create intent: %v", err)
	}

	for _, followupIntent := range intent.FollowupIntents {
		if err := importer.createFollowupIntent(followupIntent, newIntent); err != nil {
			return err
		}
	}

	return nil
}

func (importer *intentsImporter) createFollowupIntent(followupIntent Intent, parentFollowupIntent Intent) error {
	parentFollowupIntent, err := importer.intentsClient.CreateFollowupIntent(followupIntent, parentFollowupIntent)
	if err != nil {
		return fmt.Errorf("create followup intent: %v", err)
	}

	if len(followupIntent.FollowupIntents) > 0 {
		for _, followupIntent := range followupIntent.FollowupIntents {
			err := importer.createFollowupIntent(followupIntent, parentFollowupIntent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type intentData struct {
	Name            string       `json:"name"`
	UserSays        []string     `json:"usersays"`
	Responses       []string     `json:"responses"`
	FollowupIntents []intentData `json:"followup"`
	IsFallback      bool         `json:"fallback"`
}

func readIntents(dat []byte) ([]Intent, error) {
	var data struct {
		Intents []intentData `json:"intents"`
	}

	if err := yaml.Unmarshal(dat, &data); err != nil {
		return nil, fmt.Errorf("unmarshal data: %v", err)
	}

	var intents []Intent
	for _, intent := range data.Intents {
		intents = append(intents, intentDataToIntent(intent))
	}

	return intents, nil
}

func intentDataToIntent(intentData intentData) Intent {
	var messages []Message
	for _, val := range intentData.Responses {
		messages = append(messages, Message{Text: val})
	}

	var (
		parameters      []Parameter
		trainingPhrases []TrainingPhrase
		params          = make(map[string]bool)
	)

	for _, val := range intentData.UserSays {
		parts := parseTrainingPhrase(val)
		trainingPhrases = append(trainingPhrases, TrainingPhrase{Parts: parts})
		for _, p := range parts {
			if p.Alias == "" {
				continue
			}
			if _, ok := params[p.Alias]; ok {
				continue
			}
			params[p.Alias] = true
			parameters = append(parameters, Parameter{
				DisplayName:           p.Alias,
				Value:                 fmt.Sprintf("$%s", p.Alias),
				EntityTypeDisplayName: p.EntityType,
			})
		}
	}

	var followupIntents []Intent
	for _, val := range intentData.FollowupIntents {
		followupIntents = append(followupIntents, intentDataToIntent(val))
	}

	return Intent{
		DisplayName:     intentData.Name,
		IsFallback:      intentData.IsFallback,
		TrainingPhrases: trainingPhrases,
		Messages:        messages,
		Parameters:      parameters,
		FollowupIntents: followupIntents,
	}
}

var trainingPhraseRegexp = regexp.MustCompile(`@(\w+):(?:(\w+)|['"]([^'"]+)['"])`)

func parseTrainingPhrase(trainingPhrase string) []TrainingPhrasePart {
	var parts []TrainingPhrasePart

	matches := trainingPhraseRegexp.FindAllStringSubmatchIndex(trainingPhrase, -1)

	var n int
	for _, match := range matches {
		if match[0] > n {
			parts = append(parts, TrainingPhrasePart{
				Text: trainingPhrase[n:match[0]],
			})
		}

		var text string
		if match[4] >= 0 && match[5] >= 0 {
			text = trainingPhrase[match[4]:match[5]]
		} else {
			text = trainingPhrase[match[6]:match[7]]
		}

		parts = append(parts, TrainingPhrasePart{
			Text:        text,
			EntityType:  trainingPhrase[match[0]:match[3]],
			Alias:       trainingPhrase[match[2]:match[3]],
			UserDefined: true,
		})
		n = match[1]
	}

	if n < len(trainingPhrase) {
		parts = append(parts, TrainingPhrasePart{
			Text: trainingPhrase[n:],
		})
	}

	return parts
}
