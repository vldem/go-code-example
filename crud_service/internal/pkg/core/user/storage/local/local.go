package local

import (
	"context"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	storagePkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/storage"
)

const poolSize = 10

var ErrUserNotExists = errors.New("user does not exists")
var ErrUserExists = errors.New("user already exists")

var lastId uint

type storage struct {
	mu     sync.RWMutex
	data   map[uint]models.User
	poolCh chan struct{}
}

func New() storagePkg.Interface {
	return &storage{
		mu:     sync.RWMutex{},
		data:   map[uint]models.User{},
		poolCh: make(chan struct{}, poolSize),
	}
}

func (s *storage) List(ctx context.Context, _, _ uint64, _ models.SortingOrder) ([]models.User, error) {
	s.poolCh <- struct{}{}
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
		<-s.poolCh
	}()

	result := make([]models.User, 0, len(s.data))
	for _, value := range s.data {
		result = append(result, value)
	}
	return result, nil
}

func (s *storage) Add(ctx context.Context, user models.User) (uint, error) {
	s.poolCh <- struct{}{}
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
		<-s.poolCh
	}()

	if _, ok := s.data[user.Id]; ok {
		return 0, errors.Wrapf(ErrUserExists, "user-id: [%s]", strconv.FormatUint(uint64(user.Id), 10))
	}
	lastId++
	user.Id = lastId
	s.data[user.Id] = user
	return user.Id, nil
}

func (s *storage) Update(ctx context.Context, user models.User) error {
	s.poolCh <- struct{}{}
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
		<-s.poolCh
	}()

	if _, ok := s.data[user.Id]; !ok {
		return errors.Wrapf(ErrUserNotExists, "user-id: [%s]", strconv.FormatUint(uint64(user.Id), 10))
	}
	s.data[user.Id] = user
	return nil
}

func (s *storage) Delete(ctx context.Context, id uint) error {
	s.poolCh <- struct{}{}
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
		<-s.poolCh
	}()

	if _, ok := s.data[id]; !ok {
		return errors.Wrapf(ErrUserNotExists, "user-id: [%s]", strconv.FormatUint(uint64(id), 10))
	}
	delete(s.data, id)
	return nil
}

func (s *storage) Get(ctx context.Context, id uint) (*models.User, error) {
	s.poolCh <- struct{}{}
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
		<-s.poolCh
	}()

	user, ok := s.data[id]
	if !ok {
		return nil, errors.Wrapf(ErrUserNotExists, "user-id: [%s]", strconv.FormatUint(uint64(id), 10))
	}
	return &user, nil
}

func (s *storage) GetRoleIdByName(ctx context.Context, roleName string) (uint8, error) {
	roleId := models.GetRoleId(roleName)
	var err error
	if roleId == 0 {
		err = errors.Errorf("role [%s] not found", roleName)
	}
	return roleId, err
}
