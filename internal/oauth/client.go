package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func GetClient(config *oauth2.Config) *http.Client {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting user home directory: %v", err)
		return nil
	}
	tokDirectory := home + "/.credentials_gdrive"
	if _, err := os.Stat(tokDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(tokDirectory, 0700)
		if err != nil {
			log.Printf("Error creating token directory: %v", err)
			return nil
		}
		log.Printf("Token directory created at: %s", tokDirectory)
	}
	tokFile := tokDirectory + "/gdrive_token.json"
	if _, err := os.Stat(tokFile); os.IsNotExist(err) {
		tok, err := tokenFromWeb(config)
		if err != nil {
			log.Printf("Error retrieving token from web: %v", err)
			return nil
		}
		saveToken(tokFile, tok)
		return config.Client(context.Background(), tok)
	}
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Printf("Error retrieving token from file: %v", err)
		return nil
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("Error opening token file: %v", err)
		return nil, err
	}
	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	if err != nil {
		log.Printf("Error decoding token file: %v", err)
		return nil, err
	}
	log.Printf("Token retrieved successfully from file: %s", file)
	return t, nil
}

func tokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Visit the URL for the auth dialog: %v", authUrl)
	var code string
	log.Printf("Enter the authorization code:")
	if _, err := fmt.Scan(&code); err != nil {
		log.Printf("Error reading authorization code: %v", err)
		return nil, err
	}
	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging authorization code: %v", err)
		return nil, err	
	}
	log.Printf("Token retrieved successfully: %v", tok)
	return tok, nil
}

func saveToken(file string, token *oauth2.Token) error {
	f, err := os.Create(file)
	if err != nil {
		log.Printf("Error creating token file: %v", err)
		return err
	}
	defer f.Close()
	log.Printf("Saving token to file: %s", file)
	return json.NewEncoder(f).Encode(token)
}
