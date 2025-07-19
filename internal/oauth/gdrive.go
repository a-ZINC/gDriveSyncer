package oauth

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
		return nil, fmt.Errorf("credentials file not found")
	}
	data, err := os.ReadFile(credentialPath)
	if err != nil {
		return nil, err
	}
	cfg, err := google.ConfigFromJSON(data, drive.DriveFileScope)
	if err != nil {
		return nil, err
	}
	client, err := GetClient(cfg)
	if err != nil {
		return nil, err
	}
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	log.Println("Drive service created successfully")
	return srv, nil
}
