package server_repo

import (
	"log"
	"sync"

	logic "github.com/azazel3ooo/keeper/internal/logic/server"
	"github.com/azazel3ooo/keeper/internal/models"
)

// ProcessingWatcher итерируется по каналу сервера. Вызывает необходимые методы хранилища (обновление, удаление...)
// необходим для асинхронной обработки запросов и уменьшении времени ответа
func (s Server) ProcessingWatcher(wt *sync.WaitGroup) {
	defer wt.Done()

	for el := range s.processingChan {
		switch el.Operation {
		case ProcessingOperations[SetOperation]:
			r, ok := el.Data.(models.UserData)
			if !ok {
				log.Println("invalid type in ProcessingWatcher with:", el)
				continue
			}

			err := logic.Set(r, s.storage, el.User)
			if err != nil {
				log.Println(err)
			}

		case ProcessingOperations[DeleteOperation]:
			r, ok := el.Data.(models.DeleteRequest)
			if !ok {
				log.Println("invalid type in ProcessingWatcher with:", el)
				continue
			}

			err := logic.Delete(r, s.storage, el.User)
			if err != nil {
				log.Println(err)
			}

		case ProcessingOperations[UpdateOperation]:
			r, ok := el.Data.(models.UserData)
			if !ok {
				log.Println("invalid type in ProcessingWatcher with:", el)
				continue
			}

			err := logic.Update(r, s.storage, el.User)
			if err != nil {
				log.Println(err)
			}

		default:
			log.Println("unknown type in ProcessingWatcher with ", el)
		}
	}
}
