package action

import (
	"encoding/json"
	"gdriveSync/cmd"
	"gdriveSync/utils"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type Chann struct {
	Path    string
	Version int
	Hash    string
}

type PushAction struct {
	Data     map[string]interface{}
	PathChan chan *Chann
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
	ver, ok := p.Data["version"]
	if !ok {
		log.Println("Version not found in init file")
		return nil
	}
	version, err := strconv.Atoi(ver.(string))
	if err != nil {
		log.Printf("Error converting version to int: %v", err)
		return nil
	}
	log.Printf("Pushing data with version: %v", version)
	curr_data := make(map[string]interface{})
	if version != 0 {
		val, ok := p.Data[strconv.Itoa(version)].(map[string]interface{})
		if !ok {
			log.Printf("Error converting data to map for version %d", version)
			return nil
		}
		for k, v := range val {
			curr_data[k] = v
			log.Printf("Current data for version %d: %s = %v", version, k, v)
		}
	}
	p.WatchDirectory(curr_data, version)
	log.Printf("Current data: %v", curr_data)

	return nil
}

func (p *PushAction) WatchDirectory(data map[string]interface{}, version int) {
	defer close(p.PathChan)
	if cmd.Verbose {
		log.Println("Verbose mode enabled. Watching directory for changes...")
	}
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			if cmd.Verbose {
				log.Printf("Error walking directory: %v", err)
			}
			return err
		}
		if d.IsDir() {
			if cmd.Verbose {
				log.Printf("Skipping directory: %s", path)
			}
			return nil
		}
		hash, err := utils.CreateHash(path)
		if err != nil {
			log.Printf("Error creating hash for file %s: %v", path, err)
			return err
		}
		if cmd.Verbose {
			log.Printf("File: %s, Hash: %s", path, hash)
		}
		existingHash, exists := data[path]
		if exists && existingHash == hash {
			if cmd.Verbose {
				log.Printf("File %s has not changed, skipping upload.", path)
			}
			return nil
		} else {
			if cmd.Verbose {
				log.Printf("File %s has changed, preparing for upload.", path)
			}
		}
		p.PathChan <- &Chann{
			Path:    path,
			Version: version,
			Hash:    hash,
		}
		return nil
	})
}
