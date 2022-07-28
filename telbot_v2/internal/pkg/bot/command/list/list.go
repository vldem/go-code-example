package list

import (
	"fmt"
	"strings"

	commandPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/bot/command"
	userPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user"
	"github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/models"
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
	return "list"
}

func (c *command) Description() string {
	return " - list of users"
}

func (c *command) Process(_ string) string {
	result := []string{}

	users := c.user.List()
	if len(users) == 0 {
		return "no users found"
	}

	for _, user := range users {
		result = append(result, fmt.Sprintf("%d: %s / %s / %s ", user.Id, user.Email, user.Name, models.GetRoleName(user.Role)))
	}

	return strings.Join(result, "\n")
}
