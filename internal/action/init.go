package action

import (
	"encoding/json"
	"fmt"
	"gdriveSync/utils"
	"log"
	"os"
	"path"

	"google.golang.org/api/drive/v3"
)

type InitAction struct{
	Service *drive.Service
}
type BaseInit struct {
	Version int    `json:"version"`
	FolderId string `json:"folderId"`
}

func (i *InitAction) Action() error {
	cd, err := os.Getwd()
	if err != nil {
		return err
	}

	initFilePath := path.Join(cd, "gdrive_init.json")
	_, err = os.Stat(initFilePath)
	if os.IsNotExist(err) {
		file, err := os.Create(initFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		log.Printf("Init file created at: %s", initFilePath)
		folder, err := CreateBaseFolder(i.Service)
		if err != nil {
			fmt.Printf("⚠️  %s%sFOLDER CREATION FAILED:%s Could not create folder: %s%s%v%s",
				utils.Red, utils.Bold, utils.Reset, utils.Red, utils.Bold, err, utils.Reset)
			os.Exit(1)
		}
		fmt.Printf("📂 %s%sFolder created in Drive: %s%s%s\n",
			utils.Green, utils.Bold, utils.Cyan, folder.Name, utils.Reset)

		initData := BaseInit{
			Version: 0,
			FolderId: folder.Id,
		}
		jsonContent, err := json.MarshalIndent(initData, "", "  ")
		if err != nil {
			return err
		}
		file.WriteString(string(jsonContent))
	} else {
		fmt.Printf("ℹ️  %s%sInit file already exists at: %s%s%s\n",
			utils.Yellow, utils.Bold, initFilePath, utils.Reset, utils.Reset)
		return nil
	}
	return nil
}

func CreateBaseFolder(srv *drive.Service) (*drive.File, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dir := path.Base(wd)
	folder := &drive.File{
		Name:     dir,
		MimeType: "application/vnd.google-apps.folder",
	}
	f, err := srv.Files.Create(folder).Do()
	if err != nil {
		return nil, err
	}
	return f, nil
}
