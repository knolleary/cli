package user_test

import (
	"cf"
	. "cf/commands/user"
	"cf/configuration"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
)

func TestDeleteUserFailsWithUsage(t *testing.T) {
	userRepo := &testapi.FakeUserRepository{}
	reqFactory := &testreq.FakeReqFactory{}

	ui := callDeleteUser(t, []string{}, userRepo, reqFactory)
	assert.True(t, ui.FailedWithUsage)

	ui = callDeleteUser(t, []string{"foo"}, userRepo, reqFactory)
	assert.False(t, ui.FailedWithUsage)

	ui = callDeleteUser(t, []string{"foo", "bar"}, userRepo, reqFactory)
	assert.True(t, ui.FailedWithUsage)
}

func TestDeleteUserRequirements(t *testing.T) {
	userRepo := &testapi.FakeUserRepository{}
	reqFactory := &testreq.FakeReqFactory{}
	args := []string{"-f", "my-user"}

	reqFactory.LoginSuccess = false
	callDeleteUser(t, args, userRepo, reqFactory)
	assert.False(t, testcmd.CommandDidPassRequirements)

	reqFactory.LoginSuccess = true
	callDeleteUser(t, args, userRepo, reqFactory)
	assert.True(t, testcmd.CommandDidPassRequirements)
}

func TestDeleteUserWhenConfirmingWithY(t *testing.T) {
	ui, userRepo := deleteWithConfirmation(t, "Y")

	assert.Equal(t, len(ui.Outputs), 2)
	assert.Equal(t, len(ui.Prompts), 1)
	assert.Contains(t, ui.Prompts[0], "Really delete")
	assert.Contains(t, ui.Outputs[0], "Deleting user")
	assert.Contains(t, ui.Outputs[0], "my-user")
	assert.Contains(t, ui.Outputs[0], "current-user")

	assert.Equal(t, userRepo.FindByUsernameUsername, "my-user")
	assert.Equal(t, userRepo.DeleteUserGuid, "my-found-user-guid")

	assert.Contains(t, ui.Outputs[1], "OK")
}

func TestDeleteUserWhenConfirmingWithYes(t *testing.T) {
	ui, userRepo := deleteWithConfirmation(t, "Yes")

	assert.Equal(t, len(ui.Outputs), 2)
	assert.Equal(t, len(ui.Prompts), 1)
	assert.Contains(t, ui.Prompts[0], "Really delete")
	assert.Contains(t, ui.Outputs[0], "Deleting user")
	assert.Contains(t, ui.Outputs[0], "my-user")
	assert.Contains(t, ui.Outputs[0], "current-user")

	assert.Equal(t, userRepo.FindByUsernameUsername, "my-user")
	assert.Equal(t, userRepo.DeleteUserGuid, "my-found-user-guid")

	assert.Contains(t, ui.Outputs[1], "OK")
}

func TestDeleteUserWhenNotConfirming(t *testing.T) {
	ui, userRepo := deleteWithConfirmation(t, "Nope")

	assert.Equal(t, len(ui.Outputs), 0)
	assert.Contains(t, ui.Prompts[0], "Really delete")

	assert.Equal(t, userRepo.FindByUsernameUsername, "")
	assert.Equal(t, userRepo.DeleteUserGuid, "")
}

func TestDeleteUserWithForceOption(t *testing.T) {
	foundUser := cf.User{}
	foundUser.Guid = "my-found-user-guid"
	userRepo := &testapi.FakeUserRepository{FindByUsernameUser: foundUser}
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}

	ui := callDeleteUser(t, []string{"-f", "my-user"}, userRepo, reqFactory)

	assert.Equal(t, len(ui.Outputs), 2)
	assert.Equal(t, len(ui.Prompts), 0)
	assert.Contains(t, ui.Outputs[0], "Deleting user")
	assert.Contains(t, ui.Outputs[0], "my-user")

	assert.Equal(t, userRepo.FindByUsernameUsername, "my-user")
	assert.Equal(t, userRepo.DeleteUserGuid, "my-found-user-guid")

	assert.Contains(t, ui.Outputs[1], "OK")
}

func TestDeleteUserWhenUserNotFound(t *testing.T) {
	userRepo := &testapi.FakeUserRepository{FindByUsernameNotFound: true}
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}

	ui := callDeleteUser(t, []string{"-f", "my-user"}, userRepo, reqFactory)

	assert.Equal(t, len(ui.Outputs), 3)
	assert.Equal(t, len(ui.Prompts), 0)
	assert.Contains(t, ui.Outputs[0], "Deleting user")
	assert.Contains(t, ui.Outputs[0], "my-user")

	assert.Equal(t, userRepo.FindByUsernameUsername, "my-user")
	assert.Equal(t, userRepo.DeleteUserGuid, "")

	assert.Contains(t, ui.Outputs[1], "OK")
	assert.Contains(t, ui.Outputs[2], "does not exist")
}

func callDeleteUser(t *testing.T, args []string, userRepo *testapi.FakeUserRepository, reqFactory *testreq.FakeReqFactory) (ui *testterm.FakeUI) {
	ui = new(testterm.FakeUI)

	token, err := testconfig.CreateAccessTokenWithTokenInfo(configuration.TokenInfo{
		Username: "current-user",
	})
	assert.NoError(t, err)
	org_Auto := cf.OrganizationFields{}
	org_Auto.Name = "my-org"
	space_Auto := cf.SpaceFields{}
	space_Auto.Name = "my-space"
	config := &configuration.Configuration{
		Space:        space_Auto,
		Organization: org_Auto,
		AccessToken:  token,
	}

	cmd := NewDeleteUser(ui, config, userRepo)
	ctxt := testcmd.NewContext("delete-user", args)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}

func deleteWithConfirmation(t *testing.T, confirmation string) (ui *testterm.FakeUI, userRepo *testapi.FakeUserRepository) {
	ui = &testterm.FakeUI{
		Inputs: []string{confirmation},
	}
	user_Auto2 := cf.User{}
	user_Auto2.Username = "my-found-user"
	user_Auto2.Guid = "my-found-user-guid"
	userRepo = &testapi.FakeUserRepository{
		FindByUsernameUser: user_Auto2,
	}

	token, err := testconfig.CreateAccessTokenWithTokenInfo(configuration.TokenInfo{
		Username: "current-user",
	})
	assert.NoError(t, err)
	org_Auto2 := cf.OrganizationFields{}
	org_Auto2.Name = "my-org"
	space_Auto2 := cf.SpaceFields{}
	space_Auto2.Name = "my-space"
	config := &configuration.Configuration{
		Space:        space_Auto2,
		Organization: org_Auto2,
		AccessToken:  token,
	}

	cmd := NewDeleteUser(ui, config, userRepo)

	ctxt := testcmd.NewContext("delete-user", []string{"my-user"})
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}

	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
