package client

import (
	"errors"
	"fmt"
	logic "github.com/azazel3ooo/keeper/internal/logic/client"
	"github.com/azazel3ooo/keeper/internal/models"
	repo "github.com/azazel3ooo/keeper/internal/models/client_repo"
	"log"
	"net/http"
)

const BuildVersion = "v0.0.1"
const BuildDate = "20.08.22"

// Start производит инициализацию клиента, после чего запускает интерактивное меню
func Start() {
	var cfg models.Config

	err := cfg.Init("client_settings.yml")
	if err != nil {
		log.Fatal(err)
	}

	var s repo.ClientStorage
	err = s.Init(cfg.DbLocation)
	if err != nil {
		log.Fatal(err)
	}

	c := repo.NewClient(
		repo.WithStorage(&s),
		repo.WithConfig(cfg),
		repo.WithClient(http.DefaultClient),
	)

	LoopMenu(*c)
}

// LoopMenu проводит авторизацию\регистрацию пользователя, после чего запускает зацикленное меню для выполнения действий
// (добавление, удаление, обновление и получение полного списка записей)
func LoopMenu(c repo.Client) {
	// GET TOKEN STAGE
	fmt.Println("Hello in our keeper\nChoose option")
	fmt.Printf("Reqistration: type option \"r\"\n" +
		"Authorization: type option \"a\"\n" +
		"Build version: type option \"b\"\n" +
		"Build date: type option \"d\"\n")

	for !c.ReadyForActions() {
		var option string
		fmt.Printf("Type option:\n")
		_, err := fmt.Scanf("%s\n", &option)
		if err != nil {
			fmt.Printf("Please try again, error: %s\n", err.Error())
			continue
		}

		if option == "b" {
			fmt.Println(BuildVersion)
			continue
		}
		if option == "d" {
			fmt.Println(BuildDate)
			continue
		}

		var requestBody models.UserRequest
		fmt.Printf("Type your login:\n")
		_, err = fmt.Scanf("%s\n", &requestBody.Login)
		if err != nil {
			fmt.Printf("Please try again, error: %s\n", err.Error())
			continue
		}

		fmt.Printf("Type your password:\n")
		_, err = fmt.Scanf("%s\n", &requestBody.Password)
		if err != nil {
			fmt.Printf("Please try again, error: %s\n", err.Error())
			continue
		}

		var reqAddr string
		if option == "r" {
			reqAddr = c.RegistrationAddress()
		}
		if option == "a" {
			reqAddr = c.AuthorizationAddress()
		}
		if reqAddr == "" {
			fmt.Printf("Unknown option. Please, try again\n")
			continue
		}

		err = c.GetToken(requestBody, reqAddr)
		if err != nil {
			fmt.Printf("Please try again, error: %s\n", err.Error())
			continue
		}

		if option == "a" {
			err = c.ActualizeStorage()
			if err != nil {
				log.Println("can't actualize storage: " + err.Error())
			}
		}

	}
	// GET TOKEN STAGE

	// GET OPTION STAGE
	fmt.Printf("Choose action:\n" +
		"Add: type a\n" +
		"Update: type u\n" +
		"Delete: type d\n" +
		"Get data list: type g\n" +
		"Quit: type q\n")

	finished := false
	for !finished {
		var action string
		fmt.Printf("Type action:\n")
		_, err := fmt.Scanf("%s\n", &action)
		if err != nil {
			fmt.Printf("Please try again, error: %s\n", err.Error())
			continue
		}

		switch action {
		case "a":
			var req models.UserData
			fmt.Printf("Type your data:\n")
			_, err = fmt.Scanf("%s\n", &req.Data)
			if err != nil {
				fmt.Printf("Please try again, error: %s\n", err.Error())
				continue
			}

			fmt.Printf("Type your metadata:\n")
			_, err = fmt.Scanf("%s\n", &req.Comment)
			if err != nil {
				fmt.Printf("Please try again, error: %s\n", err.Error())
				continue
			}

			req.ID = logic.GenerateID()
			err = logic.ActionProcessing(req, c, c.ActionAddr(), http.MethodPost, logic.Set)
			if errors.Is(err, models.ErrExpiredToken) {
				finished = true
			}
			if err != nil {
				fmt.Printf("Please try again, error: %s\n", err.Error())
				continue
			}

		case "u":
			var req models.UserData
			fmt.Printf("Type data ID:\n")
			_, err = fmt.Scanf("%s\n", &req.ID)
			if err != nil {
				fmt.Printf("Please try again, error: %s\n", err.Error())
				continue
			}

			fmt.Printf("Type your data:\n")
			_, err = fmt.Scanf("%s\n", &req.Data)
			if err != nil {
				fmt.Printf("Please try again, error: %s\n", err.Error())
				continue
			}

			fmt.Printf("Type your metadata:\n")
			_, err = fmt.Scanf("%s\n", &req.Comment)
			if err != nil {
				fmt.Printf("Please try again, error: %s\n", err.Error())
				continue
			}

			err = logic.ActionProcessing(req, c, c.ActionAddr(), http.MethodPatch, logic.Update)
			if errors.Is(err, models.ErrExpiredToken) {
				finished = true
			}
			if err != nil {
				fmt.Printf("Please try again, error: %s\n", err.Error())
				continue
			}

		case "d":
			var req models.DeleteRequest
			fmt.Printf("Type ID:\n")
			_, err = fmt.Scanf("%s\n", &req.ID)
			if err != nil {
				fmt.Printf("Please try again, error: %s\n", err.Error())
				continue
			}

			err = logic.ActionProcessing(req, c, c.ActionAddr(), http.MethodDelete, logic.Delete)
			if errors.Is(err, models.ErrExpiredToken) {
				finished = true
			}
			if err != nil {
				fmt.Printf("Please try again, error: %s\n", err.Error())
				continue
			}

		case "g":
			data, err := c.GetAll()
			if err != nil {
				fmt.Printf("Please try again, error: %s\n", err.Error())
				continue
			}
			logic.PrintData(data)

		case "q":
			finished = true

		default:
			fmt.Printf("Unknown option. Please, try again\n")
		}
	}
	// GET OPTION STAGE

	fmt.Println("Bye!")
}
