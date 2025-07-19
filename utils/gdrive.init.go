package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
)

type Drive struct {
	Data    map[string]interface{}
	Version int
	FolderId string
}

func (d *Drive) ExtractInitData() error {
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
	err = json.NewDecoder(file).Decode(&d.Data)
	if err != nil {
		return err
	}
	if d.Data == nil {
		return err
	}
	ver, ok := d.Data["version"]
	if !ok {
		return fmt.Errorf("version not found in init file")
	}
	version, err := strconv.Atoi(fmt.Sprintf("%v", ver))
	if err != nil {
		return err
	}
	d.Version = version

	folderId, ok := d.Data["folderId"]
	if !ok {
		return fmt.Errorf("folderId not found in init file")
	}
	d.FolderId = folderId.(string)
	return nil
}
