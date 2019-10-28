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

type IntentsClient struct {
	projectID     string
	intentsClient *dialogflow.IntentsClient
}

func NewIntentsClient(projectID, credentialsFile string) (*IntentsClient, error) {
	ctx := context.Background()

	intentsClient, err := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}

	return &IntentsClient{
		projectID:     projectID,
		intentsClient: intentsClient,
	}, nil
}

func (client *IntentsClient) ListIntents() ([]Intent, error) {
	iter := client.intentsClient.ListIntents(
		context.Background(),
		&dialogflowpb.ListIntentsRequest{
			Parent: fmt.Sprintf("projects/%s/agent", client.projectID),
		},
	)

	var intents []Intent
	for {
		intent, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		intents = append(intents, dialogflowIntentToIntent(intent))
	}

	return intents, nil
}

func (client *IntentsClient) CreateIntent(intent Intent) (Intent, error) {
	if intent.DisplayName == "" {
		return Intent{}, errors.New("display name is empty")
	}

	dialogflowIntent, err := client.intentsClient.CreateIntent(
		context.Background(),
		&dialogflowpb.CreateIntentRequest{
			Parent: fmt.Sprintf("projects/%s/agent", client.projectID),
			Intent: &dialogflowpb.Intent{
				DisplayName:              intent.DisplayName,
				WebhookState:             dialogflowpb.Intent_WEBHOOK_STATE_UNSPECIFIED,
				TrainingPhrases:          toDialogflowTrainingPhrases(intent.TrainingPhrases),
				Messages:                 toDialogflowIntentMessages(intent.Messages),
				Parameters:               toDialogflowParameters(intent.Parameters),
				ParentFollowupIntentName: intent.ParentFollowupIntentName,
			},
		},
	)
	if err != nil {
		return Intent{}, err
	}

	return dialogflowIntentToIntent(dialogflowIntent), nil
}

func (client *IntentsClient) CreateFollowupIntent(intent Intent, parentFollowupIntent Intent) (Intent, error) {
	if parentFollowupIntent.Name == "" {
		return Intent{}, errors.New("parent followup intent name is empty")
	}
	intent.ParentFollowupIntentName = parentFollowupIntent.Name

	return client.CreateIntent(intent)
}

func (client *IntentsClient) DeleteIntent(intentID string) error {
	if intentID == "" {
		return errors.New("missing intent id")
	}
	err := client.intentsClient.DeleteIntent(context.Background(), &dialogflowpb.DeleteIntentRequest{
		Name: fmt.Sprintf("projects/%s/agent/intents/%s", client.projectID, intentID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (client *IntentsClient) DeleteIntents(intents []Intent) error {
	var dialogflowIntents []*dialogflowpb.Intent
	for _, intent := range intents {
		dialogflowIntents = append(dialogflowIntents, &dialogflowpb.Intent{Name: intent.Name})
	}
	_, err := client.intentsClient.BatchDeleteIntents(
		context.Background(),
		&dialogflowpb.BatchDeleteIntentsRequest{
			Parent:  fmt.Sprintf("projects/%s/agent", client.projectID),
			Intents: dialogflowIntents,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (client *IntentsClient) Close() error {
	return client.intentsClient.Close()
}

func dialogflowIntentToIntent(dialogflowIntent *dialogflowpb.Intent) Intent {
	return Intent{
		Name:                     dialogflowIntent.Name,
		DisplayName:              dialogflowIntent.DisplayName,
		Priority:                 dialogflowIntent.Priority,
		IsFallback:               dialogflowIntent.IsFallback,
		TrainingPhrases:          toTrainingPhrases(dialogflowIntent.TrainingPhrases),
		Action:                   dialogflowIntent.Action,
		InputContextNames:        nil,
		OutputContexts:           nil,
		Parameters:               nil,
		Messages:                 nil,
		RootFollowupIntentName:   dialogflowIntent.RootFollowupIntentName,
		ParentFollowupIntentName: dialogflowIntent.ParentFollowupIntentName,
		FollowupIntentInfo:       nil,
	}
}

func toTrainingPhrases(dialogflowTrainingPhrases []*dialogflowpb.Intent_TrainingPhrase) []TrainingPhrase {
	var trainingPhrases []TrainingPhrase
	for _, t := range dialogflowTrainingPhrases {
		var parts []TrainingPhrasePart
		for _, p := range t.Parts {
			parts = append(parts, TrainingPhrasePart{
				Text:        p.Text,
				EntityType:  p.EntityType,
				Alias:       p.Alias,
				UserDefined: p.UserDefined,
			})
		}
		trainingPhrases = append(trainingPhrases, TrainingPhrase{
			Name:  t.Name,
			Parts: parts,
		})
	}
	return trainingPhrases
}

func toDialogflowTrainingPhrases(trainingPhrases []TrainingPhrase) []*dialogflowpb.Intent_TrainingPhrase {
	var dialogflowTrainingPhrases []*dialogflowpb.Intent_TrainingPhrase
	for _, t := range trainingPhrases {
		parts := make([]*dialogflowpb.Intent_TrainingPhrase_Part, len(t.Parts))
		for i, p := range t.Parts {
			parts[i] = &dialogflowpb.Intent_TrainingPhrase_Part{
				Text:        p.Text,
				EntityType:  p.EntityType,
				Alias:       p.Alias,
				UserDefined: p.UserDefined,
			}
		}
		dialogflowTrainingPhrases = append(dialogflowTrainingPhrases, &dialogflowpb.Intent_TrainingPhrase{
			Type:  dialogflowpb.Intent_TrainingPhrase_EXAMPLE,
			Parts: parts,
		})
	}
	return dialogflowTrainingPhrases
}

func toDialogflowIntentMessages(messages []Message) []*dialogflowpb.Intent_Message {
	var messageTexts []string
	for _, m := range messages {
		messageTexts = append(messageTexts, m.Text)
	}
	intentMessageTexts := dialogflowpb.Intent_Message_Text{Text: messageTexts}
	wrappedIntentMessageTexts := dialogflowpb.Intent_Message_Text_{Text: &intentMessageTexts}
	intentMessage := dialogflowpb.Intent_Message{Message: &wrappedIntentMessageTexts}
	return []*dialogflowpb.Intent_Message{&intentMessage}
}

func toDialogflowParameters(parameters []Parameter) []*dialogflowpb.Intent_Parameter {
	var dialogflowParameters []*dialogflowpb.Intent_Parameter
	for _, p := range parameters {
		dialogflowParameters = append(dialogflowParameters, &dialogflowpb.Intent_Parameter{
			DisplayName:           p.DisplayName,
			EntityTypeDisplayName: p.EntityTypeDisplayName,
			Value:                 p.Value,
			Mandatory:             p.Mandatory,
			IsList:                p.IsList,
		})
	}
	return dialogflowParameters
}
