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
)

func main() {
	var act action.Actions
	pathChannel := make(chan *action.Chann, 3)
	wg := &sync.WaitGroup{}
	cmd.Execute()
	service, err := oauth.DriveClient()
	if err != nil {
		fmt.Printf("🔒 %s%sAUTH FAILED:%s Cannot create Drive client: %s%s%v%s",
			utils.Red, utils.Bold, utils.Reset, utils.Red, utils.Bold, err, utils.Reset)
		os.Exit(1)
	}
	initData := &utils.Drive{
		Data: make(map[string]interface{}),
		Version: 0,
		FolderId: "",
	}

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
				Drive:     initData,
				PathChan: pathChannel,
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
	}

	uploadService := &worker.UploadService{Service: service, Data: make(map[string]interface{})}
	if cmd.Verbose {
		log.Println("Starting upload workers...")
	}

	for i := 0; i < cmd.NumOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChannel {
				err := uploadService.Upload(path.Path, path.Version, path.Hash, path.IsChanged, path.FileId, initData.FolderId)
				if err != nil {
					fmt.Printf("▶ %s%sFAILED:%s File %s%s%s could not be uploaded: %s%s%v%s",
						utils.Yellow, utils.Bold, utils.Reset, utils.Cyan, path.Path, utils.Reset, utils.Red, utils.Bold, err, utils.Reset)
					continue
				}
			}
		}()
	}
	wg.Wait()
	log.Println("All uploads completed. Updating version in init file...")
	err = uploadService.UpdateWithVersion()
	if err != nil {
		fmt.Printf("⚠️  %s%sVERSION UPDATE FAILED:%s Could not update init file: %s%s%v%s",
			utils.Red, utils.Bold, utils.Reset, utils.Red, utils.Bold, err, utils.Reset)
		os.Exit(1)
	}
}
