package main

import (
	"log"

	"github.com/nicovogelaar/dialogflow-agent/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
