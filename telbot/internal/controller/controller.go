// This file contains controllers that call a function for corresponding bot command
package controller

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/vldem/go-code-example/telbot/internal/auth"
	"github.com/vldem/go-code-example/telbot/internal/commander"
	"github.com/vldem/go-code-example/telbot/internal/model"
	"github.com/vldem/go-code-example/telbot/internal/storage"
	"github.com/vldem/go-code-example/telbot/internal/validator"

	"github.com/pkg/errors"
)

const (
	helpCmd   = "help"
	listCmd   = "list"
	addCmd    = "add"
	updateCmd = "update"
	deleteCmd = "delete"
)

var BadArgument = errors.New("bad argument")

func AddControllers(c *commander.Commander) {
	c.RegisterController(helpCmd, helpFunc)
	c.RegisterController(listCmd, listFunc)
	c.RegisterController(addCmd, addFunc)
	c.RegisterController(updateCmd, updateFunc)
	c.RegisterController(deleteCmd, deleteFunc)
}

func listFunc(s string) string {
	data := storage.ListUser()
	res := make([]string, 0, len(data))
	for _, v := range data {
		res = append(res, v.String())
	}
	return strings.Join(res, "\n")
}

func helpFunc(s string) string {
	return "/help - list commands\n" +
		"/list - list users\n" +
		"/add <email>;<name>;<role>;<password> - add new user with email, name, role and password\n" +
		"/update <id>;<email>;<name>;<role>;<new password>;<old password> - update existing user with new email and/or name and/or password. Old password is needed to confirm changes.\n" +
		"/delete <id>;<password> - delete existing with id <id>. Please enter password to confirm deletion\n"
}

func addFunc(data string) string {
	log.Printf("add command param: <%s>", data)
	params := strings.Split(data, ";")
	if len(params) != 4 {
		return errors.Wrapf(BadArgument, "%d items: <%v>", len(params), params).Error()
	}

	//validate parameters
	err := validator.ValidateParameters(makeParametersToValidate(params))
	if err != nil {
		return errors.Wrap(err, "validation error").Error()
	}

	//add new user (email, name, role, password)
	u := model.NewUser(params[0], params[1], params[2], params[3])
	err = storage.AddUser(u)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("user %v added", u)
}

func updateFunc(data string) string {
	log.Printf("update command param: <%s>", data)
	params := strings.Split(data, ";")
	if len(params) != 6 {
		return errors.Wrapf(BadArgument, "%d items: <%v>", len(params), params).Error()
	}
	// verify password to confirm changes
	err := verifyPassword(params[0], params[5])
	if err != nil {
		return errors.Wrap(err, "delete user").Error()
	}

	id, _ := strconv.ParseUint(params[0], 10, 64)

	//validate parameters
	err = validator.ValidateParameters(makeParametersToValidate(params[1 : len(params)-1]))
	if err != nil {
		return errors.Wrap(err, "validation error").Error()
	}

	user, err := storage.GetUser(uint(id))
	if err != nil {
		return errors.Wrap(err, "update user").Error()
	}

	user.SetEmail(params[1])
	user.SetName(params[2])
	user.SetRole(params[3])
	user.SetPwd(params[4])

	err = storage.UpdateUser(user)
	if err != nil {
		return errors.Wrap(err, "update user").Error()
	}

	return fmt.Sprintf("user %v updated", user)

}

func deleteFunc(data string) string {
	log.Printf("delete command param: <%s>", data)
	params := strings.Split(data, ";")
	if len(params) != 2 {
		return errors.Wrapf(BadArgument, "%d items: <%v>", len(params), params).Error()
	}
	// verify password to confirm changes
	err := verifyPassword(params[0], params[1])
	if err != nil {
		return errors.Wrap(err, "delete user").Error()
	}

	id, _ := strconv.ParseUint(params[0], 10, 64)
	err = storage.DeleteUser(uint(id))
	if err != nil {
		return errors.Wrap(err, "delete user").Error()
	}
	return fmt.Sprintf("user %d has been deleted", id)

}

func verifyPassword(userId, pwd string) error {
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return errors.Wrapf(BadArgument, "id: <%v>", userId)
	}
	err = auth.VerifyPassword(uint(id), pwd)
	if err != nil {
		return err
	}
	return nil

}

func makeParametersToValidate(params []string) map[string]string {
	paramList := make(map[string]string, len(params))
	paramList["email"] = params[0]
	paramList["name"] = params[1]
	paramList["role"] = params[2]
	paramList["password"] = params[3]
	return paramList
}
