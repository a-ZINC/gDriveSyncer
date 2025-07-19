package action

import (
	"encoding/json"
	"gdriveSync/cmd"
	"log"
	"os"
	"path"
)

type InitAction struct{}

var initData = map[string]interface{}{
	"version": "0",
}

func (i *InitAction) Action() error {
	cd, err := os.Getwd()
	if err != nil {
		if cmd.Verbose {
			log.Printf("Error getting current directory: %v", err)
		}
		return err
	}

	initFilePath := path.Join(cd, "gdrive_init.json")
	_, err = os.Stat(initFilePath)
	if os.IsNotExist(err) {
		file, err := os.Create(initFilePath)
		if err != nil {
			if cmd.Verbose {
				log.Printf("Error creating init file: %v", err)
			}
			return err
		}
		defer file.Close()
		log.Printf("Init file created at: %s", initFilePath)
		jsonContent, err := json.MarshalIndent(initData, "", "  ")
		if err != nil {
			if cmd.Verbose {
				log.Printf("Error marshalling init data to JSON: %v", err)
			}
			return err
		}
		file.WriteString(string(jsonContent))
	} else {
		log.Printf("Init file already exists at: %s", initFilePath)
		return nil
	}
	return nil
}
