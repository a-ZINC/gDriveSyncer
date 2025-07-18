package internal

import (
	"google.golang.org/api/drive/v3"
)

type Uploader struct {
	DriveService *drive.Service
}

func NewUploader(service *drive.Service) *Uploader {
	return &Uploader{
		DriveService: service,
	}
}
