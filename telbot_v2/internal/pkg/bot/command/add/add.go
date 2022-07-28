package add

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/vldem/go-code-example/telbot_v2/internal/auth"
	commandPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/bot/command"
	userPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user"
	"github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/models"
	validatorPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/validator"
)

const msgAddUser = "add user"

type command struct {
	user userPkg.Interface
}

func New(user userPkg.Interface) commandPkg.Interface {
	return &command{
		user: user,
	}
}

func (c *command) Name() string {
	return "add"
}

func (c *command) Description() string {
	return "<email>;<name>;<role>;<password> - create user"
}

func (c *command) Process(args string) string {
	params := strings.Split(args, ";")
	if len(params) != 4 {
		return commandPkg.MsgInvalidArguments
	}

	//validate parameters
	err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate(params))
	if err != nil {
		return errors.Wrap(err, msgAddUser).Error()
	}

	if err := c.user.Create(models.User{
		Id:       0,
		Email:    params[0],
		Name:     params[1],
		Role:     models.GetRoleId(params[2]),
		Password: auth.GenHashPassword(params[3]),
	}); err != nil {
		return errors.Wrap(err, msgAddUser).Error()
	}

	return "user has been added"
}
