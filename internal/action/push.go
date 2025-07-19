package action

import (
	"encoding/json"
	"fmt"
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
	IsChanged bool
}

type PushAction struct {
	Data     map[string]interface{}
	PathChan chan *Chann
}

func (p *PushAction) Action() error {
	cd, err := os.Getwd()
	if err != nil {
		return err
	}

	initFilePath := path.Join(cd, "gdrive_init.json")
	_, err = os.Stat(initFilePath)
	if os.IsNotExist(err) {
		return err
	}
	file, err := os.Open(initFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&p.Data)
	if err != nil {
		return err
	}
	if p.Data == nil {
		return err
	}
	ver, ok := p.Data["version"]
	if !ok {
		return fmt.Errorf("version not found in init file")
	}
	version, err := strconv.Atoi(ver.(string))
	if err != nil {
		return err
	}
	curr_data := make(map[string]interface{})
	if version != 0 {
		val, ok := p.Data[strconv.Itoa(version)].(map[string]interface{})
		if !ok {
			return fmt.Errorf("no data found for version %d", version)
		}
		for k, v := range val {
			curr_data[k] = v
		}
	}
	p.WatchDirectory(curr_data, version)
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
		p.PathChan <- &Chann{
			Path:    path,
			Version: version,
			Hash:    hash,
			IsChanged: !exists || existingHash != hash,
		}
		return nil
	})
}
