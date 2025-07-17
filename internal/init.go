package internal

import (
	"log"
	"os"
	"path"
)

func CreateInitFile() error {
	cd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return err
	}

	initFilePath := path.Join(cd, "gdrive_init.json")
	_, err = os.Stat(initFilePath)
	if os.IsNotExist(err) {
		file, err := os.Create(initFilePath)
		if err != nil {
			log.Printf("Error creating init file: %v", err)
			return err
		}
		defer file.Close()
		log.Printf("Init file created at: %s", initFilePath)
		file.WriteString(
		`{
			"version": "0"
		}`,
		)
	} else {
		log.Printf("Init already created at: %s", initFilePath)
		return nil
	}
	if err != nil {
		log.Printf("Error checking init file: %v", err)
		return err
	}
	return nil
}
