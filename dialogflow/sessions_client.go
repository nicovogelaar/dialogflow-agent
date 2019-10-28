package dialogflow

import (
	"context"
	"fmt"

	"cloud.google.com/go/dialogflow/apiv2"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type SessionsClient struct {
	projectID      string
	sessionsClient *dialogflow.SessionsClient
}

func NewSessionsClient(projectID, credentialsFile string) (*SessionsClient, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create sessions client: %v", err)
	}

	return &SessionsClient{projectID: projectID, sessionsClient: sessionClient}, nil
}

func (client *SessionsClient) DetectIntentText(sessionID, text, languageCode string) (string, error) {
	ctx := context.Background()

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", client.projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := client.sessionsClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", fmt.Errorf("failed to detect intent client: %v", err)
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	return fulfillmentText, err
}

func (client *SessionsClient) Close() error {
	return client.sessionsClient.Close()
}
