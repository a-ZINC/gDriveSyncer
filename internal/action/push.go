package action

import (	
	"fmt"
	"gdriveSync/cmd"
	"gdriveSync/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Chann struct {
	Path    string
	Version int
	Hash    string
	IsChanged bool
	FileId string
}

type PushAction struct {
	Drive     *utils.Drive
	PathChan chan *Chann
}

func (p *PushAction) Action() error {
	err := p.Drive.ExtractInitData()
	if err != nil {
		return fmt.Errorf("failed to extract init data: %w", err)
	}
	curr_data := make(map[string]interface{})
	if p.Drive.Version != 0 {
		val, ok := utils.GetMap(p.Drive.Data, strconv.Itoa(p.Drive.Version))
		if !ok {
			return fmt.Errorf("version %d not found in init data", p.Drive.Version)
		}
		for k, v := range val {
			curr_data[k] = v
		}
	}
	p.WatchDirectory(curr_data, p.Drive.Version)
	return nil
}

func (p *PushAction) WatchDirectory(data map[string]interface{}, version int) {
	defer close(p.PathChan)
	var (
		pathFileId string
		pathHash string
	)
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
			return err
		}
		if cmd.Verbose {
			log.Printf("File: %s, Hash: %s", path, hash)
		}

		fileMap, exists := utils.GetMap(data, path)
		if exists {
			pathFileId = fileMap["fileId"].(string)
			pathHash = fileMap["hash"].(string)
		}
		p.PathChan <- &Chann{
			Path:    path,
			Version: version,
			Hash:    hash,
			IsChanged: !exists || pathHash != hash,
			FileId: pathFileId,
		}
		return nil
	})
}
