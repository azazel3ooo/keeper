package server_repo

import (
	"errors"
	logic "github.com/azazel3ooo/keeper/internal/logic/server"
	"github.com/azazel3ooo/keeper/internal/models"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"

	_ "github.com/azazel3ooo/keeper/docs"
)

// registration godoc
// @Description  handler for registration
// @Tags         All
// @Accept       json
// @Produce      json
// @Param        request body models.UserRequest true "Request structure"
// @Success      200	{object} models.UserResponse
// @Failure      400
// @Failure      409
// @Failure      500
// @Router       /api/v1/registration [post]
func (s *Server) registration(c *fiber.Ctx) error {
	var req models.UserRequest
	err := c.BodyParser(&req)
	if err != nil || !req.Valid() {
		return c.SendStatus(http.StatusBadRequest)
	}

	id, err := logic.Registration(req, s.storage)
	if errors.Is(err, models.ErrUserConflict) {
		return c.SendStatus(http.StatusConflict)
	} else if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	token, err := logic.GenerateToken(id)
	if err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(models.UserResponse{
		Token: token,
	})
}

// authorization godoc
// @Description  handler for authorization
// @Tags         All
// @Accept       json
// @Produce      json
// @Param        request body models.UserRequest true "Request structure"
// @Success      200	{object} models.UserResponse
// @Failure      400
// @Failure      403
// @Failure      500
// @Router       /api/v1/auth [post]
func (s *Server) authorization(c *fiber.Ctx) error {
	var req models.UserRequest
	err := c.BodyParser(&req)
	if err != nil || !req.Valid() {
		return c.SendStatus(http.StatusBadRequest)
	}

	id, err := logic.CheckUser(req, s.storage)
	if errors.Is(err, models.ErrUserDataConflict) {
		return c.SendStatus(http.StatusForbidden)
	} else if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	token, err := logic.GenerateToken(id)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(models.UserResponse{
		Token: token,
	})
}

// getAll godoc
// @Description  handler for get full list of user data
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Insert your access token" default(<Add access token here>)
// @Success      200	{object} models.UserDataResponse
// @Failure      401
// @Failure      403
// @Failure      500
// @Router       /api/v1/items [get]
func (s *Server) getAll(c *fiber.Ctx) error {
	id, err := logic.CheckToken(c.Get("Authorization"))
	if err != nil {
		status := http.StatusInternalServerError // default
		if errors.Is(err, models.ErrInvalidToken) {
			return c.SendStatus(http.StatusForbidden)
		}
		if errors.Is(err, models.ErrExpiredToken) {
			return c.SendStatus(http.StatusUnauthorized)
		}

		return c.SendStatus(status)
	}

	res, err := logic.GetAll(id, s.storage)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(models.UserDataResponse{
		Data: res,
	})
}

// set godoc
// @Description  handler for set new data in global storage
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Insert your access token" default(<Add access token here>)
// @Param        request body models.UserData true "Request structure"
// @Success      200
// @Failure      400
// @Failure      403
// @Failure      401
// @Failure      500
// @Router       /api/v1/items [post]
func (s *Server) set(c *fiber.Ctx) error {
	id, err := logic.CheckToken(c.Get("Authorization"))
	if err != nil {
		status := http.StatusInternalServerError // default
		if errors.Is(err, models.ErrInvalidToken) {
			return c.SendStatus(http.StatusForbidden)
		}
		if errors.Is(err, models.ErrExpiredToken) {
			return c.SendStatus(http.StatusUnauthorized)
		}

		return c.SendStatus(status)
	}

	var req models.UserData
	err = c.BodyParser(&req)
	if err != nil || !req.Valid() {
		return c.SendStatus(http.StatusBadRequest)
	}

	s.processingChan <- ProcessingTuple{Operation: ProcessingOperations[SetOperation], Data: req, User: id}

	return c.SendStatus(http.StatusOK)
}

// delete godoc
// @Description  handler for delete exist data in global storage
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Insert your access token" default(<Add access token here>)
// @Param        request body models.DeleteRequest true "Request structure"
// @Success      200
// @Failure      400
// @Failure      403
// @Failure      401
// @Failure      500
// @Router       /api/v1/items [delete]
func (s *Server) delete(c *fiber.Ctx) error {
	id, err := logic.CheckToken(c.Get("Authorization"))
	if err != nil {
		status := http.StatusInternalServerError // default
		if errors.Is(err, models.ErrInvalidToken) {
			return c.SendStatus(http.StatusForbidden)
		}
		if errors.Is(err, models.ErrExpiredToken) {
			return c.SendStatus(http.StatusUnauthorized)
		}

		return c.SendStatus(status)
	}

	var req models.DeleteRequest
	err = c.BodyParser(&req)
	if err != nil || !req.Valid() {
		return c.SendStatus(http.StatusBadRequest)
	}

	s.processingChan <- ProcessingTuple{Operation: ProcessingOperations[DeleteOperation], Data: req, User: id}

	return c.SendStatus(http.StatusOK)
}

// update godoc
// @Description  handler for update exist data in global storage
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Insert your access token" default(<Add access token here>)
// @Param        request body models.UserData true "Request structure"
// @Success      200
// @Failure      400
// @Failure      403
// @Failure      401
// @Failure      500
// @Router       /api/v1/items [patch]
func (s *Server) update(c *fiber.Ctx) error {
	id, err := logic.CheckToken(c.Get("Authorization"))
	if err != nil {
		status := http.StatusInternalServerError // default
		if errors.Is(err, models.ErrInvalidToken) {
			return c.SendStatus(http.StatusForbidden)
		}
		if errors.Is(err, models.ErrExpiredToken) {
			return c.SendStatus(http.StatusUnauthorized)
		}

		return c.SendStatus(status)
	}

	var req models.UserData
	err = c.BodyParser(&req)
	if err != nil || !req.Valid() {
		return c.SendStatus(http.StatusBadRequest)
	}

	s.processingChan <- ProcessingTuple{Operation: ProcessingOperations[UpdateOperation], Data: req, User: id}

	return c.SendStatus(http.StatusOK)
}
