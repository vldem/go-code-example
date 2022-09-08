// This model contains user object and methods
//go:generate mockgen -source=./user.go -destination=./mocks/user.go -package=mock_user
package user

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/vldem/homework1/internal/config"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	storagePkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/storage"
	postgresStoragePkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/storage/postgres"
	"golang.org/x/net/context"
)

type Interface interface {
	Create(ctx context.Context, user models.User) (uint, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (*models.User, error)
	List(ctx context.Context, recPerPage uint64, pageNum uint64, sortingOrder models.SortingOrder) ([]models.User, error)
	GetRoleIdByName(ctx context.Context, roleName string) (uint8, error)
}

type core struct {
	storage storagePkg.Interface
}

func New(pool *pgxpool.Pool) Interface {
	return &core{
		//storage: localStoragePkg.New(),
		storage: postgresStoragePkg.New(pool),
	}
}

func (c *core) Create(ctx context.Context, user models.User) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ShortDuration)
	defer cancel()
	timeOutCh := make(chan struct{}, 1)

	var id uint
	var err error

	go func(ch chan struct{}) {
		id, err = c.storage.Add(ctx, user)
		ch <- struct{}{}
	}(timeOutCh)

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-timeOutCh:
	}

	return id, err
}

func (c *core) Update(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.ShortDuration)
	defer cancel()

	timeOutCh := make(chan struct{}, 1)

	var err error

	go func(ch chan struct{}) {
		err = c.storage.Update(ctx, user)
		ch <- struct{}{}
	}(timeOutCh)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timeOutCh:
	}

	return err
}

func (c *core) Delete(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, config.ShortDuration)
	defer cancel()

	timeOutCh := make(chan struct{}, 1)
	var err error

	go func(ch chan struct{}) {
		err = c.storage.Delete(ctx, id)
		ch <- struct{}{}
	}(timeOutCh)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timeOutCh:
	}

	return err
}

func (c *core) Get(ctx context.Context, id uint) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ShortDuration)
	defer cancel()
	timeOutCh := make(chan struct{}, 1)

	var result *models.User
	var err error

	go func(ch chan struct{}) {
		result, err = c.storage.Get(ctx, id)
		ch <- struct{}{}
	}(timeOutCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timeOutCh:
	}
	return result, err
}

func (c *core) List(ctx context.Context, recPerPage uint64, pageNum uint64, sortingOrder models.SortingOrder) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ShortDuration)
	defer cancel()
	timeOutCh := make(chan struct{}, 1)

	var result []models.User
	var err error

	go func(ch chan struct{}) {
		result, err = c.storage.List(ctx, recPerPage, pageNum, sortingOrder)
		ch <- struct{}{}
	}(timeOutCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timeOutCh:
	}

	return result, err
}

func (c *core) GetRoleIdByName(ctx context.Context, roleName string) (uint8, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ShortDuration)
	defer cancel()
	timeOutCh := make(chan struct{}, 1)

	var result uint8
	var err error

	go func(ch chan struct{}) {
		result, err = c.storage.GetRoleIdByName(ctx, roleName)
		ch <- struct{}{}
	}(timeOutCh)

	select {
	case <-ctx.Done():
		return result, ctx.Err()
	case <-timeOutCh:
	}
	return result, err
}
