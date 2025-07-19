package worker

import (
	"encoding/json"
	"gdriveSync/cmd"
	"log"
	"os"
	"strconv"

	"google.golang.org/api/drive/v3"
)

type UploadService struct {
	DriveService *drive.Service
	Data         map[string]interface{}
}

func (p *UploadService) Upload(path string, version int, hash string) error {

	ver := strconv.Itoa(version + 1)
	f, err := os.Open(path)
	if err != nil {
		if cmd.Verbose {
			log.Printf("Error opening file %s for upload: %v", path, err)
		}
		return err
	}
	_, err = p.DriveService.Files.Create(&drive.File{
		Name: path,
	}).Media(f).Do()
	if err != nil {
		if cmd.Verbose {
			log.Printf("Error uploading file %s: %v", path, err)
		}
		return err
	}
	versionMap, ok := p.Data[ver].(map[string]interface{})
	if !ok {
		versionMap = make(map[string]interface{})
		p.Data[ver] = versionMap
	}
	versionMap[path] = hash
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
