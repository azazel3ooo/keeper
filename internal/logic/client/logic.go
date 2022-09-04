package client_logic

import (
	"github.com/azazel3ooo/keeper/internal/models"
	"github.com/azazel3ooo/keeper/internal/models/client_repo"
)

// ActionProcessing запускает переданный action локально, после чего, отправляет необходимое на сервер
func ActionProcessing(req models.Validatable, c client_repo.Client, addr, method string,
	action func(c client_repo.Client, r models.Validatable) error) error {

	err := action(c, req)
	if err != nil {
		return err
	}

	err = c.ActionToServer(req, addr, method)
	return err
}

// Set выполняет SetLocal для переданного client_repo.Client с аргументом models.Validatable
// реализовано для передачи в ActionProcessing в качестве обертки над реальным действием
func Set(c client_repo.Client, r models.Validatable) error { return c.SetLocal(r) }

// Update выполняет UpdateLocal для переданного client_repo.Client с аргументом models.Validatable
// реализовано для передачи в ActionProcessing в качестве обертки над реальным действием
func Update(c client_repo.Client, r models.Validatable) error { return c.UpdateLocal(r) }

// Delete выполняет DeleteLocal для переданного client_repo.Client с аргументом models.Validatable
// реализовано для передачи в ActionProcessing в качестве обертки над реальным действием
func Delete(c client_repo.Client, r models.Validatable) error { return c.DeleteLocal(r) }
