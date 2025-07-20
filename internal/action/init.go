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

		initData := BaseInit{
			Version: 0,
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
