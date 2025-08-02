package action

import (
	"fmt"
	"gdriveSync/cmd"
	"gdriveSync/utils"
	"strings"

	"google.golang.org/api/drive/v3"
)

type SearchAction struct {
	DriveService      *drive.Service
	MyDriveList       []*drive.File
	SharedDriveList   []*drive.File
	MyDriveListId     []*drive.File
	SharedDriveListId []*drive.File
}

func (s *SearchAction) Action() error {
	myDriveQuery := "'me' in owners and trashed = false"
	sharedDriveQuery := "sharedWithMe and trashed = false"
	myDriveNameSearch := fmt.Sprintf("name contains '%s'", cmd.SearchQuery)
	sharedDriveNameSearch := fmt.Sprintf("name contains '%s'", cmd.SearchQuery)

	// search by name
	s.ActionHelper(myDriveQuery+" and "+myDriveNameSearch, sharedDriveQuery+" and "+sharedDriveNameSearch, false)
	s.ActionHelper(myDriveQuery, sharedDriveQuery, true)

	fmt.Printf("🔍 %s%sSearch Results:%s\n", utils.Blue, utils.Bold, utils.Reset)
	if len(s.MyDriveList) != 0 || len(s.SharedDriveList) != 0 {
		s.PrintResults(s.MyDriveList, s.SharedDriveList)
	}
	if len(s.MyDriveListId) != 0 || len(s.SharedDriveListId) != 0 {
		s.PrintResults(s.MyDriveListId, s.SharedDriveList)
	}

	return nil
}

func (s *SearchAction) ActionHelper(myDriveQuery, sharedDriveQuery string, isIdSearch bool) error {
	if cmd.Type == cmd.DRIVE || cmd.Type == cmd.ALL {
		list, err := GetDriveFiles(myDriveQuery, s.DriveService)
		if err != nil {
			return err
		}
		if isIdSearch {
			s.MyDriveListId = append(s.MyDriveListId, s.FilterById(list)...)
		} else {
			s.MyDriveList = append(s.MyDriveList, list...)
		}
	}

	if cmd.Type == cmd.SHARED || cmd.Type == cmd.ALL {
		list, err := GetDriveFiles(sharedDriveQuery, s.DriveService)
		if err != nil {
			return err
		}
		if isIdSearch {
			s.SharedDriveListId = append(s.SharedDriveListId, s.FilterById(list)...)
		} else {
			s.SharedDriveList = append(s.SharedDriveList, list...)
		}
	}
	return nil
}
func looksLikeRootId(id string) bool {
    return strings.HasPrefix(id, "0") && len(id) == 19
}

func (s *SearchAction) FilterById(list []*drive.File) []*drive.File {
	var filteredList []*drive.File

	for _, file := range list {
		if file.Id == cmd.SearchQuery {
			filteredList = append(filteredList, file)
			return filteredList
		}
	}
	if len(filteredList) == 0 && looksLikeRootId(cmd.SearchQuery) {
		if dri, err := s.DriveService.Drives.Get(cmd.SearchQuery).Do(); err == nil {
			driveFile := &drive.File{
				Id:       dri.Id,
				Name:     dri.Name,
				MimeType: "application/vnd.google-apps.folder",
			}
			filteredList = append(filteredList, driveFile)
			return filteredList
		}

		if file, err := s.DriveService.Files.Get(cmd.SearchQuery).
			Fields("id, name, mimeType, parents, size").
			SupportsAllDrives(true).
			Do(); err == nil {
			filteredList = append(filteredList, file)
			return filteredList
		}
	}

	for _, file := range list {
		if strings.Contains(file.Id, cmd.SearchQuery) {
			filteredList = append(filteredList, file)
		}
	}

	return filteredList
}

func (s *SearchAction) PrintResults(myDriveList, sharedDriveList []*drive.File) {
	if cmd.Type == cmd.DRIVE || cmd.Type == cmd.ALL {
		for _, file := range myDriveList {
			var icon string
			if file.MimeType == "application/vnd.google-apps.folder" {
				icon = "📁"
			} else {
				icon = "📄"
			}

			var colorName string
			if file.MimeType == "application/vnd.google-apps.folder" {
				colorName = fmt.Sprintf("%s%s%s%s", utils.Green, utils.Bold, file.Name, utils.Reset)
			} else {
				colorName = fmt.Sprintf("%s%s%s", utils.White, file.Name, utils.Reset)
			}

			sizeLabel := ""
			if file.MimeType != "application/vnd.google-apps.folder" && file.Size > 0 {
				sizeLabel = fmt.Sprintf(" %s[%s]%s", utils.Cyan, formatSize(file.Size), utils.Reset)
			}

			parentLabel := ""
			if len(file.Parents) > 0 {
				parentLabel = fmt.Sprintf(" %s(parent: %s)%s", utils.Cyan, file.Parents[0], utils.Reset)
			}
			fmt.Printf("%s%s", utils.Blue, "─")
			fmt.Printf("   %s %s %s(%s)%s%s%s\n",
				icon, colorName, utils.Yellow, file.Id, utils.Reset,
				sizeLabel, parentLabel,
			)
		}
	}

	if cmd.Type == cmd.SHARED || cmd.Type == cmd.ALL {
		fmt.Printf("%s\n", utils.Reset)

		for _, file := range sharedDriveList {
			var icon string
			if file.MimeType == "application/vnd.google-apps.folder" {
				icon = "📁"
			} else {
				icon = "📄"
			}

			var colorName string
			if file.MimeType == "application/vnd.google-apps.folder" {
				colorName = fmt.Sprintf("%s%s%s%s", utils.Green, utils.Bold, file.Name, utils.Reset)
			} else {
				colorName = fmt.Sprintf("%s%s%s", utils.White, file.Name, utils.Reset)
			}

			sizeLabel := ""
			if file.MimeType != "application/vnd.google-apps.folder" && file.Size > 0 {
				sizeLabel = fmt.Sprintf(" %s[%s]%s", utils.Cyan, formatSize(file.Size), utils.Reset)
			}

			parentLabel := ""
			if len(file.Parents) > 0 {
				parentLabel = fmt.Sprintf(" %s(parent: %s)%s", utils.Green, file.Parents[0], utils.Reset)
			}
			fmt.Printf("%s%s", utils.Blue, "─")
			fmt.Printf("   %s %s %s(%s)%s%s%s\n",
				icon, colorName, utils.Yellow, file.Id, utils.Reset,
				sizeLabel, parentLabel,
			)
		}

	}
}
