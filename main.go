package main

import (
	"fmt"
	"gdriveSync/cmd"
	"gdriveSync/internal/action"
	"gdriveSync/internal/oauth"
	"gdriveSync/internal/worker"
	"gdriveSync/utils"
	"log"
	"os"
	"sync"

	"google.golang.org/api/drive/v3"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		os.Exit(1)
	}
	tokDirectory := home + "/.gdrive"
	if _, err := os.Stat(tokDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(tokDirectory, 0700)
		if err != nil {
			os.Exit(1)
		}
		if cmd.Verbose {
			log.Printf("Token directory created at: %s", tokDirectory)
		}
	}

	var act action.Actions
	pathChannel := make(chan *action.Chann, 3)
	wg := &sync.WaitGroup{}
	cmd.Execute()
	service, err := oauth.DriveClient()
	if err != nil {
		fmt.Printf("🔒 %s%sAUTH FAILED:%s Cannot create Drive client: %s%s%v%s\n",
			utils.Red, utils.Bold, utils.Reset, utils.Red, utils.Bold, err, utils.Reset)
		fmt.Printf("Please ensure you have a valid credentials file at %s%s%s/.gdrive/credentials.json%s\n",
			utils.Yellow, utils.Bold, os.Getenv("HOME"), utils.Reset)
		os.Exit(1)
	}
	oldMap := make(map[string]interface{})
	newMap := make(map[string]interface{})
	initData := &utils.Drive{
		OldInitData: oldMap,
		Version:     0,
		FolderId:    "",
		Service:     service,
		NewInitData: newMap,
	}
	stack := &utils.Stack[action.FolderStack]{}

	if cmd.Create {
		act = &action.InitAction{
			Service: service,
		}
		err := act.Action()
		if err != nil {
			fmt.Printf("❌ %s%sERROR:%s Init action failed: %s%s%v%s",
				utils.Red, utils.Bold, utils.Reset, utils.Red, utils.Bold, err, utils.Reset)
			os.Exit(1)
		}
		fmt.Printf("✅ %s%sInitialization complete.%s", utils.Green, utils.Bold, utils.Reset)
		os.Exit(1)
	}

	if cmd.Push {
		if cmd.Verbose {
			log.Println("Push action initiated.")
		}
		go func() {
			act = &action.PushAction{
				InitData: initData,
				PathChan: pathChannel,
				Stack:    stack,
			}
			err := act.Action()
			if err != nil {
				fmt.Printf("💥 %s%sFATAL:%s Push action crashed: %s%s%v%s",
					utils.Red, utils.Bold, utils.Reset, utils.Red, utils.Bold, err, utils.Reset)
				os.Exit(1)
			}
		}()
		if cmd.Verbose {
			log.Println("Push action completed.")
		}
		uploadService := &worker.UploadService{Service: service, NewInitData: newMap, OldInitData: oldMap, Stack: stack}
		if cmd.Verbose {
			log.Println("Starting upload workers...")
		}

		for i := 0; i < cmd.NumOfWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for path := range pathChannel {
					err := uploadService.UploadFile(path.Dir, path.Name, path.Version, path.Hash, path.IsChanged, path.FileId, path.FolderId)
					if err != nil {
						fmt.Printf("▶ %s%sFAILED:%s File %s%s%s could not be uploaded: %s%s%v%s",
							utils.Yellow, utils.Bold, utils.Reset, utils.Cyan, path.Dir, utils.Reset, utils.Red, utils.Bold, err, utils.Reset)
						continue
					}
				}
			}()
		}
		wg.Wait()
		fmt.Printf("✅ %s%sUpload completed.%s", utils.Green, utils.Bold, utils.Reset)
		err = uploadService.UpdateWithVersion()
		if err != nil {
			fmt.Printf("⚠️  %s%sVERSION UPDATE FAILED:%s Could not update init file: %s%s%v%s",
				utils.Red, utils.Bold, utils.Reset, utils.Red, utils.Bold, err, utils.Reset)
			os.Exit(1)
		}
	}

	if cmd.Show {
		if cmd.Verbose {
			log.Println("List action initiated.")
		}
		act = &action.ListAction{
			DriveService:  service,
			FolderHandler: make(map[string][]string),
			FileName:      make(map[string]*drive.File),
			Shared:        []string{},
			SharedFiles:   []string{},
			SharedFolders: []string{},
		}
		err := act.Action()
		if err != nil {
			fmt.Printf("❌ %s%sERROR:%s List action failed: %s%s%v%s",
				utils.Red, utils.Bold, utils.Reset, utils.Red, utils.Bold, err, utils.Reset)
			os.Exit(1)
		}
		if cmd.Verbose {
			log.Println("List action completed.")
		}
	}
}
