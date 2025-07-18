package main

import (
	"gdriveSync/cmd"
	"gdriveSync/internal/action"
	"gdriveSync/internal/oauth"
	"gdriveSync/internal/worker"
	"log"
	"sync"
)

func main() {
	var act action.Actions
	pathChannel := make(chan string, 3)
	numOfWorkers := 3
	wg := &sync.WaitGroup{}
	cmd.Execute()

	if cmd.Create {
		act = &action.InitAction{}
		log.Println("Creating GDrive init file...")
		err := act.Action()
		if err != nil {
			log.Printf("Error executing init action: %v", err)
			return
		}
		log.Println("Initialization complete.")
		return
	}

	if cmd.Push {
		log.Println("Starting push action...")
		act = &action.PushAction{
			Data:     make(map[string]interface{}),
			PathChan: pathChannel,
		}
		err := act.Action()
		if err != nil {
			log.Printf("Error executing push action: %v", err)
			return
		}
		log.Println("Push action completed successfully.")
	}
	log.Println("Starting upload workers...")

	service, err := oauth.DriveClient()
	if err != nil {
		log.Printf("Error creating Drive client: %v", err)
		return
	}
	uploadService := &worker.UploadService{DriveService: service}
	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChannel {
				err := uploadService.Upload(path)
				if err != nil {
					log.Printf("Error uploading file %s: %v", path, err)
				}
			}
		}()
	}
	wg.Wait()
}
