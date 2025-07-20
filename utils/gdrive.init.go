package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"google.golang.org/api/drive/v3"
)

type Drive struct {
	OldInitData map[string]interface{}
	Version     int
	FolderId    string
	Service     *drive.Service
	NewInitData map[string]interface{}
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
	err = json.NewDecoder(file).Decode(&d.OldInitData)
	if err != nil {
		return err
	}
	if d.OldInitData == nil {
		return err
	}
	ver, ok := d.OldInitData["version"]
	if !ok {
		return fmt.Errorf("version not found in init file")
	}
	version, err := strconv.Atoi(fmt.Sprintf("%v", ver))
	if err != nil {
		return err
	}
	d.Version = version
	return nil
}
