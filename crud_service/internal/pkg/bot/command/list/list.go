package list

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/vldem/homework1/internal/config"
	commandPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/bot/command"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
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
	return "list"
}

func (c *command) Description() string {
	return " [<record per page>;<page number>;<sorting field:id|name|email>;<descending: true|false>]- list of users"
}

func (c *command) Process(ctx context.Context, args string) string {
	result := []string{}
	params := strings.Split(args, ";")
	var err error
	var recPerPage, pageNum uint64
	if len(params) >= 2 {
		recPerPage, err = strconv.ParseUint(params[0], 10, 64)
		if err != nil {
			return commandPkg.MsgInvalidArguments
		}
		pageNum, err = strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return commandPkg.MsgInvalidArguments
		}
	}
	sortingOrder := models.SortingOrder{
		Field:      config.DefaultSortingField,
		Descending: false,
	}

	if len(params) == 4 {
		err := validatorPkg.ValidateSortingField(params[2])
		if err != nil {
			return commandPkg.MsgInvalidArguments
		}

		descending, err := strconv.ParseBool(params[3])
		if err != nil {
			return commandPkg.MsgInvalidArguments
		}
		sortingOrder.Field = params[2]
		sortingOrder.Descending = descending
	}

	if recPerPage == 0 {
		recPerPage = config.DefaultRecPerPage
	}
	if pageNum == 0 {
		pageNum = config.DefaultPageNum
	}

	users, err := c.client.UserList(ctx, &pb.BackendUserListRequest{
		RecPerPage: &recPerPage,
		PageNum:    &pageNum,
		Order: &pb.BackendUserListRequest_SortingOrder{
			Field:      sortingOrder.Field,
			Descending: sortingOrder.Descending,
		},
	})
	if err != nil {
		return errors.Wrap(err, "internal error").Error()
	}
	if len(users.GetUsers()) == 0 {
		return "no users found"
	}

	for _, user := range users.GetUsers() {
		result = append(result, fmt.Sprintf("%d: %s / %s / %s ", user.Id, user.Email, user.Name, user.Role))
	}

	return strings.Join(result, "\n")
}
