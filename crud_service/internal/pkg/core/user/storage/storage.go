//go:generate mockgen -source=./storage.go -destination=./mocks/storage.go -package=mock_storage

// This is a storage that contains list of users
package storage

import (
	"context"

	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
)

type Interface interface {
	Add(ctx context.Context, user models.User) (uint, error)
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, user models.User) error
	List(ctx context.Context, recPerPage, pageNum uint64, sortingOrder models.SortingOrder) ([]models.User, error)
	GetRoleIdByName(ctx context.Context, role string) (uint8, error)
}
