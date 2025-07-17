package internal

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func DriveClient() (*drive.Service, error) {
	credentialPath := "credentials.json"
	if _, err := os.Stat(credentialPath); os.IsNotExist(err) {
		log.Printf("Credentials file does not exist at: %s", credentialPath)
		return nil, fmt.Errorf("credentials file not found")
	}
	data, err := os.ReadFile(credentialPath)
	if err != nil {
		log.Printf("Error reading credentials file: %v", err)
		return nil, err
	}
	cfg, err := google.ConfigFromJSON(data, drive.DriveFileScope)
	if err != nil {
		log.Printf("Error creating config from JSON: %v", err)
		return nil, err
	}
	client := GetClient(cfg)
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Error creating Drive service: %v", err)
		return nil, err
	}
	log.Println("Drive service created successfully")
	return srv, nil
}
