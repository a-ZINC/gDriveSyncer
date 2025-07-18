package worker

import (
	"log"
	"os"

	"google.golang.org/api/drive/v3"
)

type UploadService struct {
	DriveService *drive.Service
}

func (p *UploadService) Upload(path string) error {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file for upload: %v", err)
	}
	file, err := p.DriveService.Files.Create(&drive.File{
		Name: path,
	}).Media(f).Do()
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return err
	}
	log.Printf("File uploaded successfully: %s", file.Id)
	return nil
}