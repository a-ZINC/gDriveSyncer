package worker

import (
	"encoding/json"
	"fmt"
	"gdriveSync/cmd"
	"gdriveSync/internal/action"
	"gdriveSync/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"google.golang.org/api/drive/v3"
)

type UploadService struct {
	Service     *drive.Service
	NewInitData map[string]interface{}
	OldInitData map[string]interface{}
	Stack       *utils.Stack[action.FolderStack]
}
type fileData struct {
	Hash   string `json:"hash"`
	FileId string `json:"fileId"`
}

func (p *UploadService) UploadFile(dir string, name string, version int, hash string, isChanged bool, fileId string, folderId string) error {
	ver := strconv.Itoa(version)
	path := filepath.Join(dir, name)
	var file *drive.File
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if isChanged {
		f, err := os.Open(path)
		if err != nil {
			if cmd.Verbose {
				log.Printf("Error opening file %s for upload: %v", path, err)
			}
			return err
		}
		defer f.Close()
		file, err = p.FileCreateOrUpdate(name, fileId, folderId, f)
		if err != nil {
			if cmd.Verbose {
				log.Printf("Error uploading file %s: %v", path, err)
			}
			return err
		}
		fmt.Printf("%s%s▶ SUCCESS:%s File %s%s%s uploaded!\n",
			utils.Green, utils.Bold, utils.Reset,
			utils.Cyan, path, utils.Reset)
	} else {
		fmt.Printf("%s%s▶ SKIPPED:%s File %s%s%s not changed!\n",
			utils.Green, utils.Bold, utils.Reset,
			utils.Cyan, path, utils.Reset)
	}
	if file != nil {
		fileId = file.Id
	}
	newMap, exists := utils.GetCurrentDirectoryMap(p.NewInitData, ver, wd, dir)
	if !exists {
		return fmt.Errorf("failed to get current directory map for %s", dir)
	}
	newMap[name] = fileData{
		Hash:   hash,
		FileId: fileId,
	}
	return nil
}

func (p *UploadService) UpdateWithVersion() error {
	f, err := os.OpenFile("gdrive_init.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for k, v := range p.NewInitData {
		if _, exists := p.OldInitData[k]; !exists {
			p.OldInitData[k] = v
			log.Printf("Adding new version %s to init file", k)
		} else {
			log.Printf("Version %s already exists in init file, skipping", k)
		}
	}
	if _, err = f.Seek(0, 0); err != nil {
		return err
	}
	if err = f.Truncate(0); err != nil {
		return err
	}
	err = json.NewEncoder(f).Encode(p.OldInitData)
	if err != nil {
		return err
	}
	return nil
}

func (p *UploadService) FileCreateOrUpdate(path string, fileId string, folderId string, f *os.File) (*drive.File, error) {
	var err error
	var file *drive.File
	if fileId != "" {
		file, err = p.Service.Files.Update(fileId, &drive.File{
			Name: path,
		}).Media(f).Do()
		if err != nil {
			if cmd.Verbose {
				log.Printf("Error updating file %s: %v", path, err)
			}
			return nil, err
		}
	} else {
		file, err = p.Service.Files.Create(&drive.File{
			Name:    path,
			Parents: []string{folderId},
		}).Media(f).Do()
		if err != nil {
			if cmd.Verbose {
				log.Printf("Error uploading file %s: %v", path, err)
			}
			return nil, err
		}
	}
	return file, nil
}
