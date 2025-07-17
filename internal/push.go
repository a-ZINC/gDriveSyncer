package internal

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"path"
)

type PushPayload struct {
	data map[string]interface{}
}
var data = make(map[string]interface{})

func Push() error {
	cd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return nil
	}

	initFilePath := path.Join(cd, "gdrive_init.json")
	_, err = os.Stat(initFilePath)
	if os.IsNotExist(err) {
		log.Printf("Init file does not exist at: %s", initFilePath)
		return nil
	}
	file, err := os.Open(initFilePath)
	if err != nil {
		log.Printf("Error opening init file: %v", err)
		return nil
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		log.Printf("Error decoding init file: %v", err)
		return nil
	}
	version, ok := data["version"].(int)
	if !ok {
		log.Println("Version not found in init file")
		return nil
	}
	log.Printf("Pushing data with version: %v", version)
	var curr_data map[string]interface{}
	if version != 0 {
		curr_data, ok = data["data"].(map[string]interface{})
		if !ok {
			log.Println("Data not found in init file")
			return nil
		}
	} else {
		curr_data = make(map[string]interface{})
	}
	WatchDirectory(curr_data)
	log.Printf("Current data: %v", curr_data)

	return nil
}

func WatchDirectory(data map[string]interface{}) {
	newVersion := make(map[string]interface{})
	log.Println("Watching directory for changes...")
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error walking directory: %v", err)
			return err
		}
		if d.IsDir() {
			log.Printf("Directory found: %s", path)
			return nil
		}
		log.Printf("File found: %s", path)
		hash, err := createHash(path)
		if err != nil {
			log.Printf("Error creating hash for file %s: %v", path, err)
			return err
		}
		log.Printf("Hash for file %s: %s", path, hash)
		existingHash, exists := data[path]
		if exists && existingHash == hash {
			log.Printf("File %s has not changed.", path)
		} else {
			log.Printf("File %s has changed or is new.", path)
		}

		newVersion[path] = hash
		return nil
	})
	log.Println("Directory watching not implemented yet.")
}
