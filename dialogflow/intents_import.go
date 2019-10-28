package dialogflow

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/ghodss/yaml"
)

type IntentsImporter interface {
	ImportIntents() error
}

type fileIntentsImporter struct {
	intentsClient *IntentsClient
	filename      string
}

func NewFileIntentsImporter(intentsClient *IntentsClient, filename string) IntentsImporter {
	return &fileIntentsImporter{
		intentsClient: intentsClient,
		filename:      filename,
	}
}

func (importer *fileIntentsImporter) ImportIntents() error {
	intents, err := readIntentsFromFile(importer.filename)
	if err != nil {
		return fmt.Errorf("read intents from file: %v", err)
	}

	for _, intent := range intents {
		if err = createIntent(importer.intentsClient, intent); err != nil {
			return err
		}
	}

	return nil
}

type urlIntentsImporter struct {
	intentsClient *IntentsClient
	url           string
}

func NewURLIntentsImporter(intentsClient *IntentsClient, url string) IntentsImporter {
	return &urlIntentsImporter{
		intentsClient: intentsClient,
		url:           url,
	}
}

func (importer *urlIntentsImporter) ImportIntents() error {
	intents, err := readIntentsFromURL(importer.url)
	if err != nil {
		return fmt.Errorf("read intents from url: %v", err)
	}

	for _, intent := range intents {
		if err = createIntent(importer.intentsClient, intent); err != nil {
			return err
		}
	}

	return nil
}

func createIntent(intentsClient *IntentsClient, intent Intent) error {
	newIntent, err := intentsClient.CreateIntent(intent)
	if err != nil {
		return fmt.Errorf("create intent: %v", err)
	}

	for _, followupIntent := range intent.FollowupIntents {
		if err := createFollowupIntent(intentsClient, followupIntent, newIntent); err != nil {
			return err
		}
	}

	return nil
}

func createFollowupIntent(intentsClient *IntentsClient, followupIntent Intent, parentFollowupIntent Intent) error {
	parentFollowupIntent, err := intentsClient.CreateFollowupIntent(followupIntent, parentFollowupIntent)
	if err != nil {
		return fmt.Errorf("create followup intent: %v", err)
	}

	if len(followupIntent.FollowupIntents) > 0 {
		for _, followupIntent := range followupIntent.FollowupIntents {
			err := createFollowupIntent(intentsClient, followupIntent, parentFollowupIntent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func readIntentsFromFile(filename string) ([]Intent, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	return readIntents(data)
}

func readIntentsFromURL(url string) ([]Intent, error) {
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

	intents, err := readIntents(data)
	if err != nil {
		return nil, err
	}

	return intents, err
}

type intentData struct {
	Name      string       `json:"name"`
	UserSays  []string     `json:"usersays"`
	Responses []string     `json:"responses"`
	FollowUp  []intentData `json:"followup"`
}

func readIntents(dat []byte) ([]Intent, error) {
	var data struct {
		Intents []intentData `json:"intents"`
	}

	err := yaml.Unmarshal(dat, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data: %v", err)
	}

	var intents []Intent
	for _, intent := range data.Intents {
		intents = append(intents, intentDataToIntent(intent))
	}

	return intents, err
}

func intentDataToIntent(intentData intentData) Intent {
	var intent Intent
	intent.DisplayName = intentData.Name
	for _, val := range intentData.Responses {
		intent.Messages = append(intent.Messages, Message{Text: val})
	}
	parameters := make(map[string]bool)
	for _, val := range intentData.UserSays {
		parts := parseTrainingPhrase(val)
		intent.TrainingPhrases = append(intent.TrainingPhrases, TrainingPhrase{Parts: parts})
		for _, p := range parts {
			if p.Alias == "" {
				continue
			}
			if _, ok := parameters[p.Alias]; ok {
				continue
			}
			parameters[p.Alias] = true
			intent.Parameters = append(intent.Parameters, Parameter{
				DisplayName:           p.Alias,
				Value:                 fmt.Sprintf("$%s", p.Alias),
				EntityTypeDisplayName: p.EntityType,
			})
		}
	}
	for _, val := range intentData.FollowUp {
		intent.FollowupIntents = append(intent.FollowupIntents, intentDataToIntent(val))
	}
	return intent
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
