package local

import (
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/models"
	storagePkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/storage"
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

func (s *storage) List() []models.User {
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
	return result
}

func (s *storage) Add(user models.User) error {
	s.poolCh <- struct{}{}
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
		<-s.poolCh
	}()

	if _, ok := s.data[user.Id]; ok {
		return errors.Wrapf(ErrUserExists, "user-id: [%s]", strconv.FormatUint(uint64(user.Id), 10))
	}
	lastId++
	user.Id = lastId
	s.data[user.Id] = user
	return nil
}

func (s *storage) Update(user models.User) error {
	s.poolCh <- struct{}{}
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
		<-s.poolCh
	}()

	if _, ok := s.data[user.Id]; !ok {
		return errors.Wrapf(ErrUserNotExists, "user-id: [%s]", strconv.FormatUint(uint64(user.Id), 10))
	}
	s.data[user.Id] = user
	return nil
}

func (s *storage) Delete(id uint) error {
	s.poolCh <- struct{}{}
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
		<-s.poolCh
	}()

	if _, ok := s.data[id]; !ok {
		return errors.Wrapf(ErrUserNotExists, "user-id: [%s]", strconv.FormatUint(uint64(id), 10))
	}
	delete(s.data, id)
	return nil
}

func (s *storage) Get(id uint) (*models.User, error) {
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
