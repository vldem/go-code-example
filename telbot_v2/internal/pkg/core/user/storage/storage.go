// This is a storage that contains list of users
package storage

import (
	"github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/models"
)

type Interface interface {
	Add(user models.User) error
	Delete(id uint) error
	Get(id uint) (*models.User, error)
	Update(user models.User) error
	List() []models.User
}
