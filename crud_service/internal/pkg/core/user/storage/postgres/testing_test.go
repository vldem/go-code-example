package postgres

import (
	"testing"

	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	storagePkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/storage"
)

type usersTestFixture struct {
	usersRepo storagePkg.Interface
	data      models.User
}

func setUp(t *testing.T) usersTestFixture {
	var fixture usersTestFixture
	fixture.data = models.User{
		Id:       1,
		Email:    "test01@dummy.com",
		Name:     "Test Tester",
		Role:     "Admin",
		Password: "123456",
	}
	return fixture
}

// func (f *usersRepoFixtures) tearDown() {

// }
