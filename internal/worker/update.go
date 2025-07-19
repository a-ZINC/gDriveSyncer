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
	ver := strconv.Itoa(version + 1)
	if p.Data == nil {
		log.Println("Data map is nil, initializing...")
		p.Data = make(map[string]interface{})
	}
	newMap, ok := p.Data[ver].( map[string]interface{})
	if !ok || newMap == nil {
		log.Printf("Error converting data to map for version %s", ver)
		newMap = make(map[string]interface{})
		p.Data[ver] = newMap
	}
	newMap[path] = hash
	log.Printf("Updated version for %s to %d", path, version)
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
