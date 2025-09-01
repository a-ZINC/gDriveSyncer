package action

import "google.golang.org/api/drive/v3"

type CloneAction struct {
	DriveService *drive.Service
	Destination string
	FolderId string
}

func (p *CloneAction) Action() error {
	p.DriveService.Files.List()
	return nil
}
