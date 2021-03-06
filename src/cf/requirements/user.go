package requirements

import (
	"cf/api"
	"cf/errors"
	"cf/models"
	"cf/terminal"
)

type UserRequirement interface {
	Requirement
	GetUser() models.UserFields
}

type userApiRequirement struct {
	username string
	ui       terminal.UI
	userRepo api.UserRepository
	user     models.UserFields
}

func NewUserRequirement(username string, ui terminal.UI, userRepo api.UserRepository) (req *userApiRequirement) {
	req = new(userApiRequirement)
	req.username = username
	req.ui = ui
	req.userRepo = userRepo
	return
}

func (req *userApiRequirement) Execute() (success bool) {
	var apiErr errors.Error
	req.user, apiErr = req.userRepo.FindByUsername(req.username)

	if apiErr != nil {
		req.ui.Failed(apiErr.Error())
		return false
	}

	return true
}

func (req *userApiRequirement) GetUser() models.UserFields {
	return req.user
}
