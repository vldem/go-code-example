package update

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/vldem/homework1/internal/auth"
	commandPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/bot/command"
	validatorPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/validator"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
)

var msgUpdateUser = "update user"

type command struct {
	client pb.BackendClient
}

func New(client pb.BackendClient) commandPkg.Interface {
	return &command{
		client: client,
	}
}

func (c *command) Name() string {
	return "update"
}

func (c *command) Description() string {
	return "<id>;<email>;<name>;<role>;<password>;<old password> - update user"
}

func (c *command) Process(ctx context.Context, args string) string {
	params := strings.Split(args, ";")
	if len(params) != 6 {
		return commandPkg.MsgInvalidArguments
	}

	//validate parameters
	if err := validatorPkg.ValidateUserId(params[0]); err != nil {
		return errors.Wrap(err, msgUpdateUser).Error()
	}
	if err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate(params[1 : len(params)-1])); err != nil {
		return errors.Wrap(err, msgUpdateUser).Error()
	}

	id, _ := strconv.ParseUint(params[0], 10, 64)

	if _, err := c.client.UserUpdate(ctx, &pb.BackendUserUpdateRequest{
		Id:       id,
		Email:    params[1],
		Name:     params[2],
		Role:     params[3],
		Password: auth.GenHashPassword(params[4]),
	}); err != nil {
		return errors.Wrap(err, msgUpdateUser).Error()
	}

	return "user has been updated"
}
