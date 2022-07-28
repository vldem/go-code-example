package update

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/vldem/go-code-example/telbot_v2/internal/auth"
	commandPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/bot/command"
	userPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user"
	"github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/models"
	validatorPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/validator"
)

var msgUpdateUser = "update user"

type command struct {
	user userPkg.Interface
}

func New(user userPkg.Interface) commandPkg.Interface {
	return &command{
		user: user,
	}
}

func (c *command) Name() string {
	return "update"
}

func (c *command) Description() string {
	return "<id>;<email>;<name>;<role>;<password>;<old password> - update user"
}

func (c *command) Process(args string) string {
	params := strings.Split(args, ";")
	if len(params) != 6 {
		return commandPkg.MsgInvalidArguments
	}

	//validate parameters
	err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate(params[1 : len(params)-1]))
	if err != nil {
		return errors.Wrap(err, msgUpdateUser).Error()
	}

	//get user by id
	id, _ := strconv.ParseUint(params[0], 10, 64)
	user, err := c.user.Get(uint(id))
	if err != nil {
		return errors.Wrap(err, msgUpdateUser).Error()
	}

	// verify password to confirm changes
	err = auth.VerifyPassword(*user, params[5])
	if err != nil {
		return errors.Wrap(err, msgUpdateUser).Error()
	}

	user.Email = params[1]
	user.Name = params[2]
	user.Role = models.GetRoleId(params[3])
	user.Password = auth.GenHashPassword(params[4])

	if err := c.user.Update(*user); err != nil {
		return errors.Wrap(err, msgUpdateUser).Error()
	}

	return "user has been updated"
}
