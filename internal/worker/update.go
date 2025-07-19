package worker

import (
	"encoding/json"
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
		log.Printf("Error opening file for upload: %v", err)
	}
	file, err := p.DriveService.Files.Create(&drive.File{
		Name: path,
	}).Media(f).Do()
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return err
	}
	log.Printf("File uploaded successfully: %s", file.Id)
	versionMap, ok := p.Data[ver].(map[string]interface{})
	if !ok {
		versionMap = make(map[string]interface{})
		p.Data[ver] = versionMap
		log.Printf("Adding new version %s to init file", ver)
	}
	versionMap[path] = hash
	return nil
}

func (p *UploadService) UpdateWithVersion() error {
	f, err := os.OpenFile("gdrive_init.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error opening init file for update: %v", err)
		return err
	}
	defer f.Close()
	tempData := make(map[string]interface{})
	err = json.NewDecoder(f).Decode(&tempData)
	if err != nil {
		log.Printf("Error decoding init file: %v", err)
		return err
	}
	version, ok := tempData["version"].(string)
	if !ok {
		log.Println("Version not found in init file, initializing to 0")
		version = "0"
	}
	ver, err := strconv.Atoi(version)
	if err != nil {
		log.Printf("Error converting version to int: %v", err)
		return err
	}
	verString := strconv.Itoa(ver + 1)
	tempData["version"] = verString
	log.Printf("Incremented version to: %s", verString)
	for k, v := range p.Data {
		if _, exists := tempData[k]; !exists {
			tempData[k] = v
			log.Printf("Adding new version %s to init file", k)
		} else {
			log.Printf("Version %s already exists in init file, skipping", k)
		}
	}
	if _, err = f.Seek(0, 0); err != nil {
		log.Printf("Error seeking to start of init file: %v", err)
		return err
	}
	if err = f.Truncate(0); err != nil {
		log.Printf("Error truncating init file: %v", err)
		return err
	}
	err = json.NewEncoder(f).Encode(tempData)
	if err != nil {
		log.Printf("Error encoding data to init file: %v", err)
		return err
	}
	return nil
}
