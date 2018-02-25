package main

import (
	"log"
	"os"

	"github.com/phrase/phraseapp-go/phraseapp"
	"github.com/thesoenke/translation-proxy-phraseapp/api"
)

func main() {
	token := os.Getenv("PHRASEAPP_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Please set the access token in 'PHRASEAPP_ACCESS_TOKEN'")
	}

	credentials := phraseapp.Credentials{
		Host:  "https://api.phraseapp.com",
		Token: token,
	}
	client := phraseapp.Client{
		Credentials: credentials,
	}

	// check credentials
	_, err := client.ProjectsList(0, 10)
	if err != nil {
		log.Fatal(err)
	}

	api.Run(&client)
}
