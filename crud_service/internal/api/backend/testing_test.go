package backend

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_repository "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/mocks"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
)

type backendFixture struct {
	Ctx      context.Context
	userRepo *mock_repository.MockInterface
	service  *implementation
	data     models.User
	list     []models.User
}

type userListFixture struct {
	Ctx        context.Context
	userRepo   *mock_repository.MockInterface
	service    *implementation
	list       []*pb.BackendUserListResponse_User
	data       []models.User
	recPerPage uint64
	pageNum    uint64
	order      models.SortingOrder
}

func userSetUp(t *testing.T) backendFixture {
	t.Parallel()

	f := backendFixture{Ctx: context.Background()}
	f.userRepo = mock_repository.NewMockInterface(gomock.NewController(t))
	f.service = New(f.userRepo)
	f.data = models.User{
		Id:       1,
		Email:    "test01@dummy.com",
		Name:     "Test Tester",
		Role:     "Admin",
		Password: "123456",
	}
	return f
}

func userListSetUp(t *testing.T) userListFixture {
	t.Parallel()

	f := userListFixture{
		Ctx:  context.Background(),
		data: []models.User{},
		list: []*pb.BackendUserListResponse_User{},
	}
	f.userRepo = mock_repository.NewMockInterface(gomock.NewController(t))
	f.service = New(f.userRepo)
	f.recPerPage = 10
	f.pageNum = 1
	f.order = models.SortingOrder{
		Field:      "email",
		Descending: false,
	}
	f.list = append(f.list, &pb.BackendUserListResponse_User{
		Id:    1,
		Email: "test01@dummy.com",
		Name:  "Test Tester",
		Role:  "Admin",
	})
	f.list = append(f.list, &pb.BackendUserListResponse_User{
		Id:    2,
		Email: "test02@dummy.com",
		Name:  "Test2 Tester2",
		Role:  "User",
	})
	f.data = append(f.data, models.User{
		Id:    1,
		Email: "test01@dummy.com",
		Name:  "Test Tester",
		Role:  "Admin",
	})
	f.data = append(f.data, models.User{
		Id:    2,
		Email: "test02@dummy.com",
		Name:  "Test2 Tester2",
		Role:  "User",
	})
	return f
}
