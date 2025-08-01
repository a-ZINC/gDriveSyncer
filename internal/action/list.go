package action

import (
	"fmt"
	"gdriveSync/utils"

	"google.golang.org/api/drive/v3"
)

type ListAction struct {
	DriveService  *drive.Service
	FolderHandler map[string][]string
	FileName      map[string]*drive.File
	Shared        []string
	SharedFiles   []string
	SharedFolders []string
}

func (p *ListAction) GetDriveFiles(query string) ([]*drive.File, error) {
	files, err := p.DriveService.Files.List().Q(query).Fields("nextPageToken, files(id, name, mimeType, parents, size)").PageSize(1000).Do()
	if err != nil {
		return nil, err
	}
	return files.Files, nil

}

func (p *ListAction) Action() error {
	myDriveList, err := p.GetDriveFiles("'me' in owners and trashed = false")
	if err != nil {
		return err
	}

	rootFolder, err := p.DriveService.Files.Get("root").Fields("id", "name").Do()
	if err != nil {
		return err
	}

	for _, file := range myDriveList {
		if file.SharedWithMeTime != "" {
			p.Shared = append(p.Shared, file.Id)
		}
		p.FileName[file.Id] = file
		if len(file.Parents) != 0 {
			parentId := file.Parents[0]
			if _, exists := p.FolderHandler[parentId]; exists {
				p.FolderHandler[parentId] = append(p.FolderHandler[parentId], file.Id)
				continue
			}
			p.FolderHandler[parentId] = []string{file.Id}
		}
	}

	sharedList, err := p.GetDriveFiles("sharedWithMe = true and trashed = false")
	if err != nil {
		return err
	}

	for _, file := range sharedList {
		p.FileName[file.Id] = file
		p.SharedFiles = append(p.SharedFiles, file.Id)
		
		if file.MimeType == "application/vnd.google-apps.folder" {
			p.SharedFolders = append(p.SharedFolders, file.Id)
			p.fetchSharedFolderContents(file.Id)
		}
	}

	fmt.Printf("📁 %s%sMyDrive%s (%s)\n", utils.Blue, utils.Bold, utils.Reset, "root")
	p.VisualizeFolders(rootFolder.Id, 0, "")
	fmt.Printf("\n ----------------------------------------------------------------------------------\n")
	fmt.Printf("\n📁 %s%sShared with me%s\n", utils.Blue, utils.Bold, utils.Reset)
	p.DisplaySharedItems()

	return nil
}

func (p *ListAction) fetchSharedFolderContents(folderId string) error {
	query := fmt.Sprintf("'%s' in parents and trashed = false", folderId)
	files, err := p.GetDriveFiles(query)
	if err != nil {
		return err
	}
	
	for _, file := range files {
		p.FileName[file.Id] = file
		p.FolderHandler[folderId] = append(p.FolderHandler[folderId], file.Id)

		if file.MimeType == "application/vnd.google-apps.folder" {
			p.fetchSharedFolderContents(file.Id)
		}
	}
	
	return nil
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (p *ListAction) VisualizeFolders(root string, level int, prefix string) {
	children := p.FolderHandler[root]
	if len(children) == 0 {
		return
	}

	for i, id := range children {
		isLast := i == len(children)-1
		branch := "├── "
		newPrefix := prefix + "│   "
		if isLast {
			branch = "└── "
			newPrefix = prefix + "    "
		}

		name := p.FileName[id].Name
		mime := p.FileName[id].MimeType
		size := p.FileName[id].Size

		var icon, colorName string
		if mime == "application/vnd.google-apps.folder" {
			icon = "📁"
			colorName = fmt.Sprintf("%s%s%s", utils.Green, utils.Bold, name)
		} else {
			icon = "📄"
			colorName = name
		}

		sizeLabel := ""
		if mime != "application/vnd.google-apps.folder" && size > 0 {
			sizeLabel = fmt.Sprintf(" [%s]", formatSize(size))
		}

		fmt.Printf("%s%s%s %s (%s)%s%s\n",
			prefix, branch, icon, colorName, id, utils.Cyan, sizeLabel,
		)

		if mime == "application/vnd.google-apps.folder" {
			p.VisualizeFolders(id, level+1, newPrefix)
		}
	}
}

func (p *ListAction) DisplaySharedItems() {
    var sharedFolders, sharedFiles []*drive.File
    
    for _, fileId := range p.SharedFiles {
        file := p.FileName[fileId]
        if file.MimeType == "application/vnd.google-apps.folder" {
            sharedFolders = append(sharedFolders, file)
        } else {
            sharedFiles = append(sharedFiles, file)
        }
    }
    
    totalItems := len(sharedFolders) + len(sharedFiles)
    currentIndex := 0
    
    for _, folder := range sharedFolders {
        isLast := currentIndex == totalItems-1
        branch := "├── "
        newPrefix := "│   "
        if isLast {
            branch = "└── "
            newPrefix = "    "
        }
        
        fmt.Printf("%s📁 %s%s%s%s (%s)\n", 
            branch, utils.Green, utils.Bold, folder.Name, utils.Reset, folder.Id)

        p.VisualizeFolders(folder.Id, 0, newPrefix)
        currentIndex++
    }
    
    for _, file := range sharedFiles {
        isLast := currentIndex == totalItems-1
        branch := "├── "
        if isLast {
            branch = "└── "
        }
        
        sizeLabel := ""
        if file.Size > 0 {
            sizeLabel = fmt.Sprintf(" %s[%s]%s", utils.Cyan, formatSize(file.Size), utils.Reset)
        }
        
        fmt.Printf("%s📄 %s (%s)%s\n", branch, file.Name, file.Id, sizeLabel)
        currentIndex++
    }
}

