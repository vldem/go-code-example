package delete

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	commandPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/bot/command"
	validatorPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/validator"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
)

var msgDeleteUser = "delete user"

type command struct {
	client pb.BackendClient
}

func New(client pb.BackendClient) commandPkg.Interface {
	return &command{
		client: client,
	}
}

func (c *command) Name() string {
	return "delete"
}

func (c *command) Description() string {
	return "<id>;<password> - delete user"
}

func (c *command) Process(ctx context.Context, args string) string {
	params := strings.Split(args, ";")
	if len(params) != 2 {
		return commandPkg.MsgInvalidArguments
	}

	//validate user's input
	if err := validatorPkg.ValidateUserId(params[0]); err != nil {
		return errors.Wrap(err, msgDeleteUser).Error()
	}
	if err := validatorPkg.ValidatePassword(params[1]); err != nil {
		return errors.Wrap(err, msgDeleteUser).Error()
	}

	id, _ := strconv.ParseUint(params[0], 10, 64)

	_, err := c.client.UserDelete(ctx, &pb.BackendUserDeleteRequest{
		Id:       id,
		Password: params[1],
	})
	if err != nil {
		return errors.Wrap(err, msgDeleteUser).Error()
	}

	return "user has been deleted"
}
