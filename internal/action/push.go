package action

import (
	"encoding/json"
	"gdriveSync/utils"
	"log"
	"os"
	"path"
	"path/filepath"
)

type PushAction struct {
	Data     map[string]interface{}
	PathChan chan string
}

func (p *PushAction) Action() error {
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
	err = json.NewDecoder(file).Decode(&p.Data)
	if err != nil {
		log.Printf("Error decoding init file: %v", err)
		return nil
	}
	log.Printf("Init file data: %v", p.Data)
	if p.Data == nil {
		log.Println("No data found in init file")
		return nil
	}
	version, ok := p.Data["version"].(int)
	if !ok {
		log.Println("Version not found in init file")
		return nil
	}
	log.Printf("Pushing data with version: %v", version)
	var curr_data map[string]interface{}
	if version != 0 {
		curr_data, ok = p.Data["data"].(map[string]interface{})
		if !ok {
			log.Println("Data not found in init file")
			return nil
		}
	} else {
		curr_data = make(map[string]interface{})
	}
	p.WatchDirectory(curr_data)
	log.Printf("Current data: %v", curr_data)

	return nil
}

func (p *PushAction) WatchDirectory(data map[string]interface{}) {
	defer close(p.PathChan)
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
		hash, err := utils.CreateHash(path)
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
		p.PathChan <- path
		newVersion[path] = hash
		return nil
	})
	log.Println("Directory watching not implemented yet.")
}
