package action

import (
	"fmt"
	"gdriveSync/cmd"
	"gdriveSync/utils"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"google.golang.org/api/drive/v3"
)

type Chann struct {
	Dir       string
	Version   int
	Hash      string
	IsChanged bool
	FileId    string
	Name      string
	FolderId  string
}

type FolderStack struct {
	FolderId  string
	Directory string
}

type PushAction struct {
	InitData *utils.Drive
	PathChan chan *Chann
	Stack    *utils.Stack[FolderStack]
}

func (p *PushAction) Action() error {
	err := p.InitData.ExtractInitData()
	if err != nil {
		return fmt.Errorf("failed to extract init data: %w", err)
	}
	version := strconv.Itoa(p.InitData.Version + 1)
	_, ok := utils.GetMap(p.InitData.NewInitData, version)
	if !ok {
		p.InitData.NewInitData[version] = make(map[string]interface{})
	}
	p.WatchDirectory(version)
	return nil
}

func (p *PushAction) WatchDirectory(version string) {
	defer func() {
		recover()
		close(p.PathChan)
	}()
	var (
		pathFileId string
		pathHash   string
	)
	if cmd.Verbose {
		log.Println("Verbose mode enabled. Watching directory for changes...")
	}
	wd, err := os.Getwd()
	if err != nil {
		if cmd.Verbose {
			log.Printf("Error getting current directory: %v", err)
		}
		return
	}
	filepath.WalkDir(wd, func(fullPath string, d os.DirEntry, err error) error {
		if err != nil {
			if cmd.Verbose {
				log.Printf("Error walking directory: %v", err)
			}
			return err
		}
		folder, ok := p.Stack.Peek()
		folderId := ""
		if ok && folder != nil {
			folderId = folder.FolderId
		}


		if d.IsDir() {
			parentDir := filepath.Dir(fullPath)
			p.PopIfDirectoryDiffer(parentDir)
			p.FolderCreate(fullPath, version, wd, d)
			return nil
		}

		directory := filepath.Dir(fullPath)
		_, flag := p.PopIfDirectoryDiffer(directory)
		if !flag {
			if cmd.Verbose {
				log.Printf("Directory %s not found in stack, skipping file %s", directory, fullPath)
			}
		}
		folder, ok = p.Stack.Peek()
		if ok && folder != nil {
			folderId = folder.FolderId
		}

		hash, err := utils.CreateHash(fullPath)
		if err != nil {
			return err
		}
		if cmd.Verbose {
			log.Printf("File: %s, Hash: %s", fullPath, hash)
		}
		dir, file := path.Split(fullPath)
		ver := strconv.Itoa(p.InitData.Version)
		oldMap, exists := utils.GetCurrentDirectoryMap(p.InitData.OldInitData, ver, wd, dir)
		oldMap = utils.IfExistElseCreate(oldMap, file)
		if exists {
			pathHash, _ = oldMap["hash"].(string)
			pathFileId, _ = oldMap["fileId"].(string)
		}
		p.PathChan <- &Chann{
			Dir:       dir,
			Name:      d.Name(),
			Version:   p.InitData.Version + 1,
			Hash:      hash,
			IsChanged: !exists || pathHash != hash,
			FileId:    pathFileId,
			FolderId:  folderId,
		}
		return nil
	})
}

func (p *PushAction) PopIfDirectoryDiffer(directory string) (*FolderStack, bool) {
	if p.Stack.IsEmpty() {
		return nil, false
	}
	for {
		item, ok := p.Stack.Peek()
		if !ok {
			return nil, false
		}
		if item.Directory == directory {
			if cmd.Verbose {
				log.Printf("Directory %s matches stack top, not popping.", directory)
			}
			return item, true
		}
		if item.Directory != directory {
			_, ok = p.Stack.Pop()
			if !ok {
				return nil, false
			}
		}
		if cmd.Verbose {
			log.Printf("Popped item: %s, Directory: %s", item.Directory, directory)
		}
	}
}

func (p *PushAction) FolderCreate(fullPath string, version string, wd string, d os.DirEntry) error {

	folderId := ""
	folder, ok := p.Stack.Peek()
	if ok {
		folderId = folder.FolderId
	}

	currentVersion := strconv.Itoa(p.InitData.Version)
	oldMap, exists := utils.GetCurrentDirectoryMap(p.InitData.OldInitData, currentVersion, wd, fullPath)

	var existingFolderId string
	if exists {
		if folderIdVal, ok := oldMap["folderId"]; ok {
			existingFolderId, _ = folderIdVal.(string)
		}
	}

	var file *drive.File
	var err error

	if existingFolderId != "" {
		if cmd.Verbose {
			log.Printf("Folder %s already exists with ID: %s", d.Name(), existingFolderId)
		}
		file = &drive.File{
			Id:   existingFolderId,
			Name: d.Name(),
		}
	} else {
		file, err = utils.CreateFolder(p.InitData.Service, d.Name(), folderId)
		if err != nil {
			if cmd.Verbose {
				log.Printf("Error creating folder %s: %v", fullPath, err)
			}
			return err
		}
		fmt.Printf("📂 %s%sFolder created: %s%s%s\n",
			utils.Green, utils.Bold, utils.Cyan, file.Name, utils.Reset)
	}

	if wd == fullPath {
		newMap := utils.IfExistElseCreate(p.InitData.NewInitData, version)
		val := utils.IfExistElseCreate(newMap, fullPath)
		if val["folderId"] == nil {
			val["folderId"] = file.Id
		}
	} else {
		newMap, ok := utils.GetCurrentDirectoryMap(p.InitData.NewInitData, version, wd, fullPath)
		if !ok {
			return fmt.Errorf("failed to get current directory map for %s", fullPath)
		}
		if newMap["folderId"] == nil {
			newMap["folderId"] = file.Id
		}
	}
	p.Stack.Push(FolderStack{
		FolderId:  file.Id,
		Directory: fullPath,
	})
	return nil
}
