package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/nicovogelaar/dialogflow-agent/dialogflow"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	intentsImportFilename string

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

			if err = importIntents(intentsClient, intentsImportFilename); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	intentsImportCmd.Flags().StringVarP(&intentsImportFilename, "filename", "f", "intents.yaml", "intents filename")
}

func importIntents(intentsClient *dialogflow.IntentsClient, filename string) error {
	intents, err := readIntentsFromFile(filename)
	if err != nil {
		return fmt.Errorf("read intents from file: %v", err)
	}

	for _, intent := range intents {
		if err = createIntent(intentsClient, intent); err != nil {
			return err
		}
	}

	return nil
}

func createIntent(intentsClient *dialogflow.IntentsClient, intent dialogflow.Intent) error {
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

func createFollowupIntent(
	intentsClient *dialogflow.IntentsClient,
	followupIntent dialogflow.Intent,
	parentFollowupIntent dialogflow.Intent,
) error {
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

type intent struct {
	Name      string   `yaml:"name"`
	UserSays  []string `yaml:"usersays"`
	Responses []string `yaml:"responses"`
	FollowUp  []intent `yaml:"followup"`
}

func readIntentsFromFile(filename string) ([]dialogflow.Intent, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	var data struct {
		Intents []intent `yaml:"intents"`
	}

	err = yaml.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data: %v", err)
	}

	var intents []dialogflow.Intent
	for _, intent := range data.Intents {
		intents = append(intents, toIntent(intent))
	}

	return intents, nil
}

func toIntent(intent intent) dialogflow.Intent {
	var newIntent dialogflow.Intent
	newIntent.DisplayName = intent.Name
	for _, val := range intent.Responses {
		newIntent.Messages = append(newIntent.Messages, dialogflow.Message{Text: val})
	}
	for _, val := range intent.UserSays {
		newIntent.TrainingPhrases = append(
			newIntent.TrainingPhrases,
			dialogflow.TrainingPhrase{Parts: parseTrainingPhrase(val)},
		)
	}
	parameters := make(map[string]bool)
	for _, val := range newIntent.TrainingPhrases {
		for _, part := range val.Parts {
			if part.Alias == "" {
				continue
			}
			if _, ok := parameters[part.Alias]; ok {
				continue
			}
			parameters[part.Alias] = true
			newIntent.Parameters = append(newIntent.Parameters, dialogflow.Parameter{
				DisplayName:           part.Alias,
				Value:                 fmt.Sprintf("$%s", part.Alias),
				EntityTypeDisplayName: part.EntityType,
			})
		}
	}
	for _, val := range intent.FollowUp {
		newIntent.FollowupIntents = append(newIntent.FollowupIntents, toIntent(val))
	}
	return newIntent
}

var trainingPhraseRegexp = regexp.MustCompile(`@(\w+):(?:(\w+)|['"]([^'"]+)['"])`)

func parseTrainingPhrase(trainingPhrase string) []dialogflow.TrainingPhrasePart {
	var parts []dialogflow.TrainingPhrasePart

	matches := trainingPhraseRegexp.FindAllStringSubmatchIndex(trainingPhrase, -1)

	var n int
	for _, match := range matches {
		if match[0] > n {
			parts = append(parts, dialogflow.TrainingPhrasePart{
				Text: trainingPhrase[n:match[0]],
			})
		}

		var text string
		if match[4] >= 0 && match[5] >= 0 {
			text = trainingPhrase[match[4]:match[5]]
		} else {
			text = trainingPhrase[match[6]:match[7]]
		}

		parts = append(parts, dialogflow.TrainingPhrasePart{
			Text:        text,
			EntityType:  trainingPhrase[match[0]:match[3]],
			Alias:       trainingPhrase[match[2]:match[3]],
			UserDefined: true,
		})
		n = match[1]
	}

	if n < len(trainingPhrase) {
		parts = append(parts, dialogflow.TrainingPhrasePart{
			Text: trainingPhrase[n:],
		})
	}

	return parts
}
