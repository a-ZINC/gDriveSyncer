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
	defer close(p.PathChan)
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
		folder, _ := p.Stack.Peek()

		if d.IsDir() {
			folderId := ""
			folder, ok := p.Stack.Peek()
			if ok {
				folderId = folder.FolderId
			}
			file, err := utils.CreateFolder(p.InitData.Service, d.Name(), folderId)
			if err != nil {
				if cmd.Verbose {
					log.Printf("Error creating folder %s: %v", fullPath, err)
				}
				return err
			}
			fmt.Printf("📂 %s%sFolder created in InitData: %s%s%s\n",
				utils.Green, utils.Bold, utils.Cyan, file.Name, utils.Reset)
			if wd == fullPath {
				newMap := utils.IfExistElseCreate(p.InitData.NewInitData, version)
				val := utils.IfExistElseCreate(newMap, fullPath)
				if val["folderId"] == nil {
					val["folderId"] = file.Id
				}
			} else {
				newMap, ok := utils.GetCurrentDirectoryMap(p.InitData.NewInitData, version, wd, fullPath)
				if !ok {
					log.Printf("Failed to get current directory map for %s", fullPath)
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
		directory := filepath.Dir(fullPath)
		_, flag := p.PopIfDirectoryDiffer(directory)
		if !flag {
			log.Printf("Directory %s popped from stack.", directory)
		}

		hash, err := utils.CreateHash(fullPath)
		if err != nil {
			return err
		}
		if cmd.Verbose {
			log.Printf("File: %s, Hash: %s", fullPath, hash)
		}
		dir, file := path.Split(fullPath)
		oldMap, exists := utils.GetCurrentDirectoryMap(p.InitData.OldInitData, version, wd, dir)
		if exists {
			pathHash, _ = oldMap[file].(string)
			pathFileId, _ = oldMap[file].(string)
		}
		p.PathChan <- &Chann{
			Dir:       dir,
			Name:      d.Name(),
			Version:   p.InitData.Version + 1,
			Hash:      hash,
			IsChanged: !exists || pathHash != hash,
			FileId:    pathFileId,
			FolderId:  folder.FolderId,
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
