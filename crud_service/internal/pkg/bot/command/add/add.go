package add

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/vldem/homework1/internal/auth"
	commandPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/bot/command"
	validatorPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/validator"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
)

const msgAddUser = "add user"

type command struct {
	client pb.BackendClient
}

func New(client pb.BackendClient) commandPkg.Interface {
	return &command{
		client: client,
	}
}

func (c *command) Name() string {
	return "add"
}

func (c *command) Description() string {
	return "<email>;<name>;<role>;<password> - create user"
}

func (c *command) Process(ctx context.Context, args string) string {
	params := strings.Split(args, ";")
	if len(params) != 4 {
		return commandPkg.MsgInvalidArguments
	}

	//validate parameters
	err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate(params))
	if err != nil {
		return errors.Wrap(err, msgAddUser).Error()
	}

	id, err := c.client.UserCreate(ctx, &pb.BackendUserCreateRequest{
		Email:    params[0],
		Name:     params[1],
		Role:     params[2],
		Password: auth.GenHashPassword(params[3]),
	})
	if err != nil {
		return errors.Wrap(err, msgAddUser).Error()
	}

	return fmt.Sprintf("user [%v] has been added", id)
}
