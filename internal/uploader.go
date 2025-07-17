package internal

import (
	"log"
	"os"

	"google.golang.org/api/drive/v3"
)

type Uploader struct {
	DriveService *drive.Service
}

func NewUploader() *Uploader {
	service, err := DriveClient()
	if err != nil {
		log.Printf("Error creating Drive client: %v", err)
		return nil
	}
	return &Uploader{
		DriveService: service,
	}
}

func (u *Uploader) Upload(path string) error {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file for upload: %v", err)
	}
	file, err := u.DriveService.Files.Create(&drive.File{
		Name: path,
	}).Media(f).Do()
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return err
	}
	log.Printf("File uploaded successfully: %s", file.Id)
	return nil
}
