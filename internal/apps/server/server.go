package server

import (
	"log"
	"sync"

	"github.com/azazel3ooo/keeper/internal/models"
	repo "github.com/azazel3ooo/keeper/internal/models/server_repo"
)

// Start производит инициализацию сервера и запускает его на порту, указанном в конфигурационном файле
func Start() {
	var cfg models.Config
	err := cfg.Init("server_settings.yml")
	if err != nil {
		log.Fatal(err)
	}

	var storage repo.ServerStorage
	err = storage.Init(cfg.DbLocation)
	if err != nil {
		log.Fatal(err)
	}

	processingChan := make(repo.ProcessingChan, 1000)
	s := repo.NewServer(
		repo.WithConfig(cfg),
		repo.WithProcessingChan(processingChan),
		repo.WithStorage(&storage))

	s.SetupApp()

	var watcherWG sync.WaitGroup
	watcherWG.Add(1)
	go s.ProcessingWatcher(&watcherWG)

	log.Println(s.Listen())

	close(processingChan)
	watcherWG.Wait()
}
