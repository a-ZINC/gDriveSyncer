package main

import (
	"gdriveSync/cmd"
	"gdriveSync/internal/action"
	"gdriveSync/internal/oauth"
	"gdriveSync/internal/worker"
	"log"
	"os"
	"sync"
)

func main() {
	var act action.Actions
	pathChannel := make(chan *action.Chann, 3)
	wg := &sync.WaitGroup{}
	cmd.Execute()

	if cmd.Create {
		act = &action.InitAction{}
		err := act.Action()
		if err != nil {
			log.Printf("Error executing init action: %v", err)
			return
		}
		log.Println("Initialization complete.")
		os.Exit(1)
	}

	if cmd.Push {
		log.Println("Starting push action...")
		go func() {
			act = &action.PushAction{
				Data:     make(map[string]interface{}),
				PathChan: pathChannel,
			}
			err := act.Action()
			if err != nil {
				log.Printf("Error executing push action: %v", err)
				return
			}
		}()
		log.Println("Push action completed successfully.")
	}

	service, err := oauth.DriveClient()
	if err != nil {
		log.Printf("Error creating Drive client: %v", err)
		os.Exit(1)
	}
	uploadService := &worker.UploadService{DriveService: service, Data: make(map[string]interface{})}
	if cmd.Verbose {
		log.Println("Starting upload workers...")
	}
	for i := 0; i < cmd.NumOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChannel {
				err := uploadService.Upload(path.Path, path.Version, path.Hash)
				if err != nil {
					log.Printf("Error uploading file %s: %v", path.Path, err)
					continue
				}
				log.Printf("File %s uploaded successfully with version %d", path.Path, path.Version)
			}
		}()
	}
	wg.Wait()
	log.Println("All uploads completed. Updating version in init file...")
	err = uploadService.UpdateWithVersion()
	if err != nil {
		log.Printf("Error updating version in init file: %v", err)
		return
	}
}
