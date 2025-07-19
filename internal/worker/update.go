package worker

import (
	"encoding/json"
	"fmt"
	"gdriveSync/cmd"
	"gdriveSync/utils"
	"log"
	"os"
	"strconv"

	"google.golang.org/api/drive/v3"
)

type UploadService struct {
	Service *drive.Service
	Data    map[string]interface{}
}
type fileData struct {
	Hash   string `json:"hash"`
	FileId string `json:"fileId"`
}

func (p *UploadService) Upload(path string, version int, hash string, isChanged bool, fileId string, folderId string) error {
	ver := strconv.Itoa(version + 1)
	var file *drive.File
	if isChanged {
		f, err := os.Open(path)
		if err != nil {
			if cmd.Verbose {
				log.Printf("Error opening file %s for upload: %v", path, err)
			}
			return err
		}
		file, err = p.Service.Files.Create(&drive.File{
			Name:    path,
			Parents: []string{folderId},
		}).Media(f).Do()
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
	versionMap := utils.IfExistElseCreate(p.Data, ver)
	versionMap[path] = &fileData{
		Hash:    hash,
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
	tempData := make(map[string]interface{})
	err = json.NewDecoder(f).Decode(&tempData)
	if err != nil {
		return err
	}
	version, ok := tempData["version"].(string)
	if !ok {
		if cmd.Verbose {
			log.Println("Version not found in init file, initializing to 0")
		}
		version = "0"
	}
	ver, err := strconv.Atoi(version)
	if err != nil {
		return err
	}
	verString := strconv.Itoa(ver + 1)
	tempData["version"] = verString
	for k, v := range p.Data {
		if _, exists := tempData[k]; !exists {
			tempData[k] = v
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
	err = json.NewEncoder(f).Encode(tempData)
	if err != nil {
		return err
	}
	return nil
}
