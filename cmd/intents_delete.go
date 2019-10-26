package cmd

import (
	"fmt"
	"log"

	"github.com/nicovogelaar/dialogflow-agent/dialogflow"
	"github.com/spf13/cobra"
)

var (
	intentsDeleteAll      bool
	intentsDeleteIntentID string

	intentsDeleteCmd = &cobra.Command{
		Use: "delete",
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

			if intentsDeleteAll {
				if err := deleteAllIntents(intentsClient); err != nil {
					log.Fatal(err)
				}
				return
			}

			if err = deleteIntent(intentsClient, intentsDeleteIntentID); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	intentsDeleteCmd.Flags().BoolVarP(&intentsDeleteAll, "all", "a", false, "delete all intents")
	intentsDeleteCmd.Flags().StringVarP(&intentsDeleteIntentID, "id", "i", "", "delete intent for the given intent id")
}

func deleteIntent(intentsClient *dialogflow.IntentsClient, intentID string) error {
	if err := intentsClient.DeleteIntent(intentID); err != nil {
		return fmt.Errorf("delete intent: %v", err)
	}
	return nil
}

func deleteAllIntents(intentsClient *dialogflow.IntentsClient) error {
	intents, err := intentsClient.ListIntents()
	if err != nil {
		return fmt.Errorf("list intents: %v", err)
	}
	var deleteIntents []dialogflow.Intent
	for _, intent := range intents {
		deleteIntents = append(deleteIntents, dialogflow.Intent{Name: intent.Name})
	}
	err = intentsClient.DeleteIntents(deleteIntents)
	if err != nil {
		return fmt.Errorf("delete intents: %v", err)
	}
	return nil
}
