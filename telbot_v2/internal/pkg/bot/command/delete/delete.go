package delete

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/vldem/go-code-example/telbot_v2/internal/auth"
	commandPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/bot/command"
	userPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user"
)

var msgDeleteUser = "delete user"

type command struct {
	user userPkg.Interface
}

func New(user userPkg.Interface) commandPkg.Interface {
	return &command{
		user: user,
	}
}

func (c *command) Name() string {
	return "delete"
}

func (c *command) Description() string {
	return "<id>;<password> - delete user"
}

func (c *command) Process(args string) string {
	params := strings.Split(args, ";")
	if len(params) != 2 {
		return commandPkg.MsgInvalidArguments
	}

	id, _ := strconv.ParseUint(params[0], 10, 64)
	user, err := c.user.Get(uint(id))
	if err != nil {
		return errors.Wrap(err, msgDeleteUser).Error()
	}

	// verify password to confirm changes
	err = auth.VerifyPassword(*user, params[1])
	if err != nil {
		return errors.Wrap(err, msgDeleteUser).Error()
	}

	if err := c.user.Delete(uint(id)); err != nil {
		return errors.Wrap(err, msgDeleteUser).Error()
	}

	return "user has been deleted"
}
